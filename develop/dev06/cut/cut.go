package cut

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
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
	fields    []int
	delimiter string
	separated bool
	toEnd     bool
	reader    io.Reader
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("go-cut", flag.ContinueOnError)
	fl.Func("f", "select only these fields;  also print any line that contains no delimiter character, unless the -s option is specified", app.parseFields)
	fl.StringVar(&app.delimiter, "d", "\t", "use DELIM instead of TAB for field delimiter")
	fl.BoolVar(&app.separated, "s", false, "do not print lines not containing delimiters")

	if err := fl.Parse(args); err != nil {
		fl.Usage()
		return err
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		app.reader = os.Stdin
		return nil
	}
	return nil
}

func (a *appEnv) parseFields(s string) error {
	if strings.Contains(s, ",") {
		f := strings.Split(s, ",")
		for _, v := range f {
			i, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			a.fields = append(a.fields, i-1)
		}
		return nil
	}

	if strings.Contains(s, "-") {
		var from, to int
		var err error
		if len(s) == 3 {
			from, err = strconv.Atoi(string(s[0]))
			if err != nil {
				return err
			}
			to, err = strconv.Atoi(string(s[2]))
			if err != nil {
				return err
			}
		}

		if len(s) == 2 {
			if s[0] == '-' {
				from = 1
				to, err = strconv.Atoi(string(s[1]))
				if err != nil {
					return err
				}
			}

			if s[1] == '-' {
				a.toEnd = true
				from, err = strconv.Atoi(string(s[0]))
				if err != nil {
					return err
				}
				to = from
			}
		}

		for ; from <= to; from++ {
			a.fields = append(a.fields, from-1)
		}

		return nil
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	a.fields = append(a.fields, n-1)

	return nil
}

func (app *appEnv) run() error {
	scanner := bufio.NewScanner(app.reader)
	for scanner.Scan() {
		res := strings.Split(scanner.Text(), app.delimiter)
		l := len(res)

		if l == 1 && app.separated {
			continue
		}

		if l == 1 && !app.separated {
			fmt.Println(res[0])
			continue
		}

		if app.toEnd {
			for i := app.fields[0]; i < l; i++ {
				fmt.Printf("%s%s", res[i], app.delimiter)
			}
			fmt.Println()
			continue
		}

		for _, v := range app.fields {
			if v < l {
				fmt.Printf("%s%s", res[v], app.delimiter)
			}
		}
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
