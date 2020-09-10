<h1 align="center">IPBlocks</h1>
<p align="center">
<a href="https://goreportcard.com/report/github.com/cbergoon/ipblocks"><img src="https://goreportcard.com/badge/github.com/cbergoon/ipblocks?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/ipblocks"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

Provides subnet information and functionality to divide subnets into blocks. This enables the parallelization of large scanning and network reconnaissance operations.

#### Documentation 

See the docs [here](https://godoc.org/github.com/cbergoon/ipblocks).

#### Install
```
go get github.com/cbergoon/ipblocks
```
```
go install github.com/cbergoon/ipblocks/cmd/ipblockscli
```

#### CLI Usage 
```
ipblockscli -s 192.168.9.0/24 -b /26 
```
```
ipblockscli -v -s 192.168.9.0/24 -b /26 
```

#### Example Usage
Below is an example demonstrating calculation of subnet information and dividing a `/24` subnet into 16 blocks of `/28` address ranges within the `/24`. 

```go
package main

import (
  "crypto/sha256"
  "log"

  "github.com/cbergoon/ipblocks"
)


func main() {
	address := net.IPv4(192, 168, 1, 0)
	originalBlock, _ := ipblocks.NewIPMaskInfo(address, 24, false)
	fmt.Println(originalBlock)

	dividedBlocks, _ := originalBlock.CalculateBlocks(28)
	fmt.Println(dividedBlocks)

	blockRange, _ := originalBlock.CalculateRange(28)
	fmt.Println(blockRange)
}

```

