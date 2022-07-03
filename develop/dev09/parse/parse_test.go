package link

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHtml(t *testing.T) {
	cases := []struct {
		HtmlPath string
		Links    []Link
	}{
		{
			"./testfiles/ex1.html",
			[]Link{
				{Href: "/other-page", Text: "A link to another page"},
			},
		},
		{
			"./testfiles/ex2.html",
			[]Link{
				{Href: "https://www.twitter.com/joncalhoun", Text: `Check me out on twitter`},
				{Href: "https://github.com/gophercises", Text: `Gophercises is on Github!`},
			},
		},
		{
			"./testfiles/ex3.html",
			[]Link{
				{Href: "#", Text: "Login"},
				{Href: "/lost", Text: "Lost? Need help?"},
				{Href: "https://twitter.com/marcusolsson", Text: "@marcusolsson"},
			},
		},
		{
			"./testfiles/ex4.html",
			[]Link{
				{Href: "/dog-cat", Text: "dog cat"},
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
	// Output: [{Href:/other-page Text:A link to another page}]
}
