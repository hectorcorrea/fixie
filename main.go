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
const rssFile = "./blog/rss.xml"

var port int
var serverMode bool
var parser MarkdownParser
var layout string
var blogs BlogPosts

type SiteMeta struct {
	Title       string
	Author      string
	Description string
	Link        string
}

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

	blogs = BlogPosts{}
	layout = readFile(layoutFile)
	meta := htmlMeta(layout)
	processMarkdownFiles()
	blogs.CreateHomepage(layout, blogFile)
	blogs.CreateRssPage(meta, rssFile)

	fmt.Printf("Done\r\n")
}

func processMarkdownFiles() {
	fmt.Printf("Processing .md files...\r\n")
	err := filepath.WalkDir(".", processFile)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}

func processFile(fileName string, d fs.DirEntry, err error) error {
	if filepath.Ext(fileName) != ".md" || fileName == "README.md" {
		return nil
	}

	// Create the HTML version of the Markdown file
	fmt.Printf("  %s\r\n", fileName)

	// Keep track of blog posts (used for the blog homepage later on)
	isBlog := strings.HasPrefix(fileName, "blog/")
	if isBlog {
		content := readFile(fileName)
		blogPost := BlogPost{Filename: fileName, Content: content}
		blogs.Append(blogPost)
		md2HtmlFile(fileName, layout)
	} else {
		md2HtmlFile(fileName, layout)
	}
	return nil
}

func showSyntax() {
	fmt.Printf("fixie - a one gear blog engine\r\n")
	flag.PrintDefaults()
	fmt.Printf(`
Process all markdown files (.md) on the current directory and generates the HTML version
for each of them.

Fixie is very opinioned.

If there is a layout.html on the current folder this file will be used as the layout for
the generated HTML files. The layout.html file must include a {{CONTENT}} token where
you expect the content of each Markdown file to be inserted.

Files under the blog folder are considered blog posts and will be handled slightly
different. They will be converted from Markdown to HTML using the same layout as the one
for all other files, but also

1. The title of each blog post is taken from the first line of the Markdown if and only if
the first line is a H1 Heading (e.g. # My Blog Post). Otherwise the file name is used as
the title.

2. A file /blog/index.html will created with a list of all the files processed sorted
descending by created date. The created date for blog posts is calculated from the file
name or the folder for each file. If the file name starts with a date in the format
YYYY-MM-DD this will be used as the created date for the blog post. If the file is
inside a folder that starts with YYYY-MM-DD that will be used as the created date.

3. An blog.rss file will generated for all blog posts processed.
`)
	fmt.Printf("\r\n")
	fmt.Printf("\r\n")
}
