package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const layoutFile = "./layout.html"
const blogFile = "./blog/index.html"

var port int
var serverMode bool
var parser MarkdownParser

func init() {
	flag.IntVar(&port, "port", 9001, "Listening port when on server mode")
	flag.BoolVar(&serverMode, "server", false, "Pass true to launch a local web server")
	flag.Usage = func() { showSyntax() }
	flag.Parse()
}

func main() {
	fmt.Printf("fixie - a one gear blog engine\r\n")
	if serverMode == true {
		server(port)
		os.Exit(0)
	}
	processFiles()
}

func processFiles() {
	layout := readFile(layoutFile)
	blogPosts := BlogPosts{}

	err := filepath.WalkDir(".", func(fileName string, d fs.DirEntry, err error) error {
		if filepath.Ext(fileName) != ".md" || fileName == "README.md" {
			return nil
		}

		fmt.Printf("Processing file: %s\r\n", fileName)

		// Create the HTML version of the Markdown file
		content := readFile(fileName)
		htmlFile := strings.TrimSuffix(fileName, ".md") + ".html"
		md2HtmlFile(layout, content, htmlFile)

		isBlog := strings.HasPrefix(fileName, "blog/")
		if isBlog {
			blogPost := BlogPost{Filename: fileName, Content: content}
			blogPosts.Append(blogPost)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	blogPosts.CreateHomepage(layout, blogFile)
	fmt.Printf("Done\r\n")
}

func md2Html(layout string, content string) string {
	contentHtml := parser.ToHtml(content)
	if layout == "" {
		return contentHtml
	}
	return strings.Replace(layout, "{{ content }}", contentHtml, 1)
}

func md2HtmlFile(layout string, content string, htmlFile string) {
	html := md2Html(layout, content)
	saveFile(htmlFile, html)
}

func mdTitle(content string, defaultTitle string) string {
	title := parser.Title(content)
	if title != "" {
		return title
	}
	return defaultTitle
}

func showSyntax() {
	fmt.Printf("fixie - a one gear blog engine\r\n")
	flag.PrintDefaults()
	fmt.Printf(`
NOTES:
`)
	fmt.Printf("\r\n")
	fmt.Printf("\r\n")
}
