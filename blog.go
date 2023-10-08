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

func SortDescending(blogPosts []BlogPost) []BlogPost {
	slices.SortFunc(blogPosts, func(a, b BlogPost) int {
		return cmp.Compare(b.DateCreated(), a.DateCreated())
	})
	return blogPosts
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
