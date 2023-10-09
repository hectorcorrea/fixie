package main

import (
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
	contentHtml := parser.ToHtml(content)
	if layout == "" {
		return contentHtml
	}
	return strings.Replace(layout, "{{CONTENT}}", contentHtml, 1)
}

func mdTitle(content string, defaultTitle string) string {
	title := parser.Title(content)
	if title != "" {
		return title
	}
	return defaultTitle
}

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

// Extract site metadata from an HTML page
func htmlMeta(html string) SiteMeta {
	meta := SiteMeta{
		Title:       htmlTitle(html),
		Description: htmlMetaAttr(html, "description"),
		Author:      htmlMetaAttr(html, "author"),
		Link:        htmlLinkCanonical(html),
	}
	return meta
}

// Extracts the value of a meta attribute in an HTML page
func htmlMetaAttr(html string, metaName string) string {
	re := regexp.MustCompile(`<meta name="` + metaName + `" content="(.*)"`)
	return regExpMatch(html, re)
}

// Extracts the title of an HTML page
func htmlTitle(html string) string {
	re := regexp.MustCompile(`<title>(.*)</title>`)
	return regExpMatch(html, re)
}

// Extracts the canonical link value of an HTML page
func htmlLinkCanonical(html string) string {
	re := regexp.MustCompile(`<link rel="canonical" href="(.*)"`)
	return regExpMatch(html, re)
}

// Returs the first match of an regex in the provided text
func regExpMatch(text string, re *regexp.Regexp) string {
	matches := re.FindStringSubmatch(text)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}
