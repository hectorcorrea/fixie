package main

import (
	"cmp"
	"fmt"
	"slices"
)

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
	year := 0
	for _, blog := range blogs {
		if blog.YearCreated() != year {
			content += fmt.Sprintf("## %d\r\n", blog.YearCreated())
			year = blog.YearCreated()
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
