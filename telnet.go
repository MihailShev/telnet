package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	conn, ctx := makeConn()

	defer closeConn(conn)

	inputCh := make(chan string)
	stopCh := make(chan struct{})

	go interruptHandle(conn)
	go readStdin(inputCh)
	go send(conn, inputCh)
	go handleResponse(conn, stopCh)

	select {
	case <-ctx.Done():
		fmt.Println("\nExit by timeout")
	case <-stopCh:
		fmt.Println("\nConnection closed")
	}

}

func makeConn() (conn net.Conn, ctx context.Context) {
	var err error

	address, port, timeout := readArgs()
	remote := fmt.Sprintf("%s:%s", address, port)
	ctx = context.Background()
	opt := ""

	if timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		opt = fmt.Sprintf(" with timeout %d milliseconds", timeout)
	}

	dialer := net.Dialer{}
	conn, err = dialer.DialContext(ctx, "tcp", remote)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connect to %s%s\n", remote, opt)

	return
}

func send(conn net.Conn, inputCh <-chan string) {
	for message := range inputCh {
		_, err := conn.Write([]byte(message))

		if err != nil {
			log.Fatal(err)
		}
	}
}

func handleResponse(conn net.Conn, stop chan<- struct{}) {
	var buf = make([]byte, 1024)

	for {
		read, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				stop <- struct{}{}
				break
			} else {
				log.Fatal(err)
			}
		}

		fmt.Print(string(buf[:read]))
	}
}

func readStdin(ch chan<- string) {
	var end = "\n"

	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		if message == end {
			ch <- fmt.Sprintf("%s%s", message, "\r\n")
		} else {
			ch <- message
		}
	}
}

func readArgs() (address string, port string, timeout int64) {
	for i, arg := range os.Args {
		switch {
		case i == 1:
			address = arg
		case i == 2:
			port = arg
		case i == 3:
			timeout = parseTimeout(arg)
		case i > 3:
			break
		}
	}

	return
}

func interruptHandle(conn net.Conn) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	for range sigCh {
		fmt.Println("\nReceived signal to interrupt, close connection")
		closeConn(conn)
		os.Exit(0)
	}
}

func closeConn(conn net.Conn) {
	err := conn.Close()

	if err != nil {
		log.Fatal(err)
	}
}

func parseTimeout(val string) int64 {
	timeout, err := strconv.ParseInt(val, 10, 64)

	if err != nil {
		log.Fatalf("Cannot parse timeout value %s%s\n", val, err)
	}

	return timeout
}
