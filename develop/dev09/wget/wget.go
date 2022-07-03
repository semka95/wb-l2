package wget

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	link "go-wget/parse"
)

// Sitemap represents data needed to build sitemap using BuildSitemap
type Sitemap struct {
	rootLink     string
	visitedLinks map[string]struct{}
}

// NewSitemap creates instance of Sitemap
func NewSitemap(rootLink string) Sitemap {
	v := map[string]struct{}{
		rootLink: {},
	}

	s := Sitemap{
		rootLink:     rootLink,
		visitedLinks: v,
	}
	return s
}

// BuildSitemap visits links in queue and parses links in other queue and
// recursively visits them. If depth is greater than zero it restricts number
// of recursive calls. All visited links recorded to Sitemap.visitedLinks
func (s *Sitemap) BuildSitemap(queue []string, depth int) error {
	if depth == 0 {
		return nil
	}

	discoveredLinks := make([]string, 0)

	for _, v := range queue {
		resp, err := http.Get(v)
		if err != nil {
			return err
		}
		fmt.Println(path.Base(resp.Request.URL.Path))

		l, err := s.parseLinks(resp.Body)
		if err != nil {
			return err
		}
		discoveredLinks = append(discoveredLinks, l...)
	}

	if len(discoveredLinks) > 0 {
		depth--
		return s.BuildSitemap(discoveredLinks, depth)
	}

	return nil
}

// parseLinks reads html data from io.Reader and creates array of links.
// It only parses links with Sitemap.rootLink domain
func (s *Sitemap) parseLinks(r io.Reader) ([]string, error) {
	res, err := link.ParseHTML(r)
	if err != nil {
		return nil, err
	}

	links := make([]string, 0)

	for _, v := range res {
		href := v.Href
		visited := true

		switch {
		case strings.HasPrefix(href, s.rootLink):
			visited = s.isVisited(href)
		case strings.HasPrefix(href, "/"):
			href = s.rootLink + href
			visited = s.isVisited(href)
		}

		if !visited {
			links = append(links, href)
		}
	}

	return links, nil
}

// isVisited checks if url visited
func (s *Sitemap) isVisited(href string) (visited bool) {
	if _, visited = s.visitedLinks[href]; !visited {
		s.visitedLinks[href] = struct{}{}
	}

	return visited
}
