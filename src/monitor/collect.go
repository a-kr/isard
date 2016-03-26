package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func Collect(item string) error {
	collectScript := path.Join(*dataDir, item, "collect.sh")
	collectedData := path.Join(*dataDir, item, "data.txt")

	log.Printf("[%s] Running %s", item, collectScript)
	var stderr bytes.Buffer
	cmd := exec.Command(collectScript)
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		log.Printf("[%s] ERROR running %s: %s", item, collectScript, err)
		log.Printf("[%s] ERROR stderr: %s", item, stderr.String())
		return err
	}

	log.Printf("[%s] Output: %s", item, out)

	outLine := strings.TrimSpace(string(out))
	f, err := os.OpenFile(collectedData, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	f.WriteString(outLine + "\n")
	f.Close()
	return nil
}
