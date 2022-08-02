package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"

	color "github.com/fatih/color"
	"github.com/songgao/water"
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
)

var (
	remoteIp      = ""
	tunnelIpCount = 20
	os_type       = 0
)

func setupTunnel(tunName, remoteIp, localIP string) error {
	if os_type == OS_MAC {
		return runMacSetup(tunName, remoteIp, localIP)
	} else {
		return errors.New("Unimplemented for this OS")
	}
}

func runMacSetup(tunName, remoteIp, localIp string) error {
	cmd := exec.Command("ifconfig", tunName, remoteIp, localIp, "up")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run
	if nil != err() {
		return err()
	}
	return nil
}

func createTunnel(remoteIp string) (*water.Interface, error) {
	// localIp := fmt.Sprintf("10.1.0.%d", tunnelIpCount)
	tunnelIpCount += 1
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		return nil, err
	}
	// err = setupTunnel(ifce.Name(), remoteIp, localIp)
	log.Printf("New tunnel up as %s", ifce.Name())
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
		fmt.Println("Syntax: ./listener ipaddr port")
		fmt.Println("Example ./listener 172.20.20.1 445")
	}

	ip := net.JoinHostPort(os.Args[1], os.Args[2])
	remoteIp = fmt.Sprintf(os.Args[1])

	log.Printf("%s: enter 'new' to create a new tunnel to %s:%s\n", warningMsg, os.Args[1], os.Args[2])

	lis, err := net.Listen("tcp4", ip)
	if err != nil {
		log.Println("Unable to open listener: ", err)
		log.Fatal()
	}
	defer lis.Close()

	go func() {
		for {
			var command string
			fmt.Scanln(&command)
			if command == "new" {
				log.Println("Creating a new tunnel...")
				tun, err := createTunnel("10.1.0.10")
				if err != nil {
					log.Printf("%s: %s", errorMsg, err)
				}
				go func(ifce *water.Interface) {
					packet := make([]byte, 2000)
					for {
						n, err := ifce.Read(packet)
						if err != nil {
							log.Fatal(err)
						}
						log.Printf("Packet Received: % x\n", packet[:n])
					}
				}(tun)
			} else if command == "cleanup" {
				log.Println("Begining Tunnel cleanup...")
				// cleanup()
			} else {
				log.Printf("'%s' is not a valid command", command)
			}
		}
	}()

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
