package handler

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"ops-probe-mcp/checker"
)

type DomainCheckResult struct {
	Domain    string               `json:"domain"`
	Timestamp string               `json:"timestamp"`
	HTTP      *checker.HTTPResult  `json:"http,omitempty"`
	DNS       *checker.DNSResult   `json:"dns,omitempty"`
	SSL       *checker.SSLResult   `json:"ssl,omitempty"`
	Ports     *checker.PortsResult `json:"ports,omitempty"`
	Error     string               `json:"error,omitempty"`
}

func CheckDomain(domain string, options map[string]interface{}) (*DomainCheckResult, error) {
	result := &DomainCheckResult{
		Domain:    domain,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	domain = strings.TrimSpace(domain)
	if domain == "" {
		result.Error = "domain is required"
		return result, nil
	}

	if !strings.Contains(domain, ".") {
		result.Error = "invalid domain format"
		return result, nil
	}

	checkHTTP := getBoolOption(options, "check_http", true)
	checkDNS := getBoolOption(options, "check_dns", true)
	checkSSL := getBoolOption(options, "check_ssl", true)
	checkPorts := getBoolOption(options, "check_ports", false)

	timeout := 10 * time.Second
	if t, ok := options["timeout"].(float64); ok {
		timeout = time.Duration(t) * time.Second
	}

	if checkDNS {
		dnsResult, err := checker.DNSCheck(domain)
		if err == nil {
			result.DNS = dnsResult
		}
	}

	if checkHTTP {
		httpURL := domain
		if !strings.HasPrefix(httpURL, "http") {
			httpURL = "https://" + domain
		}

		httpResult, err := checker.HTTPCheck(httpURL, timeout)
		if err == nil {
			result.HTTP = httpResult
		}
	}

	if checkSSL {
		sslURL := domain
		if !strings.HasPrefix(sslURL, "https") {
			if !strings.HasPrefix(sslURL, "http") {
				sslURL = "https://" + domain
			} else {
				sslURL = strings.Replace(sslURL, "http://", "https://", 1)
			}
		}

		sslResult, err := checker.SSLCheck(sslURL, timeout)
		if err == nil {
			result.SSL = sslResult
		}
	}

	if checkPorts {
		ports := []int{80, 443}
		if p, ok := options["ports"].([]interface{}); ok {
			ports = make([]int, 0, len(p))
			for _, v := range p {
				if f, ok := v.(float64); ok {
					ports = append(ports, int(f))
				}
			}
		}
		portsResult := checker.PortCheck(domain, ports, timeout/2)
		result.Ports = portsResult
	}

	return result, nil
}

func CheckDomainsBatch(domains []interface{}, options map[string]interface{}) ([]*DomainCheckResult, error) {
	results := make([]*DomainCheckResult, 0, len(domains))

	for _, d := range domains {
		domain, ok := d.(string)
		if !ok {
			continue
		}
		result, err := CheckDomain(domain, options)
		if err != nil {
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

func DNSLookup(domain string, recordType string) (*checker.DNSResult, error) {
	result, err := checker.DNSCheck(domain)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func SSLCertificateCheck(domain string) (*checker.SSLResult, error) {
	sslURL := domain
	if !strings.HasPrefix(sslURL, "https") {
		if !strings.HasPrefix(sslURL, "http") {
			sslURL = "https://" + domain
		} else {
			sslURL = strings.Replace(sslURL, "http://", "https://", 1)
		}
	}

	result, err := checker.SSLCheck(sslURL, 10*time.Second)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getBoolOption(options map[string]interface{}, key string, defaultValue bool) bool {
	if options == nil {
		return defaultValue
	}
	if v, ok := options[key].(bool); ok {
		return v
	}
	return defaultValue
}

func ParseURL(rawURL string) (string, error) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	return parsed.Hostname(), nil
}

func ResultToJSON(result interface{}) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
