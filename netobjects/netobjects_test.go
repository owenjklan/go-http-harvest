package netobjects

import (
	"net/netip"
	"testing"
)

func TestNetworkAndBroadcastAddressAreIgnored(t *testing.T) {
	// Use .0 on the end to prove we're handling it properly
	// (By never returning a .0 address)
	baseAddress := "127.0.0.0"
	var returnedAddresses []netip.Addr
	var tn, _ = NewTargetNetwork(baseAddress)

	// Have it spit out all the generated addresses
	for {
		var nextIPAddr, ok = tn.NextHostAddress()
		if !ok {
			break
		}
		returnedAddresses = append(returnedAddresses, nextIPAddr)
	}

	// Confirm the length, should be 254
	if len(returnedAddresses) != 254 {
		t.Errorf("got %d addresses, expected 254", len(returnedAddresses))
	}

	// Iterate and ensure we don't have "127.0.0.1" or "127.0.0.255" in the results
	for _, address := range returnedAddresses {
		if address == netip.MustParseAddr("127.0.0.0") {
			t.Errorf("got .0 host address when we shouldn't have: %v", address)
		} else if address == netip.MustParseAddr("127.0.0.255") {
			t.Errorf("got .255 host address when we shouldn't have: %v", address)
		}
	}
}
