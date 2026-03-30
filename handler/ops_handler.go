package handler

import (
	"strings"
	"time"

	"ops-probe-mcp/checker"
)

func PingCheck(host string, count int, timeout time.Duration) (*checker.PingResult, error) {
	return checker.PingCheck(host, count, timeout)
}

func HTTPHealthCheck(url string, expectedStatus int, expectedContent string, timeout time.Duration, resolveIP string) (*checker.HealthCheckResult, error) {
	return checker.HTTPHealthCheck(url, expectedStatus, expectedContent, timeout, resolveIP)
}

func SSLCertDetailCheck(domain string, resolveIP string) (*checker.SSLCertDetail, error) {
	sslURL := domain
	if !strings.HasPrefix(sslURL, "https") {
		if !strings.HasPrefix(sslURL, "http") {
			sslURL = "https://" + domain
		} else {
			sslURL = strings.Replace(sslURL, "http://", "https://", 1)
		}
	}
	return checker.SSLCertDetailCheck(sslURL, 10*time.Second, resolveIP)
}

func WhoisCheck(domain string, timeout time.Duration) (*checker.WhoisResult, error) {
	return checker.WhoisCheck(domain, timeout)
}
