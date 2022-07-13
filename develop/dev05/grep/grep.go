package grep

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
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
	reader       io.ReadCloser
	input        []string
	pattern      string
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("go-grep", flag.ContinueOnError)
	fl.IntVar(&app.after, "A", 0, "Print NUM lines of trailing context after matching lines.")
	fl.IntVar(&app.before, "B", 0, "Print NUM lines of leading context before matching lines.")
	fl.IntVar(&app.context, "C", 0, "Print NUM lines of output context.")
	fl.IntVar(&app.count, "c", -1, "Suppress normal output; instead print a count of matching lines for each input file.  With the -v, count non-matching lines.")
	fl.BoolVar(&app.ignoreCase, "i", false, " Ignore case distinctions in patterns and input data, so that characters that differ only in case match each other.")
	fl.BoolVar(&app.invert, "v", false, " Invert the sense of matching, to select non-matching lines.")
	fl.BoolVar(&app.isFixed, "F", false, " Interpret PATTERNS as fixed strings, not regular expressions.")
	fl.BoolVar(&app.printLineNum, "n", false, "Prefix each line of output with the 1-based line number within its input file.")

	if err := fl.Parse(args); err != nil {
		fl.Usage()
		return err
	}

	if app.after == 0 {
		app.after = app.context
	}
	if app.before == 0 {
		app.before = app.context
	}

	app.pattern = fl.Arg(0)
	if app.isFixed {
		app.pattern = fmt.Sprintf("^%s$", app.pattern)
	}
	if app.ignoreCase {
		app.pattern = "(?i)" + app.pattern
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		app.reader = os.Stdin
		return nil
	}

	file, err := os.Open(fl.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s: %v\n", fl.Arg(0), err)
		return err
	}
	app.reader = file

	return nil
}

func (app *appEnv) run() error {
	defer app.reader.Close()
	scanner := bufio.NewScanner(app.reader)
	for scanner.Scan() {
		app.input = append(app.input, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	r, err := regexp.Compile(app.pattern)
	if err != nil {
		return err
	}

	app.printResult(app.findMatched(r))

	return nil
}

func (app *appEnv) findMatched(r *regexp.Regexp) []int {
	matched := make([]int, 0, len(app.input)/3)

	for i, v := range app.input {
		if app.count == 0 {
			break
		}
		if r.MatchString(v) && !app.invert || !r.MatchString(v) && app.invert {
			matched = append(matched, i)
			app.count--
		}
	}

	return matched
}

func (app *appEnv) printResult(matched []int) {
	printed := make(map[int]struct{})
	matchedM := make(map[int]struct{})
	for _, v := range matched {
		matchedM[v] = struct{}{}
	}

	var lastPrinted int
	for _, v := range matched {
		if app.before > 0 || app.after > 0 {
			if v-lastPrinted > 2 {
				fmt.Println("--")
			}
			if _, ok := printed[v]; ok {
				continue
			}

			start := v - app.before
			if v-app.before < 0 {
				start = 0
			}
			finish := v + app.after
			if v+app.after > len(app.input)-1 {
				finish = len(app.input) - 1
			}

			for ; start <= finish; start++ {
				if _, ok := printed[start]; ok {
					continue
				}
				if _, ok := matchedM[start]; ok && start != v {
					break
				}
				if app.printLineNum {
					if _, ok := matchedM[start]; ok {
						fmt.Printf("%d:%s\n", start+1, app.input[start])
					} else {
						fmt.Printf("%d-%s\n", start+1, app.input[start])
					}
					printed[start] = struct{}{}
					lastPrinted = start
					continue
				}
				fmt.Println(app.input[start])
				printed[start] = struct{}{}
				lastPrinted = start

			}
			continue
		}

		if app.printLineNum {
			fmt.Printf("%d:%s\n", v+1, app.input[v])
			continue
		}
		fmt.Println(v)
	}
}
