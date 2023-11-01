package cziploc

import (
	"testing"
)

func TestIpFetch_FindRawIp(t *testing.T) {
	ipf := NewIpFetch("/var/teamsacs/data/qqwry.dat")
	r := ipf.FindRawIp("114.114.114.114")
	t.Log(r)

}
