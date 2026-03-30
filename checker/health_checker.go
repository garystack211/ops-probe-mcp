package checker

import (
	"net/http"
	"strings"
	"time"
)

type HealthCheckResult struct {
	URL            string `json:"url"`
	StatusCode     int    `json:"status_code"`
	Healthy        bool   `json:"healthy"`
	ResponseTimeMs int64  `json:"response_time_ms"`
	ContentMatch   bool   `json:"content_match,omitempty"`
	Error          string `json:"error,omitempty"`
}

func HTTPHealthCheck(url string, expectedStatus int, expectedContent string, timeout time.Duration) (*HealthCheckResult, error) {
	result := &HealthCheckResult{
		URL:     url,
		Healthy: false,
	}

	if timeout == 0 {
		timeout = 10 * time.Second
	}

	client := &http.Client{Timeout: timeout}
	start := time.Now()

	resp, err := client.Get(url)
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
