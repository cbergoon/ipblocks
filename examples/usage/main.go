package main

import (
	"fmt"
	"net"

	"github.com/cbergoon/ipblocks"
)

func main() {
	address := net.IPv4(192, 168, 1, 0)
	originalBlock, _ := ipblocks.NewIPMaskInfo(address, 24, false)
	fmt.Println("Original Block: ", originalBlock)

	dividedBlocks, _ := originalBlock.CalculateBlocks(28)
	fmt.Println("Divided Blocks: ")
	for _, bl := range dividedBlocks {
		fmt.Println(bl)
	}

	blockRange, _ := originalBlock.CalculateRange(28)
	fmt.Println("Block Range: ", blockRange)
}
