package unpacker

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// CLI runs the go-grab-xkcd command line app and returns its exit status.
func CLI(args []string) int {
	var app appEnv
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		err := app.fromStdin()
		if err != nil {
			return 2
		}
	} else {
		err := app.fromArgs(args)
		if err != nil {
			return 2
		}
	}

	if err := app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}
	return 0
}

type appEnv struct {
	input string
}

func (app *appEnv) fromArgs(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "no string provided")
		return fmt.Errorf("no string provided")
	}

	app.input = args[0]
	return nil
}

func (app *appEnv) fromStdin() error {
	var stdin []byte
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		stdin = append(stdin, scanner.Bytes()...)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	app.input = string(stdin)
	return nil
}

func (app *appEnv) run() error {
	fmt.Println(unpack(app.input))
	return nil
}

func unpack(input string) (string, error) {
	runes := []rune(input)
	var b strings.Builder

	for i := 0; i < len(runes); i++ {
		if unicode.IsDigit(runes[i]) && runes[i] < unicode.MaxASCII {
			if i == 0 {
				return "", nil
			}

			var num strings.Builder
			num.WriteRune(runes[i])
			letter := runes[i-1]

			for j := i + 1; j < len(runes)-1 && unicode.IsDigit(runes[j]) && runes[j] < unicode.MaxASCII; j++ {
				num.WriteRune(runes[j])
				i++
			}

			res, err := strconv.Atoi(num.String())
			if err != nil {
				return "", err
			}

			for j := 0; j < res-1; j++ {
				b.WriteRune(letter)
			}

			continue
		}
		_, err := b.WriteRune(runes[i])
		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}
