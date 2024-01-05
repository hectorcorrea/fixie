package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

type BlogPost struct {
	Filename string
	Content  string
	Metadata Metadata
}

type Metadata struct {
	Title       string  `xml:"title"`     // Not used
	Slug        string  `xml:"slug"`      // Used when calculating redirects for legacy posts
	Summary     string  `xml:"summary"`   // Not used
	CreatedOn   string  `xml:"createdOn"` // Used when file name does not have a date
	UpdatedOn   string  `xml:"updatedOn"` // Not used
	PostedOn    string  `xml:"postedOn"`  // Not used
	OldSequence string  `xml:"oldSequence"`
	Fields      []Field `xml:"fields"` // Used when calculating redirects for legacy posts
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
	if len(b.Metadata.CreatedOn) >= 10 {
		// Use the date part (YYYY-MM-DD) from the metadata
		return b.Metadata.CreatedOn[0:10]
	}

	date := dateFromFilename(b.Filename)
	if date != "" {
		// Use the data part from the filename
		return date
	}

	return "1970-01-01"
}

func (b BlogPost) DatePosted() string {
	if len(b.Metadata.PostedOn) >= 10 {
		// Use the date part (YYYY-MM-DD) from the metadata
		return b.Metadata.PostedOn[0:10]
	}

	date := dateFromFilename(b.Filename)
	if date != "" {
		// Use the data part from the filename
		return date
	}

	return "1970-01-01"
}

func (b BlogPost) YearCreated() int {
	yearString := b.DateCreated()[0:4]
	year, _ := strconv.Atoi(yearString)
	return year
}

func (b BlogPost) YearPosted() int {
	yearString := b.DatePosted()[0:4]
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

func (b BlogPost) ToSearchEntry() string {
	entry := fmt.Sprintf(`{
		"id": "%s",
		"name": "%s",
		"text": "%s"
	}`, b.LinkUrl(), b.Title(), b.SearchText())
	return entry
}

func (b BlogPost) SearchText() string {
	// Extract all alphanumeric words and a few selected special characters
	plainText := ""
	var prevChar rune
	for _, c := range strings.ToLower(b.Content) {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			plainText += string(c)
		} else if c == '.' || c == '/' || c == '@' || c == ':' || c == '-' {
			plainText += string(c)
		} else if c == '#' && (prevChar == 'C' || prevChar == 'c') {
			plainText += string(c)
		} else {
			plainText += " "
		}
		prevChar = c
	}
	// Remove duplicates
	tokens := []string{}
	for _, token := range strings.Split(plainText, " ") {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		if !slices.Contains(tokens, token) {
			tokens = append(tokens, token)
		}
	}
	return strings.Join(tokens, " ")
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
	for _, field := range b.Metadata.Fields {
		if field.Name == "oldId" {
			oldId, _ := strconv.Atoi(field.Value)
			return oldId
		}
	}
	return 0
}

// Creates the redirect files required to support legacy URLs indicated in the metadata file
// Redirect legacy URLs
//
//	./blog/slug/index.html
//	./blog/slug/10									(old id)
//	./blog/slug/2008-11-25-00001		(date created + sequence)
//
// to new format
//
//	./blog/2008-11-30/slug					(date posted)
func (b BlogPost) createRedirectFiles() bool {
	newUrl := fmt.Sprintf("/blog/%s/%s", b.DatePosted(), b.Metadata.Slug)
	content := `<head><meta http-equiv="Refresh" content="0; URL=URL-GOES-HERE" /></head>`
	content = strings.Replace(content, "URL-GOES-HERE", newUrl, 1)

	redirectFile1 := fmt.Sprintf("./blog/%s/index.html", b.Metadata.Slug)
	redirectFolder := fmt.Sprintf("./blog/%s", b.Metadata.Slug)
	createDir(redirectFolder)
	saveFile(redirectFile1, content)

	if b.OldId() != 0 {
		redirectFile2 := fmt.Sprintf("./blog/%s/%d.html", b.Metadata.Slug, b.OldId())
		saveFile(redirectFile2, content)
	}

	if b.Metadata.OldSequence != "" {
		redirectFile3 := fmt.Sprintf("./blog/%s/%s-%s.html", b.Metadata.Slug, b.DateCreated(), b.Metadata.OldSequence)
		saveFile(redirectFile3, content)
	}
	return true
}
