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

func (blogs BlogPosts) Empty() bool {
	return len(blogs) == 0
}

func (blogs BlogPosts) Content() string {
	content := ""
	blogs.SortDescending()
	for _, blog := range blogs {
		content += blog.LinkMarkdown() + "\r\n"
	}
	return content
}

func (blogs BlogPosts) CreateHomepage(layout string, filename string) {
	if len(blogs) == 0 {
		return
	}
	fmt.Printf("Creating: %s\r\n", blogFile)
	md2HtmlFile(layout, blogs.Content(), filename)
}

func (b BlogPost) DateCreated() string {
	date := dateFromFilename(b.Filename)
	if date == "" {
		date = "1900-01-01"
	}
	return date
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
