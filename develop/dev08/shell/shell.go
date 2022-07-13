package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/mitchellh/go-ps"
)

func CLI() int {
	var app appEnv
	app.run()
	return 0
}

type appEnv struct {
	out io.Writer
}

func (app *appEnv) run() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		command := scanner.Text()
		if strings.Contains(command, "|") {
			if err := app.execPipe(command); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		} else {
			app.out = os.Stdout
			if err := app.execCommand(command); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}

		path, _ := filepath.Abs(".")
		fmt.Printf("%s\n$: ", path)
		scanner.Scan()
	}
}

func (app *appEnv) execCommand(command string) error {
	c := strings.Split(command, " ")

	switch c[0] {
	case "cd":
		if len(c) < 2 {
			dir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			return os.Chdir(dir)
		}
		return os.Chdir(c[1])
	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		fmt.Fprintln(app.out, pwd)
		return nil
	case "echo":
		for i := 1; i < len(c); i++ {
			fmt.Fprint(app.out, c[i], " ")
		}
		fmt.Fprintln(app.out)
		return nil
	case "kill":
		if len(c) < 2 {
			return fmt.Errorf("kill: not enough arguments")
		}

		pid, err := strconv.Atoi(c[1])
		if err != nil {
			p, err := ps.Processes()
			if err != nil {
				return err
			}
			for _, v := range p {
				if v.Executable() == c[1] {
					pid = v.Pid()
					break
				}
			}
			if pid == 0 {
				return fmt.Errorf("kill: can't find '%s' process", c[1])
			}
		}
		p, err := os.FindProcess(pid)
		if err != nil {
			return err
		}

		return p.Kill()
	case "ps":
		p, err := ps.Processes()
		if err != nil {
			return err
		}
		for _, v := range p {
			fmt.Fprintf(app.out, "%d\t%s\n", v.Pid(), v.Executable())
		}
		return nil
	case "exec":
		if len(c) < 2 {
			return fmt.Errorf("exec: not enough arguments")
		}
		binary, err := exec.LookPath(c[1])
		if err != nil {
			return err
		}
		env := os.Environ()
		return syscall.Exec(binary, c[1:], env)
	case "exit":
		fmt.Fprint(app.out, "exiting from the shell\n")
		os.Exit(0)
	default:
		return fmt.Errorf("command not found: %s", c[0])
	}
	return nil
}

func (app *appEnv) execPipe(command string) error {
	c := strings.Split(command, " | ")
	if len(c) < 2 {
		return fmt.Errorf("pipe: not enough commands: '%v'", c)
	}

	var b bytes.Buffer
	for i := 0; i < len(c); i++ {
		com := exec.Command(c[i])
		commArgs := strings.Split(c[i], " ")
		if len(commArgs) > 1 {
			com = exec.Command(commArgs[0], commArgs[1:]...)
		}

		com.Stdin = bytes.NewReader(b.Bytes())
		b.Reset()
		com.Stdout = &b

		err := com.Start()
		if err != nil {
			return err
		}
		err = com.Wait()
		if err != nil {
			return err
		}
	}

	fmt.Fprint(app.out, b.String())

	return nil
}
