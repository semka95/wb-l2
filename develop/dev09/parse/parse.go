// link is a package for parsing HTML link tags (<a href="..."...</a>),
package link

import (
	"io"
	"net/url"

	"golang.org/x/net/html"
)

// Link represents HTML link tag
type Link struct {
	Href *url.URL
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
			u, err := parseHref(n.Attr)
			if err != nil {
				return
			}
			link := Link{
				Href: u,
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

// parseHref extracts href attribute from link tag
func parseHref(attrs []html.Attribute) (*url.URL, error) {
	var href string

	for _, a := range attrs {
		if a.Key == "href" {
			href = a.Val
			break
		}
	}
	u, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	return u, nil
}
