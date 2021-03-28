package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	OpenPorts()
}

func OpenPorts() {
	// Open port
	conn, err := net.DialTimeout("tcp", "scanme.nmap.org:80", 3 * time.Second)
	fmt.Printf("Port 80 open: %v\n", err == nil)
	conn.Close()

	// Closed port
	conn, err = net.DialTimeout("tcp", "scanme.nmap.org:81", 3 * time.Second)
	fmt.Printf("Port 81 open: %v", err == nil)
}
