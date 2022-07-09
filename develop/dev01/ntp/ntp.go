package ntp

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

// CLI runs the go-ntp command line app and returns its exit status.
func CLI(args []string) int {
	var app appEnv
	err := app.fromArgs(args)
	if err != nil {
		return 2
	}

	if err := app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}
	return 0
}

type appEnv struct {
	host string
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("ntp", flag.ContinueOnError)
	fl.StringVar(&app.host, "host", "0.beevik-ntp.pool.ntp.org", "ntp host")

	if err := fl.Parse(args); err != nil {
		return err
	}

	return nil
}

func (app *appEnv) run() error {
	t, err := ntp.Time(app.host)
	if err != nil {
		return err
	}

	fmt.Println(t.UTC().Format(time.UnixDate))
	return nil
}
