package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

const layoutFile = "./layout.html"
const blogFile = "./blog/index.html"
const rssFile = "./blog/rss.xml"
const searchIndexFile = "searchIndex.js"

var port int
var serverMode bool

func init() {
	flag.IntVar(&port, "port", 9001, "Listening port when on server mode")
	flag.BoolVar(&serverMode, "server", false, "Pass true to launch a local web server")

	flag.Usage = func() { showSyntax() }
	flag.Parse()
}

func main() {
	fmt.Printf("fixie - a one gear static site generator\r\n\r\n")
	processMarkdownFiles()
	if serverMode == true {
		server(port)
	}
}

func processMarkdownFiles() {
	layout := loadLayout(layoutFile)
	siteMetadata := NewSiteMeta(layout)
	blogs := BlogPosts{}

	fmt.Printf("Processing .md files...\r\n")
	filepath.WalkDir(".", func(filename string, d fs.DirEntry, err error) error {
		if filepath.Ext(filename) != ".md" || filename == "README.md" {
			return nil
		}

		// Create the HTML version of the Markdown file
		fmt.Printf("  %s\r\n", filename)
		md2HtmlFile(filename, layout)

		// Keep track of blog entries
		isBlog := strings.HasPrefix(filename, "blog/")
		if isBlog {
			blogPost := LoadBlogPost(filename)
			blogPost.createRedirectFiles()
			blogs.Append(blogPost)
		}
		return nil
	})

	if len(blogs) == 0 {
		fmt.Printf("No blog entries (./blog/) were found\r\n")
	} else {
		fmt.Printf("%d blog entries were processed\r\n", len(blogs))
		blogs.CreateHomepage(layout, blogFile)
		blogs.CreateRssPage(siteMetadata, rssFile)
		//
		//
		//
		// TODO: Include the pages that are not blogs on the index
		//
		//
		//
		blogs.CreateSearchIndex(searchIndexFile)
	}
	return
}

func loadLayout(layoutFile string) string {
	if fileExist(layoutFile) {
		fmt.Printf("Using layout file: %s\r\n", layoutFile)
		return readFile(layoutFile)
	}

	fmt.Printf("No layout file (%s) was found\r\n", layoutFile)
	return ""
}

func showSyntax() {
	fmt.Printf("fixie - a one gear static site generator\r\n")
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
