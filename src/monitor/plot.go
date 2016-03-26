package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"text/template"
	"time"
)

const (
	GnuplotTimeFormat       = "2006-01-02 15:04:05"
	XTicLabelWithMinutes    = "%H:%M\\n%d.%m"
	XTicLabelWithDayAndYear = "%d.%m\\n%Y"
)

type PlotOptions struct {
	Width, Height int
	LastNHours    int
	LastNDays     int
	XTics         int
	XFormat       string
	InputFile     string
	FromDateTime  string
	SaveAsLatest  bool
}

func (opts *PlotOptions) FillDefaults() {
	if opts.Width == 0 {
		opts.Width = 1400
	}
	if opts.Height == 0 {
		opts.Height = 800
	}
	if opts.LastNHours == 0 {
		if opts.LastNDays == 0 {
			opts.LastNHours = 72
		} else {
			opts.LastNHours = opts.LastNDays * 24
		}
	}
	opts.FromDateTime = time.Now().Add(-time.Duration(opts.LastNHours) * time.Hour).Format(GnuplotTimeFormat)
	if opts.XTics == 0 {
		minPixelsPerTicLabel := 30
		fittingLabels := (opts.Width - 150) / minPixelsPerTicLabel
		secondsPerTic := opts.LastNHours * 60 * 60 / fittingLabels

		opts.XFormat = XTicLabelWithMinutes

		if secondsPerTic < 60*15 {
			secondsPerTic = 60 * 15
		} else if secondsPerTic < 60*30 {
			secondsPerTic = 60 * 30
		} else if secondsPerTic < 60*60 {
			secondsPerTic = 60 * 60
		} else if secondsPerTic < 60*60*24 {
			daysPerTicf := math.Ceil(float64(secondsPerTic) / float64(60*60))
			secondsPerTic = int(daysPerTicf * 60 * 60)
		} else {
			// round up to nearest multiple of a day
			daysPerTicf := math.Ceil(float64(secondsPerTic) / float64(60*60*24))
			secondsPerTic = int(daysPerTicf * 60 * 60 * 24)
			opts.XFormat = XTicLabelWithDayAndYear
		}
		opts.XTics = secondsPerTic
	}
}

func (opts *PlotOptions) ProcessTemplate(templateText string) string {
	var result bytes.Buffer
	t := template.Must(template.New("gnuplot").Parse(templateText))
	err := t.Execute(&result, opts)
	if err != nil {
		panic(err)
	}
	return result.String()
}

type PngImage struct {
	Data []byte
}

func Plot(item string, opts PlotOptions) (*PngImage, error) {
	opts.FillDefaults()

	collectedData := path.Join(*dataDir, item, "data.txt")
	plotTemplate := path.Join(*dataDir, item, "plot.gnuplot")
	plotTemplateData, err := ioutil.ReadFile(plotTemplate)
	if err != nil {
		return nil, err
	}

	opts.InputFile = collectedData
	gnuplotScript := opts.ProcessTemplate(string(plotTemplateData))

	gnuplotScriptFile, err := ioutil.TempFile("", "monitor_gnuplot_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(gnuplotScriptFile.Name())
	gnuplotScriptFile.WriteString(gnuplotScript + "\n")
	gnuplotScriptFile.Close()

	var stderr bytes.Buffer
	log.Printf("Running gnuplot %s", gnuplotScriptFile.Name())
	cmd := exec.Command("gnuplot", gnuplotScriptFile.Name())
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		log.Printf("ERROR running gnuplot %s: %s", gnuplotScriptFile.Name(), err)
		log.Printf("ERROR stderr: %s", stderr.String())
		log.Printf("ERROR plot data:\n%s", gnuplotScript)
		return nil, err
	}

	if opts.SaveAsLatest {
		latestImage := path.Join(*dataDir, item, "latest.png")
		err := ioutil.WriteFile(latestImage, out, 0666)
		if err != nil {
			return nil, err
		}
	}

	return &PngImage{out}, nil
}
