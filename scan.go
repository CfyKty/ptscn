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
		fmt.Println(fullAddress)
		conn, err := net.Dial("tcp", fullAddress)
		if err == nil {
			conn.Close() // exception could occur
			results <- p
		} else {

			results <- 0
		}
	}
}

func main() {
	flag.IntVar(&maxPort, "ports", 65535, "Ports to scan. Defaults to full range")
	turboPtr := flag.Bool("turbo", false, "Increases scan speed at the cost of accuracy. Will overwrite manual worker settings")
	workers := flag.Int("workers", 100, "Number of workers to user. Default is 100. Can cause inaccuracy if too high")

	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Please specify a domain or IP address to scan.")
		os.Exit(1)
	}
	addresses = flag.Args()[0]

	if *turboPtr {
		*workers = 140
	}

	fmt.Printf("Turbo activated. Scanning with %d workers \n", *workers)
	fmt.Printf("Scanning %s...\n", addresses)

	ports := make(chan int, *workers) // worker number
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
