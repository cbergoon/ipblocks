package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/cbergoon/ipblocks"
)

func main() {

	s := flag.String("s", "192.168.1.0/24", "CIDR representation of subnet to divide")
	b := flag.String("b", "/28", "size of block to divide subnet into")
	v := flag.Bool("v", false, "output verbose information describing subnet")

	flag.Parse()

	saddress := *s
	sblock := *b // uint8(strconv.ParseUint(sblock, 10, 64))
	sverbose := *v

	if strings.HasPrefix(sblock, "/") {
		sblock = sblock[1:]
	}

	divideMask, err := strconv.ParseUint(sblock, 10, 64)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	ip, ipnet, err := net.ParseCIDR(saddress)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	maskSize, _ := ipnet.Mask.Size()

	if sverbose {
		fmt.Printf("Calculating %s/%d as blocks of /%d \n", ip, maskSize, divideMask)
	}

	original, err := ipblocks.NewIPMaskInfo(ip, uint8(maskSize), false)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if sverbose {
		fmt.Println()
		fmt.Printf("Subnet: %s \n", original)
	}

	dividedBlocks, _ := original.CalculateBlocks(uint8(divideMask))
	if sverbose {
		fmt.Println()
		fmt.Println("Divided Blocks: ")
		for _, bl := range dividedBlocks {
			fmt.Println(bl)
		}
	}

	blockRange, err := original.CalculateRange(uint8(divideMask))
	if err != nil {
		log.Fatalf("error: %v")
	}
	if sverbose {
		fmt.Println()
		fmt.Println("Block Range: ")
	}
	for _, br := range blockRange {
		fmt.Println(br)
	}

}
