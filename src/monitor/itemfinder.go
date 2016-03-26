package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Item struct {
	Name     string
	CronRule string
}

func FindItems() ([]Item, error) {
	result := []Item{}

	files, err := ioutil.ReadDir(*dataDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		collectScript := path.Join(*dataDir, f.Name(), "collect.sh")
		_, err := os.Stat(collectScript)
		if err != nil {
			continue
		}
		item := Item{
			Name: f.Name(),
		}

		cronPath := path.Join(*dataDir, f.Name(), "cron.txt")
		cronData, err := ioutil.ReadFile(cronPath)
		if err != nil {
			continue
		}
		cronDataStr := strings.Split(string(cronData), "\n")[0]
		item.CronRule = strings.TrimSpace(cronDataStr)

		result = append(result, item)
	}
	return result, nil
}
