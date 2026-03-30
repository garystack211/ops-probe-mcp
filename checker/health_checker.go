package checker

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HealthCheckResult struct {
	URL            string `json:"url"`
	ResolvedIP     string `json:"resolved_ip,omitempty"`
	StatusCode     int    `json:"status_code"`
	Healthy        bool   `json:"healthy"`
	ResponseTimeMs int64  `json:"response_time_ms"`
	ContentMatch   bool   `json:"content_match,omitempty"`
	Error          string `json:"error,omitempty"`
}

// HTTPHealthCheck checks an HTTP endpoint.
// If resolveIP is non-empty, the TCP connection is made directly to that IP
// (bypassing DNS), while the Host header and TLS SNI still use the original hostname.
// This is useful for testing origin servers behind a CDN.
func HTTPHealthCheck(rawURL string, expectedStatus int, expectedContent string, timeout time.Duration, resolveIP string) (*HealthCheckResult, error) {
	result := &HealthCheckResult{
		URL:     rawURL,
		Healthy: false,
	}

	if timeout == 0 {
		timeout = 10 * time.Second
	}

	var client *http.Client

	if resolveIP != "" {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			result.Error = "invalid url: " + err.Error()
			return result, nil
		}
		hostname := parsed.Hostname()
		result.ResolvedIP = resolveIP

		dialer := &net.Dialer{Timeout: timeout}
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{ServerName: hostname},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// addr is "hostname:port" — replace hostname with the specified IP
				_, port, err := net.SplitHostPort(addr)
				if err != nil {
					port = "80"
					if parsed.Scheme == "https" {
						port = "443"
					}
				}
				return dialer.DialContext(ctx, network, net.JoinHostPort(resolveIP, port))
			},
		}
		client = &http.Client{Timeout: timeout, Transport: transport}
	} else {
		client = &http.Client{Timeout: timeout}
	}

	start := time.Now()
	resp, err := client.Get(rawURL)
	result.ResponseTimeMs = time.Since(start).Milliseconds()

	if err != nil {
		result.Error = err.Error()
		return result, nil
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	// 检查状态码
	if expectedStatus > 0 {
		result.Healthy = resp.StatusCode == expectedStatus
	} else {
		result.Healthy = resp.StatusCode >= 200 && resp.StatusCode < 400
	}

	// 检查内容匹配
	if expectedContent != "" {
		body := make([]byte, 4096)
		n, _ := resp.Body.Read(body)
		result.ContentMatch = strings.Contains(string(body[:n]), expectedContent)
		result.Healthy = result.Healthy && result.ContentMatch
	}

	return result, nil
}
