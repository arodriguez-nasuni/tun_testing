package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Syntax: ./dialer ipaddr port [message]")
		fmt.Println("Example ./dialer 172.20.20.1 445 the fox ran fast")
	}

	message := "Hello World"
	if len(os.Args) > 3 {
		message = ""
		for _, arg := range os.Args[3:] {
			message = message + arg
		}
	}
	ip := net.JoinHostPort(os.Args[1], os.Args[2])

	conn, err := net.Dial("tcp4", ip)
	if err != nil {
		log.Println("Unable to connect: ", err)
		log.Fatal()
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Println("Unable to write: ", err)
		log.Fatal()
	}

	echo := make([]byte, 100)
	nbr, err := conn.Read(echo)
	if err != nil {
		log.Println("Unable to read: ", err)
		log.Fatal()
	}

	log.Println("Got message: ", string(echo[:nbr]))
}
