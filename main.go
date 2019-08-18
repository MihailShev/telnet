package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	dialer := &net.Dialer{}
	address, port := readArgs()

	target := fmt.Sprintf("%s:%s", address, port)
	conn, err := dialer.Dial("tcp", target)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Connect to %s\n", target)

	input := make(chan string)

	go readStdin(input)
	go handle(conn)

	for m := range input {
		_, err := conn.Write([]byte(m))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func handle(conn net.Conn) {
	for {

		message, err := ioutil.ReadAll(conn)

		if err != nil {
			log.Fatal(err)
		}

		if len(message) > 0 {
			fmt.Println("Received message:")
			fmt.Println(string(message))
		}

		conn.SetReadDeadline(time.Now())

		if _, err = conn.Read(make([]byte, 0)); err == io.EOF {
			fmt.Println("Connection is closed")
		} else {
			fmt.Println("connection is live")
			conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		}
	}
}

func readArgs() (address string, port string) {
	for i, arg := range os.Args {
		switch {
		case i == 1:
			address = arg
		case i == 2:
			port = arg
		case i > 2:
			break
		}
	}

	return
}

func readStdin(ch chan<- string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		ch <- message + "\r\n"
	}
}
