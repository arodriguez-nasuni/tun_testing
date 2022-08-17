package main

import (
	"log"

	// "net"
	"log"

	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
)

func main() {
	tun, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		panic(err)
	}
	bs := make([]byte, 1500)

	i := 0
	log.SetLevel(log.DebugLevel)
	log.Infof("ready to start wireshark on: %s, isTap: %v", tun.Name(), tun.IsTAP())

	for {
		_, err := tun.Read(bs)
		if err != nil {
			log.Panic(err)
		}

		ver := (bs[0] & 0xf0) >> 4
		if ver == 4 {
			var header ipv4.Header
			header.Parse(bs)
			log.Infof("src:%s, dst:%s, protocol:%d", header.Src.String(), header.Dst.String(), header.Protocol)

		} else if ver == 6 {
			log.Infof("get ipv6 package")
		} else {
			log.Infof("unknown package")
		}

		i++
	}
}
