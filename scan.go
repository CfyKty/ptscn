package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
)

var addresses string
var maxPort int

func worker(ports, results chan int) {
	for p := range ports {
		fullAddress := fmt.Sprintf("%s:%d", addresses, p)
		//fmt.Println(fullAddress)
		conn, err := net.Dial("tcp", fullAddress)
		if err == nil {
			conn.Close()
			results <- p
		} else {

			results <- 0
		}
	}
}

func main() {
	flag.IntVar(&maxPort, "Ports", 1024, "Ports to scan")
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Incorrect ")
		os.Exit(1)
	}
	addresses = flag.Args()[0]

	fmt.Printf("Scanning %s...\n", addresses)

	ports := make(chan int, 100) // worker number
	results := make(chan int)
	var openports []int

	for i := 1; i <= cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= maxPort; i++ {
			ports <- i
		}
	}()

	for i := 0; i < maxPort; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)

	fmt.Println("Scan complete")
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
