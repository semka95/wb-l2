package wget

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	link "go-wget/parse"
)

// Sitemap represents data needed to build sitemap using BuildSitemap
type Sitemap struct {
	rootLink     string
	visitedLinks map[string]struct{}
	directory    string
}

// NewSitemap creates instance of Sitemap
func NewSitemap(rootLink string, directory string) Sitemap {
	v := map[string]struct{}{
		rootLink: {},
	}

	s := Sitemap{
		rootLink:     rootLink,
		visitedLinks: v,
		directory:    directory,
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
		defer resp.Body.Close()

		mediatype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
		if err != nil {
			fmt.Printf("can't parse link type: %s", err.Error())
		}
		ext, err := mime.ExtensionsByType(mediatype)
		if err != nil || len(ext) == 0 {
			fmt.Printf("can't parse link type: %s", err.Error())
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		r := bytes.NewReader(body)
		fileName := path.Join(s.directory, path.Base(resp.Request.URL.Path)+ext[0])
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		size, err := io.Copy(file, r)
		if err != nil {
			return err
		}

		fmt.Printf("\nDownloaded a file %s with size %d bytes\n", fileName, size)

		r.Seek(0, 0)
		l, err := s.parseLinks(r)
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
		href := v.Href.Host + v.Href.Path
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
