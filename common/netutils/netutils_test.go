package netutils

import (
	"net"
	"testing"

	"github.com/c-robinson/iplib"
)

func TestParseIpNet(t *testing.T) {
	n, err := ParseIpNet("10.12.1.0/24")
	if err != nil {
		t.Error(err)
	}
	t.Log(n.String())
	t.Log(n.Mask())
	var maskhash = n.Mask().String()
	t.Log(iplib.HexStringToIP(maskhash))
	t.Log(iplib.IP4ToUint32(n.IP()))
	t.Log(net.IP(n.Mask()).String())
	t.Log(n.Mask().Size())
}
