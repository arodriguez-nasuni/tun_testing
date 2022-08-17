package main

import (
	"log"
	wintun "win_tunneler/tun"
)

const (
	DEFAULT_MTU = 1500
)

func main() {
	// Test is to create tunnel.
	tun, err := wintun.CreateTUN("testertun", 1500)
	if err != nil {
		log.Fatal(err)
	}
	name, err := tun.Name()
	if err != nil {
		tun.Close()
		log.Printf("closed")
	}
	log.Printf("Name: %s", name)
	tun.Close()
	log.Printf("Finished")

}
