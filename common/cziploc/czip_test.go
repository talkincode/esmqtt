package cziploc

import (
	"testing"
)

func TestFetchAliyunIpData(t *testing.T) {
	r, err := FetchCZApiIpData("114.114.114.114")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestFetchDat(t *testing.T) {
	err := FetchDat("/var/teamsacs/data/qqwry.dat")
	if err != nil {
		t.Fatal(err)
	}
}
