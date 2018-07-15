package gluadns

import (
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"
	glua "github.com/yuin/gopher-lua"
)

var api = map[string]glua.LGFunction{
	"exchange": exchange,
}

func exchange(L *glua.LState) int {
	var sIP string

	q := L.CheckString(1)
	t := L.CheckString(2)
	host := L.CheckString(3)
	if L.GetTop() > 3 {
		sIP = L.CheckString(4)
	}

	rrType, ok := dns.StringToType[t]
	if !ok {
		L.Push(glua.LNil)
		L.Push(glua.LString("RRType is not support"))
		return 2
	}

	m := new(dns.Msg)
	m.SetQuestion(q, rrType)
	if sIP != "" {
		if err := setClientSubnet(m, sIP); err != nil {
			L.Push(glua.LNil)
			L.Push(glua.LString(err.Error()))
			return 2
		}
	}

	if !strings.Contains(host, ":") {
		host = net.JoinHostPort(host, "53")
	}

	r, err := dns.Exchange(m, host)
	if err != nil {
		L.Push(glua.LNil)
		L.Push(glua.LString(err.Error()))
		return 2
	}

	if r.Truncated {
		L.Push(glua.LNil)
		L.Push(glua.LString("DNS response is truncated"))
		return 2
	}

	ltable := new(glua.LTable)
	vtable := new(glua.LTable)
	for _, rr := range r.Answer {
		rh := rr.Header()
		lrr := new(glua.LTable)
		lrr.RawSetString("ttl", glua.LNumber(rh.Ttl))
		lrr.RawSetString("name", glua.LString(rh.Name))
		lrr.RawSetString("type", glua.LString(dns.TypeToString[rh.Rrtype]))
		switch v := rr.(type) {
		case *dns.A:
			lrr.RawSetString("value", glua.LString(v.A.String()))
		case *dns.AAAA:
			lrr.RawSetString("value", glua.LString(v.AAAA.String()))
		case *dns.CNAME:
			lrr.RawSetString("value", glua.LString(v.Target))
		}
		vtable.Append(lrr)
	}

	ltable.RawSetString("rcode", glua.LNumber(r.Rcode))
	ltable.RawSetString("answer", vtable)

	L.Push(ltable)
	L.Push(glua.LNil)
	return 2
}

func setClientSubnet(m *dns.Msg, sourceIPv4 string) error {
	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_SUBNET)
	e.Code = dns.EDNS0SUBNET
	e.Family = 1         // 1 for IPv4 source address, 2 for IPv6
	e.SourceNetmask = 32 // 32 for IPV4, 128 for IPv6
	e.SourceScope = 0
	ip := net.ParseIP(sourceIPv4)
	if ip == nil {
		return fmt.Errorf("source ip format is not a valid IPv4 address %s", sourceIPv4)
	}

	e.Address = ip.To4() // for IPv4
	o.Option = append(o.Option, e)

	m.Extra = append(m.Extra, o)
	return nil
}
