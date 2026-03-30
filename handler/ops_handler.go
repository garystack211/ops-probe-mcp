package handler

import (
	"time"

	"ops-probe-mcp/checker"
)

func PingCheck(host string, count int, timeout time.Duration) (*checker.PingResult, error) {
	return checker.PingCheck(host, count, timeout)
}

func HTTPHealthCheck(url string, expectedStatus int, expectedContent string, timeout time.Duration) (*checker.HealthCheckResult, error) {
	return checker.HTTPHealthCheck(url, expectedStatus, expectedContent, timeout)
}

func WhoisCheck(domain string, timeout time.Duration) (*checker.WhoisResult, error) {
	return checker.WhoisCheck(domain, timeout)
}
