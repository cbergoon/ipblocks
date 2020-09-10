package ipblocks

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"net"
)

const (
	//Bit masks representing IPv4 octets
	firstOctetBitMask  = 0xff000000
	secondOctetBitMask = 0x00ff0000
	thirdOctetBitMask  = 0x0000ff00
	fourthOctetBitMask = 0x000000ff
)

// IPMaskInfo represents critical information describing a subnet or range block.
type IPMaskInfo struct {
	Address         net.IP
	Mask            uint8
	MaskAddress     net.IP
	NetworkAddress  net.IP
	WildcardAddress net.IP
	StartAddress    net.IP
	EndAddress      net.IP
}

// NewIPMaskInfo calculates information about the subnet or range block.
func NewIPMaskInfo(addr net.IP, mask uint8, rangeBlock bool) (*IPMaskInfo, error) {
	if mask < 0 || mask > 31 {
		return nil, errors.New("error: mask must be integer between 0 and 31")
	}

	ip4 := addr.To4()

	maskValue := caclulateMaskValue(mask)
	maskAddress := caclulateMaskAddress(maskValue)
	wildcardAddress := caclulateWildcardAddress(maskAddress)
	networkAddress := caclulateNetworkAddress(ip4, maskAddress) // might need plus one
	startAddress := caclulateStartAddress(networkAddress, rangeBlock)
	endAddress := caclulateEndAddress(startAddress, wildcardAddress)

	return &IPMaskInfo{
		Address:         addr,
		Mask:            mask,
		MaskAddress:     byteArrayToIPv4(maskAddress),
		NetworkAddress:  byteArrayToIPv4(networkAddress),
		WildcardAddress: byteArrayToIPv4(wildcardAddress),
		StartAddress:    byteArrayToIPv4(startAddress),
		EndAddress:      byteArrayToIPv4(endAddress),
	}, nil
}

// CalculateBlocks divides an existing IPMaskInfo type into smaller consecutive subnets contained
// within the existing subnet.
func (ipmi *IPMaskInfo) CalculateBlocks(mask uint8) ([]*IPMaskInfo, error) {
	if ipmi.Mask > mask {
		return nil, errors.New("error: block mask size must be less than initial")
	}
	var blocks []*IPMaskInfo
	block, err := NewIPMaskInfo(ipmi.StartAddress, mask, true)
	if err != nil {
		return nil, err
	}
	blocks = append(blocks, block)
	for !equals(block.EndAddress.To4(), ipmi.EndAddress.To4()) {
		nextAddress := make([]byte, 4, 4)
		copy(nextAddress, block.EndAddress.To4())
		nextAddress = addOne(nextAddress)
		nextAddressIP := net.IPv4(nextAddress[0], nextAddress[1], nextAddress[2], nextAddress[3])
		block, err = NewIPMaskInfo(nextAddressIP, mask, true)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

// CalculateRange returns an array of strings containing the CIDR represenation of subnet contained
// within the existing subnet.
func (ipmi *IPMaskInfo) CalculateRange(mask uint8) ([]string, error) {
	addresses := []string{}
	dividedBlocks, err := ipmi.CalculateBlocks(mask)
	if err != nil {
		return nil, err
	}
	for _, block := range dividedBlocks {
		address := fmt.Sprintf("%s/%d", block.NetworkAddress, block.Mask)
		addresses = append(addresses, address)
	}
	return addresses, nil
}

// String returns plain text representation of IPMaskInfo
func (ipmi *IPMaskInfo) String() string {
	return fmt.Sprintf("Address: %s Mask: %d Mask Address: %s Wildcard Address: %s Network Address: %s Start Address: %s End Address: %s", ipmi.Address, ipmi.Mask, ipmi.MaskAddress, ipmi.WildcardAddress, ipmi.NetworkAddress, ipmi.StartAddress, ipmi.EndAddress)
}

func addOne(addr []byte) []byte {
	result := addr
	for i := 3; i >= 0; i-- {
		if addr[i] < byte(255) {
			result[i] = addr[i] + 1
			break
		} else {
			result[i] = byte(0)
		}
	}
	return result
}

func equals(addr1, addr2 []byte) bool {
	for i := 0; i < 4; i++ {
		if addr1[i] != addr2[i] {
			return false
		}
	}
	return true
}

func byteArrayToIPv4(ipBytes []byte) net.IP {
	return net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
}

func caclulateMaskValue(mask uint8) uint32 {
	mv := (uint32(math.Pow(2, float64(mask))) - 1) << uint32(32-mask)
	return mv
}

func caclulateMaskAddress(maskValue uint32) []byte {
	ma := []byte{
		math.MaxUint8 << uint32(8-(bits.OnesCount32(maskValue&firstOctetBitMask))),
		math.MaxUint8 << uint32(8-(bits.OnesCount32(maskValue&secondOctetBitMask))),
		math.MaxUint8 << uint32(8-(bits.OnesCount32(maskValue&thirdOctetBitMask))),
		math.MaxUint8 << uint32(8-(bits.OnesCount32(maskValue&fourthOctetBitMask))),
	}
	return ma
}

func caclulateWildcardAddress(maskAddress []byte) []byte {
	wa := []byte{
		maskAddress[0] ^ math.MaxUint8,
		maskAddress[1] ^ math.MaxUint8,
		maskAddress[2] ^ math.MaxUint8,
		maskAddress[3] ^ math.MaxUint8,
	}
	return wa
}

func caclulateNetworkAddress(address, maskAddress []byte) []byte {
	na := []byte{
		address[0] & maskAddress[0],
		address[1] & maskAddress[1],
		address[2] & maskAddress[2],
		address[3] & maskAddress[3],
	}
	return na
}

func caclulateStartAddress(networkAddress []byte, rangeBlock bool) []byte {
	sa := make([]byte, 4, 4)
	copy(sa, networkAddress)
	if !rangeBlock {
		sa[3] = sa[3] + 1
	}
	return sa
}

func caclulateEndAddress(startAddress, wildcardAddress []byte) []byte {
	ea := []byte{
		startAddress[0] | wildcardAddress[0],
		startAddress[1] | wildcardAddress[1],
		startAddress[2] | wildcardAddress[2],
		startAddress[3] | wildcardAddress[3],
	}
	return ea
}
