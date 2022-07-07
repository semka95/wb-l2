package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	li, err := net.Listen("tcp", ":5023")
	if err != nil {
		log.Panic(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
		}

		go telnetConn(conn)
	}
}

func telnetConn(conn net.Conn) {
	defer conn.Close()
	errCh := make(chan error)

	go func() {
		for {
			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				errCh <- err
				return
			}

			log.Print("got: ", data)
			fmt.Fprintf(conn, "returning your message %s", data)
		}
	}()

	if err := <-errCh; err == io.EOF {
		log.Println("connection dropped")
	} else {
		log.Printf("got error: %v", err)
	}
}
