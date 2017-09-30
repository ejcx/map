package port

import (
	"testing"
)

var (
	p = new(PortScan)
)

func TestParseSetCIDR(t *testing.T) {
	err := p.ParseSetCIDR("10.0.0.0/24")
	if err != nil {
		t.Errorf("Could not parse address: %s", err)
	}
}

func TestEnumerate(t *testing.T) {
	err := p.ParseSetCIDR("10.0.0.0/24")
	if err != nil {
		t.Errorf("Could not parse address: %s", err)
	}
	ipList := ipEnumerator(p.CIDR)
	if len(ipList) != 256 {
		t.Error("invalid subnet enumeration")
	}
}
