package checker

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type HTTPResult struct {
	Reachable      bool   `json:"reachable"`
	StatusCode     int    `json:"status_code"`
	ResponseTimeMs int64  `json:"response_time_ms"`
	HTTPSEnabled   bool   `json:"https_enabled"`
	HTTPSValid     bool   `json:"https_valid"`
	FinalURL       string `json:"final_url,omitempty"`
	Server         string `json:"server,omitempty"`
	ContentType    string `json:"content_type,omitempty"`
	Error          string `json:"error,omitempty"`
}

func HTTPCheck(url string, timeout time.Duration) (*HTTPResult, error) {
	result := &HTTPResult{
		Reachable: false,
	}

	if timeout == 0 {
		timeout = 10 * time.Second
	}

	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	start := time.Now()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		result.Error = err.Error()
		return result, nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; DomainChecker/1.0)")
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		result.Error = err.Error()
		return result, nil
	}
	defer resp.Body.Close()

	result.ResponseTimeMs = time.Since(start).Milliseconds()
	result.Reachable = true
	result.StatusCode = resp.StatusCode
	result.Server = resp.Header.Get("Server")
	result.ContentType = resp.Header.Get("Content-Type")
	result.FinalURL = resp.Request.URL.String()
	result.HTTPSEnabled = resp.Request.URL.Scheme == "https"

	return result, nil
}

func HTTPSCheck(url string, timeout time.Duration) (*HTTPResult, *SSLResult, error) {
	httpResult, err := HTTPCheck(url, timeout)
	if err != nil {
		return httpResult, nil, err
	}

	if !httpResult.Reachable {
		return httpResult, nil, nil
	}

	sslResult, err := SSLCheck(url, timeout)
	if err != nil {
		return httpResult, sslResult, nil
	}

	httpResult.HTTPSValid = sslResult.Valid

	return httpResult, sslResult, nil
}

func TLSConnectionCheck(host string, port string, timeout time.Duration) (*SSLResult, error) {
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	conn, err := tls.Dial("tcp", host+":"+port, &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS10,
	})

	if err != nil {
		return &SSLResult{
			Valid:    false,
			Error:    err.Error(),
			Issues:   []string{"TLS connection failed: " + err.Error()},
			Protocol: "unknown",
		}, nil
	}
	defer conn.Close()

	state := conn.ConnectionState()
	cert := state.PeerCertificates[0]

	sslResult := &SSLResult{
		Valid:         true,
		Issuer:        cert.Issuer.CommonName,
		Subject:       cert.Subject.CommonName,
		ExpiresAt:     cert.NotAfter,
		DaysRemaining: int(time.Until(cert.NotAfter).Hours() / 24),
		Protocol:      getProtocolVersion(state.Version),
		Issues:        []string{},
	}

	if time.Until(cert.NotAfter).Hours()/24 < 30 {
		sslResult.Issues = append(sslResult.Issues, "Certificate expires within 30 days")
	}

	return sslResult, nil
}

func getProtocolVersion(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (0x%x)", version)
	}
}
