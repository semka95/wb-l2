package sortstrings

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// CLI runs the go-sort command line app and returns its exit status.
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

type appEnv struct {
	isNumeric       bool
	isReverse       bool
	deleteDuplicate bool
	column          int
	reader          io.ReadCloser
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("sortfile", flag.ContinueOnError)
	fl.IntVar(&app.column, "k", 1, "sort via a column")
	fl.BoolVar(&app.isNumeric, "n", false, "compare according to string numerical value")
	fl.BoolVar(&app.isReverse, "r", false, "reverse the result of comparisons")
	fl.BoolVar(&app.deleteDuplicate, "u", false, "delete duplicate strings")

	if err := fl.Parse(args); err != nil {
		fl.Usage()
		return err
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		app.reader = os.Stdin
		return nil
	}

	file, err := os.Open(fl.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s: %v\n", fl.Arg(0), err)
		return err
	}
	app.reader = file

	return nil
}

func (app *appEnv) run() error {
	defer app.reader.Close()
	data := make([]string, 0)

	scanner := bufio.NewScanner(app.reader)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if app.column == 1 && !app.isNumeric {
		data = app.sort(data)
		writeToOutput(data)
		return nil
	}

	data = app.sortColumns(data)
	writeToOutput(data)

	return nil
}

func (app *appEnv) sort(data []string) []string {
	if app.isReverse {
		sort.Sort(sort.Reverse(sort.StringSlice(data)))
	} else {
		sort.Strings(data)
	}

	if app.deleteDuplicate {
		data = delDuplicate(data)
	}

	return data
}

func (app *appEnv) sortColumns(data []string) []string {
	t := stringTable{
		data:      make([][]string, 0, len(data)),
		column:    app.column - 1,
		isNumeric: app.isNumeric,
	}

	for _, v := range data {
		t.data = append(t.data, strings.Fields(v))
	}

	if app.isReverse {
		sort.Sort(sort.Reverse(t))
	} else {
		sort.Sort(t)
	}

	for i, v := range t.data {
		data[i] = strings.Join(v, " ")
	}

	if app.deleteDuplicate {
		data = delDuplicate(data)
	}

	return data
}
