package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {

	ip := net.JoinHostPort(os.Args[1], os.Args[2])

	lis, err := net.Listen("tcp", ip)
	if err != nil {
		log.Println("Unable to open listener: ", err)
		log.Fatal()
	}
	defer lis.Close()

	for {
		// Wait for a connection.
		conn, err := lis.Accept()
		log.Println("Just accepted a connection")
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			io.Copy(c, c)
			c.Close()
		}(conn)
	}
}
