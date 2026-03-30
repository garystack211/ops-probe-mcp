package checker

import (
	"net"
	"time"

	"github.com/miekg/dns"
)

type DNSResult struct {
	Resolved bool     `json:"resolved"`
	IPv4     []string `json:"ipv4"`
	IPv6     []string `json:"ipv6"`
	CNAME    *string  `json:"cname"`
	MX       []string `json:"mx,omitempty"`
	TXT      []string `json:"txt,omitempty"`
	Error    string   `json:"error,omitempty"`
}

func DNSCheck(domain string, servers ...string) (*DNSResult, error) {
	result := &DNSResult{
		Resolved: false,
		IPv4:     []string{},
		IPv6:     []string{},
	}

	m := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id:               dns.Id(),
			RecursionDesired: true,
		},
	}

	m.Question = []dns.Question{
		{Name: dns.Fqdn(domain), Qtype: dns.TypeA, Qclass: dns.ClassINET},
	}

	client := &dns.Client{
		Timeout: 5 * time.Second,
	}

	if len(servers) > 0 {
		client.Net = "tcp"
	} else {
		client.Net = "udp"
	}

	var in *dns.Msg
	var err error

	if len(servers) > 0 {
		in, _, err = client.Exchange(m, servers[0]+":53")
	} else {
		in, _, err = client.Exchange(m, "8.8.8.8:53")
	}

	if err != nil {
		result.Error = err.Error()
		return result, nil
	}

	if in == nil || len(in.Answer) == 0 {
		return result, nil
	}

	for _, ans := range in.Answer {
		switch a := ans.(type) {
		case *dns.A:
			result.IPv4 = append(result.IPv4, a.A.String())
			result.Resolved = true
		case *dns.AAAA:
			result.IPv6 = append(result.IPv6, a.AAAA.String())
			result.Resolved = true
		case *dns.CNAME:
			result.CNAME = &a.Target
			result.Resolved = true
		case *dns.MX:
			result.MX = append(result.MX, a.Mx)
			result.Resolved = true
		case *dns.TXT:
			result.TXT = append(result.TXT, a.Txt...)
			result.Resolved = true
		}
	}

	if len(result.IPv4) == 0 && len(result.IPv6) == 0 {
		ips, err := net.LookupIP(domain)
		if err == nil {
			for _, ip := range ips {
				if ip.To4() != nil {
					result.IPv4 = append(result.IPv4, ip.String())
				} else {
					result.IPv6 = append(result.IPv6, ip.String())
				}
				result.Resolved = true
			}
		}
	}

	return result, nil
}
