package main

import (
	"io"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"text/template"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

func HttpStartServer(addr string) {
	http.HandleFunc("/", HttpView(HttpIndex))
	http.HandleFunc("/plot.png", HttpView(HttpPlot))
	http.HandleFunc("/data.csv", HttpView(HttpCsv))
	http.HandleFunc("/favicon.png", HttpView(HttpFavicon))
	log.Fatal(http.ListenAndServe(addr, nil))
}

func HttpView(handler HttpHandler) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 500)
			}
		}()
		handler(w, r)
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	IndexTemplate = template.Must(template.New("http_index").Parse(`
	<!DOCTYPE html>
	<html>
		<head>
			<title>Monitor</title>
			<link rel="icon" type="image/png" href="/favicon.png" />
		</head>
		<body>
		<style>
		div.item {
			width: auto;
			float: left;
			margin: 8px;
			padding: 8px;
			background-color: #eee;
		}
		a.taglink {
			color: gray;
			border-bottom: dashed 1px gray;
			text-decoration: none;
		}
		.item strong {
			margin-right: 20px;
		}
		</style>
		<div>
			<a href="{{ .BaseUrl }}hours=1">1h</a>
			<a href="{{ .BaseUrl }}hours=4">4h</a>
			<a href="{{ .BaseUrl }}hours=24">24h</a>
			<a href="{{ .BaseUrl }}hours=72">72h</a>
			<a href="{{ .BaseUrl }}hours=168">7d</a>
			<a href="{{ .BaseUrl }}hours=744">31d</a>
			<a href="{{ .BaseUrl }}hours=2232">3m</a>
		</div>
		<hr/>
		{{ range .Items }}
		<div class="item">
			<p>
				<strong>{{ .Name }}</strong>
				{{ range .Tags }}
					<a href="?tag={{ . }}" class="taglink">#{{ . }}</a>
				{{ end }}
			</p>
			<img src="/plot.png?item={{ .Name }}&hours={{ $.Hours }}&width=490&height=300" width="490" height="300"/>
		</div>
		{{ end }}
		</body>
	</html>
	`))
)

func filterItems(items []Item, names []string) []Item {
	result := make([]Item, 0)
	for _, x := range items {
		for _, s := range names {
			if x.Name == s {
				result = append(result, x)
				break
			}
		}
	}
	return result
}

func filterByTag(items []Item, tag string) []Item {
	result := make([]Item, 0)
	for _, x := range items {
		for _, t := range x.Tags {
			if t == tag {
				result = append(result, x)
				break
			}
		}
	}
	return result
}

func HttpIndex(w http.ResponseWriter, r *http.Request) {
	items, err := FindItems()
	panicOnError(err)

	r.ParseForm()

	baseurl := "/?"

	if r.FormValue("tag") != "" {
		items = filterByTag(items, r.FormValue("tag"))
		baseurl = "/?tag=" + r.FormValue("tag") + "&"
	}
	if len(r.Form["item"]) > 0 {
		items = filterItems(items, r.Form["item"])
	}

	hours := IntOrDefault(r.FormValue("hours"), 24)

	ctx := &struct {
		Hours   int
		Items   []Item
		BaseUrl string
	}{hours, items, baseurl}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = IndexTemplate.Execute(w, ctx)
	panicOnError(err)
}

func HttpPlot(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Printf("%s", r.Form)
	itemName := r.Form.Get("item")
	if itemName == "" {
		http.Error(w, "`item` required", 400)
		return
	}

	opts := PlotOptions{}
	opts.Width = IntOrDefault(r.Form.Get("width"), 1400)
	opts.Height = IntOrDefault(r.Form.Get("height"), 600)
	opts.LastNHours = IntOrDefault(r.Form.Get("hours"), 24)
	opts.LastNDays = IntOrDefault(r.Form.Get("days"), 0)

	png, err := Plot(itemName, opts)
	panicOnError(err)

	w.Header().Set("Content-Type", "image/png")
	w.Write(png.Data)
}

func HttpFavicon(w http.ResponseWriter, r *http.Request) {
	icon := Favicon()
	w.Header().Set("Content-Type", "image/png")
	err := png.Encode(w, icon)
	panicOnError(err)
}

func HttpCsv(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Printf("%s", r.Form)
	itemName := r.Form.Get("item")
	if itemName == "" {
		http.Error(w, "`item` required", 400)
		return
	}
	text, err := ReadData(itemName)
	panicOnError(err)
	io.WriteString(w, text)
}
