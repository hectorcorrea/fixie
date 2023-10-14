package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type BlogPost struct {
	Filename string
	Content  string
	Metadata Metadata
}

type Metadata struct {
	Title     string  `xml:"title"`     // Not used
	Slug      string  `xml:"slug"`      // Used when calculating redirects for legacy posts
	Summary   string  `xml:"summary"`   // Not used
	CreatedOn string  `xml:"createdOn"` // Used when file name does not have a date
	UpdatedOn string  `xml:"updatedOn"` // Not used
	PostedOn  string  `xml:"postedOn"`  // Not used
	Fields    []Field `xml:"fields"`    // Used when calculating redirects for legacy posts
}

type Field struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

func LoadBlogPost(filename string) BlogPost {
	content := readFile(filename)

	blog := BlogPost{Filename: filename, Content: content}
	blog.Metadata = blog.fetchMetadata()
	return blog
}

func (b BlogPost) DateCreated() string {
	date := dateFromFilename(b.Filename)
	if date != "" {
		return date
	}

	if len(b.Metadata.CreatedOn) >= 10 {
		// Use the date part (YYYY-MM-DD) from the metadata
		return b.Metadata.CreatedOn[0:10]
	}

	return "1970-01-01"
}

func (b BlogPost) YearCreated() int {
	yearString := b.DateCreated()[0:4]
	year, _ := strconv.Atoi(yearString)
	return year
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

func (b BlogPost) MetadataFile() string {
	return strings.TrimSuffix(b.Filename, ".md") + ".xml"
}

func (b BlogPost) fetchMetadata() Metadata {
	reader, err := os.Open(b.MetadataFile())
	if err != nil {
		return Metadata{}
	}
	defer reader.Close()

	byteValue, err := io.ReadAll(reader)
	var metadata Metadata
	xml.Unmarshal(byteValue, &metadata)
	return metadata
}

func (b BlogPost) OldId() int {
	oldId := 0
	for _, field := range b.Metadata.Fields {
		if field.Name == "oldId" {
			oldId, _ = strconv.Atoi(field.Value)
		}
	}
	return oldId
}

func (b BlogPost) SequenceNumber() string {
	return sequenceFromString(filepath.Dir(b.Filename))
}

// Creates the redirect files required to support legacy URLs indicated in the metadata file
func (b BlogPost) createRedirectFiles() bool {
	oldId := b.OldId()
	if oldId == 0 {
		return false
	}

	redirectFolder := fmt.Sprintf("./blog/%s", b.Metadata.Slug)
	redirectFile1 := fmt.Sprintf("%s/%d.html", redirectFolder, oldId)
	redirectFile2 := fmt.Sprintf("%s/index.html", redirectFolder)

	content := `<head><meta http-equiv="Refresh" content="0; URL=URL-GOES-HERE" /></head>`

	sequenceNumber := b.SequenceNumber()
	folder := b.DateCreated()
	if sequenceNumber != "" {
		folder += "-" + sequenceNumber
	}
	newUrl := fmt.Sprintf("/blog/%s/%s", folder, b.Metadata.Slug)
	content = strings.Replace(content, "URL-GOES-HERE", newUrl, 1)

	createDir(redirectFolder)
	saveFile(redirectFile1, content)
	saveFile(redirectFile2, content)
	return true
}
