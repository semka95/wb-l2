package link

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHtml(t *testing.T) {
	u1, _ := url.Parse("/other-page")
	u2, _ := url.Parse("https://www.twitter.com/joncalhoun")
	u3, _ := url.Parse("https://github.com/gophercises")
	u4, _ := url.Parse("#")
	u5, _ := url.Parse("/lost")
	u6, _ := url.Parse("https://twitter.com/marcusolsson")
	u7, _ := url.Parse("/dog-cat")
	cases := []struct {
		HtmlPath string
		Links    []Link
	}{
		{
			"./testdata/ex1.html",
			[]Link{
				{Href: u1},
			},
		},
		{
			"./testdata/ex2.html",
			[]Link{
				{Href: u2},
				{Href: u3},
			},
		},
		{
			"./testdata/ex3.html",
			[]Link{
				{Href: u4},
				{Href: u5},
				{Href: u6},
			},
		},
		{
			"./testdata/ex4.html",
			[]Link{
				{Href: u7},
			},
		},
	}

	for _, test := range cases {
		t.Run(fmt.Sprintf("%s file", test.HtmlPath), func(t *testing.T) {
			f, err := os.Open(test.HtmlPath)
			assert.NoError(t, err)
			defer f.Close()

			got, err := ParseHTML(f)
			assert.NoError(t, err)

			assert.EqualValues(t, test.Links, got)
		})
	}
}

func ExampleParseHTML() {
	htmlFile := `<html><body><h1>Hello!</h1><a href="/other-page">A link to another page</a></body></html>`

	r := strings.NewReader(htmlFile)
	result, err := ParseHTML(r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", result)
	// Output: [{Href:/other-page}]
}
