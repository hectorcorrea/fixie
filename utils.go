package main

import (
	"os"
	"path/filepath"
	"regexp"
)

func readFile(filename string) string {
	bytes, _ := os.ReadFile(filename)
	return string(bytes)
}

func saveFile(filename string, content string) {
	os.WriteFile(filename, []byte(content), 0644)
}

func dateFromFilename(fullPath string) string {
	// Get the date from the file name...
	date := dateFromString(filepath.Base(fullPath))
	if date == "" {
		// ...see if we can get a date from the path
		date = dateFromString(filepath.Dir(fullPath))
	}
	return date
}

func dateFromString(value string) string {
	reDate := regexp.MustCompile("\\d\\d\\d\\d\\-\\d\\d\\-\\d\\d")
	match := string(reDate.Find([]byte(value)))
	return match
}
