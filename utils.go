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
	parser := NewMarkdownParser(content)
	contentHtml := parser.Html()
	if layout != "" {
		// If there is layout merge the HTML content with the layout
		contentHtml = strings.Replace(layout, "{{CONTENT}}", contentHtml, 1)
	}

	// Replace the Title in the original HTML with the one from the Markdown
	reTitle := regexp.MustCompile(`<title>(.*)</title>`)
	htmlTitle := regExpMatch(layout, reTitle)
	if htmlTitle != "" && parser.Title() != "" {
		contentHtml = strings.Replace(contentHtml, "<title>"+htmlTitle+"</title>", "<title>"+parser.Title()+"</title>", 1)
	}

	// Replace the Description in the original HTML with the one from the Markdown
	reDescription := regexp.MustCompile(`<meta name="description" content="(.*)">`)
	htmlDesc := regExpMatch(layout, reDescription)
	if htmlDesc != "" && parser.Description() != "" {
		oldHtml := `<meta name="description" content="` + htmlDesc + `">`
		newHtml := `<meta name="description" content="` + parser.Description() + `">`
		contentHtml = strings.Replace(contentHtml, oldHtml, newHtml, 1)
	}

	// TODO: Replace the canonical link on the HTML (<link rel="canonical" href="https://something" />)
	// with the correct link for each page.
	//
	// The issue is that at this stage I don't know what base URL I should use, heck I don't even know
	// the URL for this page (i.e. the one based from the Title of the post) so that complicates things
	// a little bit. I might need to implement this logic as part of the blog/page processing code.

	return contentHtml
}

func mdTitle(content string, defaultTitle string) string {
	parser := NewMarkdownParser(content)
	if parser.Title() != "" {
		return parser.Title()
	}
	return defaultTitle
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	notExist := errors.Is(err, os.ErrNotExist)
	return !notExist
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
