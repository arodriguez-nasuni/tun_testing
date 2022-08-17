package main

import (
	"log"
	"win_tunneler/tun"
)

const (
	DEFAULT_MTU = 1500
)

func main() {
	// Test is to create tunnel.
	tun, err := CreateTUN("testertun", 1500)
	if err != nil {
		log.Fatal(err)
	}

}
