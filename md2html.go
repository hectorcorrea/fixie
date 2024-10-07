package main

import (
	"regexp"
	"strings"
)

var reBold, reItalic, reStrike, reLink, reCode, reImg, reImgHtml *regexp.Regexp

type MarkdownParser struct {
	markdown     string   // original Markdown
	html         string   // resulting HTML
	title        string   // internal
	descriptions []string // internal
	pre          bool     // internal
	quote        bool     // internal
	li           bool     // internal
}

func init() {
	// text **in bold**
	reBold = regexp.MustCompile("(\\*\\*)(.*?)(\\*\\*)")

	// text *in italic*
	reItalic = regexp.MustCompile("(\\*)(.*?)(\\*)")

	// ~~striked text~~
	reStrike = regexp.MustCompile("(~~)(.*?)(~~)")

	// [some text](http://somewhere.org)
	// \\[							starts with [
	//    ([^\\]]*?)		followed by any character except ] using lazy match
	// \\]							followed by an ending ]
	// \\(							followed by (
	//    (.*?)					followed by any character using lazy match
	// \\)      				ends with )
	//
	// Using ([^\\]]*?) instead of (.*?) to prevent the parser from
	// picking up text in brackets that is not an URL, for example:
	// [1] [book x](http://link/to/bookx)
	reLink = regexp.MustCompile("\\[([^\\]]*?)\\]\\((.*?)\\)")

	// `hello world`
	reCode = regexp.MustCompile("(`)(.*?)(`)")

	// ![image caption](http://somewhere.org/image.png)
	reImg = regexp.MustCompile("!\\[([^\\]]*?)\\]\\((.*?)\\)")

	// <img ... />
	reImgHtml = regexp.MustCompile(`<img (.*) />`)
}

func NewMarkdownParser(markdown string) MarkdownParser {
	parser := MarkdownParser{markdown: markdown}
	parser.processMarkdown()
	return parser
}

func (p *MarkdownParser) processMarkdown() {
	p.html = ""
	p.pre = false
	p.quote = false
	p.li = false
	p.descriptions = []string{}
	p.title = ""
	lines := strings.Split(p.markdown, "\n")
	for _, line := range lines {
		if p.isQuote(line) {
			if p.quote {
				// already in blockquote
			} else {
				// start a new blockquote
				p.html += "<blockquote>\n"
				p.quote = true
			}
		} else {
			if p.quote {
				// end current blockquote
				p.html += "</blockquote>\n"
				p.quote = false
			}
		}

		if p.isListItem(line) {
			if p.li {
				// already inside a list
			} else {
				// start a new list
				p.html += "<ul>\n"
				p.li = true
			}
		} else {
			if p.li {
				// end current list
				p.html += "</ul>\n"
				p.li = false
			}
		}

		l := strings.TrimSpace(line)

		if p.isH1(l) {
			p.html += "<h1>" + substr(l, 2) + "</h1>\n"
			if p.title == "" {
				p.title = chomp(strings.TrimPrefix(l, "# "))
			}
		} else if p.isH2(l) {
			p.html += "<h2>" + substr(l, 3) + "</h2>\n"
		} else if p.isH3(l) {
			p.html += "<h3>" + substr(l, 4) + "</h3>\n"
		} else if p.isPreTerminal(l) {
			p.html += "<pre class=\"terminal\">\n"
			p.pre = true
		} else if p.isPreCode(l) {
			p.html += "<pre class=\"code\">\n"
			p.pre = true
		} else if p.isPre(l) {
			if p.pre {
				p.html += "</pre>\n"
				p.pre = false
			} else {
				p.html += "<pre>\n"
				p.pre = true
			}
		} else if l == "" {
			// html += "<br/>\n"
			p.html += "\n"
		} else {
			if p.pre {
				// we use the original line in pre to preserve spaces
				p.html += p.parsedLine(line, true) + "\n"
			} else if p.quote {
				p.html += p.parsedLine(substr(l, 2), false) + "<br/>\n"
			} else if p.li {
				p.html += "<li>" + p.parsedLine(substr(l, 2), false) + "\n"
			} else {
				// A regular paragraph
				p.html += "<p>" + p.parsedLine(l, false) + "</p>\n"

				// Aggregate the first few paragraps to use as the description
				if len(p.descriptions) < 3 {
					p.descriptions = append(p.descriptions, p.cleanLine(l))
				}
			}
		}
	}
}

func (p MarkdownParser) Html() string {
	return p.html
}

func (p MarkdownParser) Title() string {
	return p.title
}

func (p MarkdownParser) Description() string {
	return strings.Join(p.descriptions, " ")
}

func (p MarkdownParser) isH1(line string) bool {
	if p.pre {
		return false
	}
	return strings.HasPrefix(line, "# ")
}

func (p MarkdownParser) isH2(line string) bool {
	if p.pre {
		return false
	}
	return strings.HasPrefix(line, "## ")
}

func (p MarkdownParser) isH3(line string) bool {
	if p.pre {
		return false
	}
	return strings.HasPrefix(line, "### ")
}

func substr(line string, i int) string {
	if i >= len(line) {
		return ""
	}
	return line[i:]
}

func chomp(text string) string {
	text = strings.TrimSuffix(text, "\r\n")
	text = strings.TrimSuffix(text, "\r")
	text = strings.TrimSuffix(text, "\n")
	return text
}

func (p MarkdownParser) isPreTerminal(line string) bool {
	return strings.HasPrefix(line, "```terminal")
}

func (p MarkdownParser) isPreCode(line string) bool {
	return strings.HasPrefix(line, "```code")
}

func (p MarkdownParser) isPre(line string) bool {
	return strings.HasPrefix(line, "```")
}

func (p MarkdownParser) isQuote(line string) bool {
	return strings.HasPrefix(line, "> ") || line == ">"
}

func (p MarkdownParser) isListItem(line string) bool {
	return strings.HasPrefix(line, "* ")
}

func (p MarkdownParser) cleanLine(line string) string {
	// Ditch HTML image tags (although we allow them in the Markdown
	// we don't want it on the "clean lines")
	line = reImgHtml.ReplaceAllString(line, " ")
	line = strings.Replace(line, "<", "&lt;", -1)
	line = strings.Replace(line, ">", "&gt;", -1)
	line = strings.Replace(line, "\"", " ", -1)
	return line
}

func (p MarkdownParser) parsedLine(line string, pre bool) string {
	// Encode the < and > characters in the original Markdown...
	line = strings.Replace(line, "<", "&lt;", -1)
	line = strings.Replace(line, ">", "&gt;", -1)

	// and then allow for <img ... />
	line = strings.Replace(line, "&lt;img ", "<img ", -1)
	line = strings.Replace(line, "/&gt;", "/>", -1)

	// allow for <sup> </sup>
	line = strings.Replace(line, "&lt;sup&gt;", "<sup>", -1)
	line = strings.Replace(line, "&lt;/sup&gt;", "</sup>", -1)

	if pre {
		// don't do any extra processing if we are on code block
		return line
	}

	// and finally convert the Markdown to HTML
	line = reImg.ReplaceAllString(line, "<img src=\"$2\" alt=\"$1\" title=\"$1\" />")
	line = reBold.ReplaceAllString(line, "<b>$2</b>")
	line = reItalic.ReplaceAllString(line, "<i>$2</i>")
	line = reStrike.ReplaceAllString(line, "<strike>$2</strike>")
	line = reLink.ReplaceAllString(line, "<a href=\"$2\">$1</a>")
	line = reCode.ReplaceAllString(line, "<code>$2</code>")
	return line
}
