package wget

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/schollz/progressbar/v3"
)

// appEnv represents parsed command line arguments
type appEnv struct {
	link       *url.URL
	outputFile string
	depth      int
	recursive  bool
	resources  bool
}

// CLI runs the go-wget command line app and returns its exit status
func CLI(args []string) int {
	var app appEnv
	err := app.fromArgs(args)
	if err != nil {
		return 2
	}
	if err = app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}
	return 0
}

// fromArgs parses command line arguments into appEnv struct
func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("go-wget", flag.ContinueOnError)
	fl.StringVar(&app.outputFile, "O", "", "Path to output file")
	fl.IntVar(&app.depth, "l", -1, "Maximum number of links to follow when building downloading the site. By default depth is not set")
	fl.BoolVar(&app.recursive, "r", false, "Turn on recursive retriving")
	fl.BoolVar(&app.resources, "p", false, "Download all the files that are necessary to properly display a given HTML page")

	if err := fl.Parse(args); err != nil {
		return err
	}

	u, err := url.Parse(fl.Arg(0))
	if err != nil {
		return err
	}
	app.link = u
	app.depth++

	if app.outputFile == "" {
		app.outputFile = path.Base(app.link.Path)
	}

	return nil
}

func (app *appEnv) run() error {
	if app.recursive {
		queue := []string{app.link.String()}
		if err := os.Mkdir(app.link.Host, os.ModePerm); err != nil {
			return err
		}
		sm := NewSite(app.link.String(), app.link.Host)
		err := sm.DownloadSite(queue, app.depth)
		if err != nil {
			return err
		}
		return nil
	}

	return downloadFile(app.link.String(), app.outputFile)
}

func downloadFile(url string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Opaque = req.URL.Path
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	size, err := io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("\nDownloaded a file %s with size %d bytes\n", filePath, size)

	return nil
}
