package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Ullaakut/nmap"
	"github.com/cbergoon/ipblocks"
)

func run(wg *sync.WaitGroup, addr string) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(addr),
		nmap.WithPorts("80,443,843"),
		nmap.WithContext(ctx),
	)
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	result, warnings, err := scanner.Run()
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	if warnings != nil {
		log.Printf("Warnings: \n %v", warnings)
	}

	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		fmt.Printf("Host %q:\n", host.Addresses[0])

		for _, port := range host.Ports {
			fmt.Printf("\tPort %d/%s %s %s\n", port.ID, port.Protocol, port.State, port.Service.Name)
		}
	}

	fmt.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
}

func main() {

	originalBlock, _ := ipblocks.NewIPMaskInfo([]byte{192, 168, 1, 0}, 24, false)
	blockRange, _ := originalBlock.CalculateRange(28)

	wg := &sync.WaitGroup{}
	for _, db := range blockRange {
		fmt.Printf("Scanning: %s \n", db)
		wg.Add(1)
		go run(wg, db)
	}

	wg.Wait()
}
