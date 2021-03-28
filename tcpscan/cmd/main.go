package main

import (
	"flag"
	"fmt"
	"github.com/ThreeToes/golang-talk/tcpscan"
	"sort"
	"strconv"
	"strings"
)

func main() {
	addrF := flag.String("addr", "scanme.nmap.org", "Address to scan")
	portsF := flag.String("ports", "", "Ports to test")
	taskCountF := flag.Int("workers", 10, "Number of worker routines")
	flag.Parse()

	if *portsF == "" {
		fmt.Println("Must provide a value for ports (eg '80,400-1000')")
		return
	}
	fmt.Printf("Scanning address %s\n", *addrF)
	ports, err := getPorts(*portsF)
	if err != nil {
		fmt.Printf("Error in portspec: %v", err)
		return
	}
	openPorts := tcpscan.GetOpenPorts(*addrF, ports, *taskCountF)
	if len(openPorts) == 0 {
		fmt.Println("No open ports")
	} else {
		fmt.Printf("Open ports on %s:\n", *addrF)
		for _, p := range openPorts {
			fmt.Printf("\t* %d\n", p)
		}
	}
}

// getPorts Get ports from a spec string
func getPorts(portSpec string) ([]int, error) {
	split := strings.Split(portSpec, ",")
	var toScan = map[int]bool{}
	for _, s := range split {
		ports, err := enumeratePorts(s)
		if err != nil {
			return nil, err
		}
		for _, p := range ports {
			toScan[p] = true
		}
	}
	var ret []int
	for k, _ := range toScan {
		ret = append(ret, k)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})
	return ret, nil
}

// enumeratePorts process port ranges or single ports, eg 90-100 or 80
func enumeratePorts(portSpec string) ([]int, error) {
	split := strings.Split(portSpec, "-")

	lower, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, err
	}
	if len(split) == 1 {
		return []int {lower}, nil
	}
	upper, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, err
	}
	var ports []int
	for port := lower; port <= upper; port++ {
		ports = append(ports, port)
	}
	return ports, nil
}
