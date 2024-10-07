package main

import "regexp"

type SiteMeta struct {
	Title       string
	Author      string
	Description string
	BaseLink    string
}

// Extract site metadata from an HTML page
func NewSiteMeta(html string) SiteMeta {
	meta := SiteMeta{
		Title:       htmlTitle(html),
		Description: htmlMetaAttr(html, "description"),
		Author:      htmlMetaAttr(html, "author"),
		BaseLink:    htmlBaseLink(html),
	}
	return meta
}

// Extracts the value of a meta attribute in an HTML page
func htmlMetaAttr(html string, metaName string) string {
	re := regexp.MustCompile(`<meta name="` + metaName + `" content="(.*)"`)
	return regExpMatch(html, re)
}

// Extracts the title of an HTML page
func htmlTitle(html string) string {
	re := regexp.MustCompile(`<title>(.*)</title>`)
	return regExpMatch(html, re)
}

// Extracts the base link value of an HTML page
func htmlBaseLink(html string) string {
	re := regexp.MustCompile(`<base href="(.*)"`)
	return regExpMatch(html, re)
}
