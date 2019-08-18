package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main()  {
	dialer := &net.Dialer{}
	conn, err := dialer.Dial("tcp", "0.0.0.0:3001")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Println("start connection to 0.0.0.0:3001")

	_, err = conn.Write([]byte("GET /selection HTTP/1.1\r\n\r\n"))


	if err != nil {
		log.Fatal()
	}
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(message)
}