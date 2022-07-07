package telnet

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
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
	timeout  time.Duration
	addresss string
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("go-telnet", flag.ContinueOnError)
	fl.DurationVar(&app.timeout, "timeout", time.Second*10, "timeout to connect to server")

	if err := fl.Parse(args); err != nil {
		fl.Usage()
		return err
	}

	app.addresss = net.JoinHostPort(fl.Arg(0), fl.Arg(1))

	return nil
}

func (app *appEnv) run() error {
	d := net.Dialer{
		Timeout: app.timeout,
	}
	c, err := d.Dial("tcp", app.addresss)
	if err != nil {
		return err
	}
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	g := new(errgroup.Group)

	g.Go(func() error {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("n1 stop ch")
				return nil
			default:
				fmt.Print("$: ")
				t, err := reader.ReadString('\n')
				if err != nil {
					return err
				}

				_, err = fmt.Fprint(c, t)
				if err != nil {
					return err
				}
			}
		}
	})

	g.Go(func() error {
		reader := bufio.NewReader(c)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("n2 got stop")
				return nil
			default:
				t, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				fmt.Printf("got from server: %s", t)
			}
		}
	})

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		fmt.Println("\ngot interrupt signal")
		cancel()
	}()

	err = g.Wait()
	if err != io.EOF {
		return err
	}
	fmt.Println("server closed, exiting program")

	return nil
}
