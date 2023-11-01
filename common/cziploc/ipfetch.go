package cziploc

import (
	"strings"

	"github.com/talkincode/esmqtt/common/iploc"
)

type IpFetch struct {
	// czipdata path
	DatPath   string
	IpLocator *iploc.Locator
}

func NewIpFetch(datPath string) *IpFetch {
	p := &IpFetch{DatPath: datPath}
	p.IpLocator, _ = iploc.Open(p.DatPath)
	return p
}

// FindRawIp find ip
func (p *IpFetch) FindRawIp(ip string) *IpAddress {
	if p.IpLocator == nil {
		return nil
	}
	return detailToIpAddress(p.IpLocator.Find(ip))
}

// FindIp find ip
func (p *IpFetch) FindIp(ip string) *IpAddress {
	if p.IpLocator == nil {
		return nil
	}
	nip, err := iploc.ParseIP(ip)
	if err != nil {
		return nil
	}
	return detailToIpAddress(p.IpLocator.FindIP(nip))
}

// DetailToIpAddress iploc.Detail to IpAddress
func detailToIpAddress(detail *iploc.Detail) *IpAddress {
	if detail == nil {
		return nil
	}
	r := &IpAddress{
		Province: detail.GetProvince(),
		City:     detail.GetCity(),
		Isp:      strings.ReplaceAll(detail.GetRegion(), "CZ88.NET", ""),
	}
	if r.Province == "" {
		r.Province = detail.GetCountry()
	}
	if r.City == "" {
		r.City = detail.GetCountry()
	}
	return r
}
