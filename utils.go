package main

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func md2HtmlFile(mdFile string, layout string) {
	content := readFile(mdFile)
	htmlFile := strings.TrimSuffix(mdFile, ".md") + ".html"
	html := md2Html(layout, content)
	saveFile(htmlFile, html)
}

func md2Html(layout string, content string) string {
	var parser MarkdownParser
	contentHtml := parser.ToHtml(content)
	if layout != "" {
		// If there is layout merge the HTML content with the layout
		contentHtml = strings.Replace(layout, "{{CONTENT}}", contentHtml, 1)
	}

	// If there is an HTML title tag in the original layout and there is a title in the markdown
	// replace the original title with the one from the markdown
	reTitle := regexp.MustCompile(`<title>(.*)</title>`)
	originalTitle := regExpMatch(layout, reTitle)
	markdownTitle := parser.Title(content)
	if originalTitle != "" && markdownTitle != "" {
		contentHtml = strings.Replace(contentHtml, "<title>"+originalTitle+"</title>", "<title>"+markdownTitle+"</title>", 1)
	}

	return contentHtml
}

func mdTitle(content string, defaultTitle string) string {
	var parser MarkdownParser
	title := parser.Title(content)
	if title != "" {
		return title
	}
	return defaultTitle
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func readFile(filename string) string {
	bytes, _ := os.ReadFile(filename)
	return string(bytes)
}

func saveFile(filename string, content string) {
	os.WriteFile(filename, []byte(content), 0644)
}

func createDir(dirname string) {
	os.MkdirAll(dirname, os.ModePerm)
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
	reDate := regexp.MustCompile(`\d\d\d\d\-\d\d\-\d\d`)
	match := string(reDate.Find([]byte(value)))
	return match
}

// Returns the first match of an regex in the provided text
func regExpMatch(text string, re *regexp.Regexp) string {
	matches := re.FindStringSubmatch(text)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}
