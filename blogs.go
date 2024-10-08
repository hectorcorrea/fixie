package main

import (
	"cmp"
	"fmt"
	"slices"
)

type BlogPosts []BlogPost

func (blogs BlogPosts) SortDescending() {
	slices.SortFunc(blogs, func(a, b BlogPost) int {
		return cmp.Compare(b.DatePosted(), a.DatePosted())
	})
}

func (blogs *BlogPosts) Append(blog BlogPost) {
	// https://stackoverflow.com/a/24726368/446681
	*blogs = append(*blogs, blog)
}

func (blogs BlogPosts) Content() string {
	content := "# Blog Posts\r\n"
	blogs.SortDescending()
	year := 0
	for _, blog := range blogs {
		if blog.YearPosted() != year {
			content += fmt.Sprintf("## %d\r\n", blog.YearPosted())
			year = blog.YearPosted()
		}
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

	rss := NewRss(meta.Title, meta.Description, meta.BaseLink)
	for _, blog := range blogs {
		rss.Add(blog.Title(), blog.Summary(), blog.LinkUrl(), blog.DatePosted())
	}
	xml, err := rss.ToXml()
	if err != nil {
		fmt.Printf("ERROR producing RSS file: %s\r\n", err)
	}
	saveFile(filename, xml)
}
