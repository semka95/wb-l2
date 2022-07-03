// link is a package for parsing HTML link tags (<a href="..."...</a>),
package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents HTML link tag
type Link struct {
	Href string
	Text string
}

// ParseHTML parses given html file and returns slice of links
func ParseHTML(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := make([]Link, 0)

	var parseNode func(node *html.Node)
	parseNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			link := Link{
				Href: parseHref(n.Attr),
				Text: parseLinkText(n),
			}

			links = append(links, link)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c)
		}
	}

	parseNode(doc)

	return links, nil
}

// parseLinkText extracts text from link tag
func parseLinkText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type != html.ElementNode {
		return ""
	}

	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += parseLinkText(c)
	}

	return strings.Join(strings.Fields(text), " ")
}

// parseHref extracts href attribute from link tag
func parseHref(attrs []html.Attribute) string {
	var href string

	for _, a := range attrs {
		if a.Key == "href" {
			href = a.Val
			break
		}
	}

	return href
}
