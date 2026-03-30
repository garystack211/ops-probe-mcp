package checker

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type PortResult struct {
	Open    bool   `json:"open"`
	Latency string `json:"latency,omitempty"`
	Error   string `json:"error,omitempty"`
}

type PortsResult struct {
	PortResults map[string]*PortResult `json:"port_results"`
}

func PortCheck(host string, ports []int, timeout time.Duration) *PortsResult {
	result := &PortsResult{
		PortResults: make(map[string]*PortResult),
	}

	if timeout == 0 {
		timeout = 5 * time.Second
	}

	for _, port := range ports {
		portStr := strconv.Itoa(port)
		start := time.Now()

		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, portStr), timeout)

		latency := time.Since(start)

		if err != nil {
			result.PortResults[portStr] = &PortResult{
				Open:  false,
				Error: err.Error(),
			}
		} else {
			conn.Close()
			result.PortResults[portStr] = &PortResult{
				Open:    true,
				Latency: latency.String(),
			}
		}
	}

	return result
}

func SinglePortCheck(host string, port int, timeout time.Duration) *PortResult {
	result := &PortResult{
		Open: false,
	}

	if timeout == 0 {
		timeout = 5 * time.Second
	}

	start := time.Now()

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), timeout)

	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer conn.Close()

	result.Open = true
	result.Latency = time.Since(start).String()

	return result
}

func CommonPorts() []int {
	return []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 465, 587, 993, 995, 3306, 3389, 5432, 6379, 8080, 8443}
}

func CheckDomainPorts(domain string, ports []int, timeout time.Duration) (*PortsResult, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve domain: %w", err)
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no IP addresses found for domain")
	}

	host := ips[0].String()
	return PortCheck(host, ports, timeout), nil
}
