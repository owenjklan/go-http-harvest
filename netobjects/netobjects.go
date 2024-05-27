package netobjects

import (
	"fmt"
	"go-http-harvest/byteseq"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

type DiscoveredHost struct {
	IPAddress       net.IP
	Port            int
	ResponseHeaders map[string]string
	StatusCode      int
}

type TargetNetwork struct {
	BaseAddress           net.IP
	baseAddressComponents [3]string
	baseAddressBytes      [3]byte
	lowestByteSeq         *byteseq.RandomByteSeq
	ResolveAddress        net.IP
	DNSResult             string
}

func NewTargetNetwork(baseAddress string) (*TargetNetwork, error) {
	addressComponents, err := determineBaseAddressString(baseAddress)
	var byteComponents [3]byte
	if err != nil {
		return nil, err
	}

	// populate byte components
	for index, component := range addressComponents {
		byteValue, _ := strconv.ParseUint(component, 10, 8)
		byteComponents[index] = byte(byteValue)
	}

	return &TargetNetwork{
		BaseAddress:           net.ParseIP(baseAddress),
		baseAddressComponents: addressComponents,
		baseAddressBytes:      byteComponents,
		lowestByteSeq:         byteseq.NewRandomSeq([]byte{0x00, 0xFF}),
		ResolveAddress:        nil,
		DNSResult:             "",
	}, nil
}

func (tn *TargetNetwork) NextHostAddress() (netip.Addr, bool) {
	if !tn.lowestByteSeq.HasMore() {
		return netip.Addr{}, false
	}

	lastByte, err := tn.lowestByteSeq.NextValue()
	if err != nil {
		return netip.Addr{}, false
	}

	// Construct IP address from components
	var components = [4]byte(append(tn.baseAddressBytes[0:3], lastByte))

	var parsedIP = netip.AddrFrom4(components)
	return parsedIP, true
}

func determineBaseAddressString(suppliedIP string) ([3]string, error) {
	addressComponents := strings.Split(suppliedIP, ".")

	if len(addressComponents) != 4 {
		return [3]string{}, fmt.Errorf("supplied IPv4 address is invalid")
	}
	upperOctets := [3]string(addressComponents[0:4])

	return upperOctets, nil
}
