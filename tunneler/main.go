package main

import (
	"errors"
	"fmt"
	// "net"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	color "github.com/fatih/color"
	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
	parse "github.com/ilyaigpetrov/parse-tcp-go"

)

var (
	errorMsg   = color.New(color.FgRed).SprintFunc()("ERROR")
	successMsg = color.New(color.FgGreen).SprintFunc()("SUCCESS")
	warningMsg = color.New(color.FgYellow).SprintFunc()("INFO")
)

const (
	OS_MAC     = 1
	OS_WINDOWS = 2
	OS_LINUX   = 3

	BUFFERSIZE = 1600
	MTU        = "1300"
)

var (
	remoteIp      = ""
	localIp 	  = ""
	tunnelIpCount = 10
	os_type       = 0
)

func setupTunnel(tunName, ip string) error {
	if os_type == OS_MAC {
		return runMacSetup(tunName, ip)
	} else if os_type == OS_LINUX {
		return runLinuxSetup(tunName, ip)
	} else {
		return errors.New("Unimplemented for this OS")
	}
}

func runBin(bin string, args ...string) {
    cmd := exec.Command(bin, args...)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    fatalIf(cmd.Run())
}

func fatalIf(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func runMacSetup(tunName, ip string) error {
	log.Fatalf("%s: UNIMPLEMENTED", errorMsg)
	// runBin("ifconfig", tunName, remoteIp, localIp, "up")
	return nil
}

func runLinuxSetup(tunName, ip string) error {
    runBin("/bin/ip", "link", "set", "dev", tunName, "mtu", MTU)
    runBin("/bin/ip", "addr", "add", ip, "dev", tunName)
    runBin("/bin/ip", "link", "set", "dev", tunName, "up")  
	return nil
}

func createTunnel(ip string) (*water.Interface, error) {
	tunnelIpCount += 1
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		return nil, err
	}
	err = setupTunnel(ifce.Name(), ip)
	log.Printf("New tunnel up as %s listening on %s\n\n", ifce.Name(), ip)
	if err != nil {
		return nil, err
	}
	return ifce, nil
}

func main() {
	running_os := runtime.GOOS

	switch running_os {
	case "windows":
		log.Printf("%s: Windows", successMsg)
		os_type = OS_WINDOWS
	case "darwin":
		log.Printf("%s: MAC operating system", successMsg)
		os_type = OS_MAC
	case "linux":
		log.Printf("%s: RUNNING: Linux", successMsg)
		os_type = OS_LINUX
	default:
		log.Fatalf("%s: this OS is not supported", errorMsg)
	}


	if len(os.Args) < 3 {
		fmt.Println("Syntax: ./tunneler local_addr dest_addr")
		fmt.Println("Example ./listener 11.11.11.2/24 127.0.0.1")
	}

	remoteIp = fmt.Sprintf("%s", os.Args[2])
	localIp = fmt.Sprintf("%s", os.Args[1])

	// laddr := net.IPAddr{IP: net.ParseIP(fmt.Sprintf("%s,%d",remoteIp,445))}
	// raddr := &net.IPAddr{IP: net.ParseIP(fmt.Sprintf("%s,%d",remoteIp,445))}
	// conn, err := net.DialIP("ip4:tcp", nil, raddr)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// defer conn.Close()

	log.Printf("%s: Making tunnel from %s to %s\n", warningMsg, localIp, remoteIp)

	src_iface, err := createTunnel(localIp)
	if err != nil {
		log.Fatalf("%s: %s", errorMsg, err)
	}

	go func() {
		for  {
			packet := make([]byte, BUFFERSIZE)
			for {
				n, err := src_iface.Read(packet)
				if err != nil {
					log.Fatalf("Error: %s",err)
				}
				header, _ := ipv4.ParseHeader(packet[:n])
				if header.Protocol == 6 {
					tcpPacket, err := parse.ParseTCPPacket(packet[:n])
					if err != nil {
						log.Fatalf("Error: %s", err)
					}
					log.Println("Received TCP Packet:")
					tcpPacket.Print()
					// if _, err := conn.Write(packet[:n]); err != nil {
					// 	fmt.Println(err.Error())
					// 	return
					// }
				}
			}
		}
	}()

	func() {
		for{
			time.Sleep(1 * time.Second)
		}
		
	}()

}
