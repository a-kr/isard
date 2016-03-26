package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func dieOnError(err error) {
	if err != nil {
		log.Panicf("Error: %s", err)
	}
}

func defaultDataDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dieOnError(err)
	return path.Join(dir, "../data")
}

func IntOrDefault(x string, defaultValue int) int {
	i, err := strconv.Atoi(x)
	if err != nil {
		i = defaultValue
	}
	return i
}

func SameStringArrays(a1 []string, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}
	items := make(map[string]int)
	for _, x := range(a1) {
		items[x] = 1
	}
	for _, x := range(a2) {
		items[x] += 1
	}
	for _, n := range items {
		if n != 2 {
			return false
		}
	}
	return true
}
