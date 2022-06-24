package sortfile

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
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

	file, err := os.OpenFile(fl.Arg(0), os.O_RDWR, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s: %v\n", fl.Arg(0), err)
		return err
	}
	app.file = file

	return nil
}

func (app *appEnv) run() error {
	defer app.file.Close()

	return app.sort()
}

func (app *appEnv) sort() error {
	scanner := bufio.NewScanner(app.file)
	file := make([]string, 0)

	for scanner.Scan() {
		file = append(file, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	sort.Strings(file)
	var b bytes.Buffer
	err := app.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = app.file.Seek(0, 0)
	if err != nil {
		return err
	}

	for _, v := range file {
		b.Write([]byte(v))
		b.Write([]byte("\n"))
	}

	_, err = app.file.Write(b.Bytes())
	if err != nil {
		return err
	}

	err = app.file.Sync()
	if err != nil {
		return err
	}

	return nil
}
