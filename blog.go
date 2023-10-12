package main

import (
	"cmp"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

type BlogPost struct {
	Filename string
	Content  string
}

type BlogPosts []BlogPost

func (blogs BlogPosts) SortDescending() {
	slices.SortFunc(blogs, func(a, b BlogPost) int {
		return cmp.Compare(b.DateCreated(), a.DateCreated())
	})
}

func (blogs *BlogPosts) Append(blog BlogPost) {
	// https://stackoverflow.com/a/24726368/446681
	*blogs = append(*blogs, blog)
}

func (blogs BlogPosts) Content() string {
	content := "# Blog Posts\r\n"
	blogs.SortDescending()
	for _, blog := range blogs {
		content += "* " + blog.LinkMarkdown() + "\r\n"
	}
	return content
}

func (blogs BlogPosts) CreateHomepage(layout string, filename string) {
	if len(blogs) == 0 {
		return
	}
	fmt.Printf("Creating blog homepage: %s\r\n", filename)
	html := md2Html(layout, blogs.Content())
	saveFile(filename, html)
}

func (blogs BlogPosts) CreateRssPage(meta SiteMeta, filename string) {
	if len(blogs) == 0 {
		return
	}
	fmt.Printf("Creating blog RSS: %s\r\n", filename)

	rss := NewRss(meta.Title, meta.Description, meta.Link)
	for _, blog := range blogs {
		rss.Add(blog.Title(), blog.Summary(), blog.LinkUrl(), blog.DateCreated())
	}
	xml, err := rss.ToXml()
	if err != nil {
		fmt.Printf("ERROR producing RSS file: %s\r\n", err)
	}
	saveFile(filename, xml)
}

func (b BlogPost) DateCreated() string {
	date := dateFromFilename(b.Filename)
	if date == "" {
		fmt.Printf("No created date found for blog: %s\r\n", b.Filename)
		return "1900-01-01"
	}
	return date
}

func (b BlogPost) Summary() string {
	return "TODO"
}

func (b BlogPost) DefaultTitle() string {
	return strings.TrimSuffix(filepath.Base(b.Filename), ".md")
}

func (b BlogPost) Title() string {
	return mdTitle(b.Content, b.DefaultTitle())
}

func (b BlogPost) LinkUrl() string {
	return "/" + strings.TrimSuffix(b.Filename, ".md")
}

func (b BlogPost) LinkMarkdown() string {
	return fmt.Sprintf("[%s](%s)", b.Title(), b.LinkUrl())
}
