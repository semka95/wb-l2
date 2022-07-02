package grep

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
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
	after        int
	before       int
	context      int
	count        int
	ignoreCase   bool
	invert       bool
	isFixed      bool
	printLineNum bool
	reader       io.Reader
	input        []string
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("go-grep", flag.ContinueOnError)
	fl.IntVar(
		&app.after, "A", 0, "Print NUM lines of trailing context after matching lines.",
	)
	fl.IntVar(
		&app.before, "B", 0, "Print NUM lines of leading context before matching lines.",
	)
	fl.IntVar(
		&app.context, "C", 0, "Print NUM lines of output context.",
	)
	fl.IntVar(
		&app.count, "c", -1, "Suppress normal output; instead print a count of matching lines for each input file.  With the -v, count non-matching lines.",
	)
	fl.BoolVar(
		&app.ignoreCase, "i", false, " Ignore case distinctions in patterns and input data, so that characters that differ only in case match each other.",
	)
	fl.BoolVar(
		&app.invert, "v", false, " Invert the sense of matching, to select non-matching lines.",
	)
	fl.BoolVar(
		&app.isFixed, "F", false, " Interpret PATTERNS as fixed strings, not regular expressions.",
	)
	fl.BoolVar(
		&app.printLineNum, "n", false, "Prefix each line of output with the 1-based line number within its input file.",
	)

	if err := fl.Parse(args); err != nil {
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
	scanner := bufio.NewScanner(app.reader)
	for scanner.Scan() {
		app.input = append(app.input, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
