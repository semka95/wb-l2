package sortfile

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

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
	file            *os.File
	data            []string
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("sortfile", flag.ContinueOnError)
	fl.IntVar(&app.column, "k", 1, "sort via a column")
	fl.BoolVar(&app.isNumeric, "n", false, "compare according to string numerical value")
	fl.BoolVar(&app.isReverse, "r", false, "reverse the result of comparisons")
	fl.BoolVar(&app.deleteDuplicate, "u", false, "delete duplicate")

	if err := fl.Parse(args); err != nil {
		fl.Usage()
		return err
	}

	file, err := os.Open(fl.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s: %v\n", fl.Arg(0), err)
		return err
	}
	app.file = file

	return nil
}

func (app *appEnv) run() error {
	defer app.file.Close()

	scanner := bufio.NewScanner(app.file)
	for scanner.Scan() {
		app.data = append(app.data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if app.column == 1 {
		app.sort()
	}

	if app.column > 1 {
		app.sortColumns()
	}

	app.writeToOutput()

	return nil
}

func (app *appEnv) sort() {
	sort.Strings(app.data)
}

type table struct {
	data   [][]string
	column int
}

func (t table) Len() int {
	return len(t.data)
}

func (t table) Less(i, j int) bool {
	return (t.data[i][t.column] < t.data[j][t.column])
}

func (t table) Swap(i, j int) {
	t.data[i], t.data[j] = t.data[j], t.data[i]
}

func (app *appEnv) sortColumns() {
	t := table{
		data:   make([][]string, 0, len(app.data)),
		column: app.column - 1,
	}
	for _, v := range app.data {
		t.data = append(t.data, strings.Fields(v))
	}

	sort.Sort(t)

	for i, v := range t.data {
		app.data[i] = strings.Join(v, " ")
	}
}

func (app *appEnv) writeToOutput() {
	for _, v := range app.data {
		fmt.Fprintf(os.Stdout, "%s\n", v)
	}
}
