package checker

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type SSLResult struct {
	Valid         bool      `json:"valid"`
	Issuer        string    `json:"issuer"`
	Subject       string    `json:"subject,omitempty"`
	ExpiresAt     time.Time `json:"expires_at"`
	DaysRemaining int       `json:"days_remaining"`
	Protocol      string    `json:"protocol,omitempty"`
	Issues        []string  `json:"issues"`
	Error         string    `json:"error,omitempty"`
}

func SSLCheck(targetURL string, timeout time.Duration, resolveIP ...string) (*SSLResult, error) {
	result := &SSLResult{
		Valid:  false,
		Issues: []string{},
	}

	if timeout == 0 {
		timeout = 10 * time.Second
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		result.Error = "Invalid URL: " + err.Error()
		return result, nil
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port == "" {
		port = "443"
	}

	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS10,
	}

	var conn *tls.Conn
	if len(resolveIP) > 0 && resolveIP[0] != "" {
		ip := resolveIP[0]
		dialer := &net.Dialer{Timeout: timeout}
		rawConn, dialErr := dialer.DialContext(context.Background(), "tcp", net.JoinHostPort(ip, port))
		if dialErr != nil {
			err = dialErr
		} else {
			conn = tls.Client(rawConn, tlsConfig)
			if err2 := conn.Handshake(); err2 != nil {
				rawConn.Close()
				err = err2
				conn = nil
			}
		}
	} else {
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", host+":"+port, tlsConfig)
	}

	if err != nil {
		if strings.Contains(err.Error(), "certificate has expired") {
			result.Issues = append(result.Issues, "Certificate has expired")
			result.DaysRemaining = -1
		} else {
			result.Error = err.Error()
			result.Issues = append(result.Issues, "TLS connection failed: "+err.Error())
		}
		return result, nil
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		result.Error = "No certificates found"
		return result, nil
	}

	cert := state.PeerCertificates[0]

	result.Valid = true
	result.ExpiresAt = cert.NotAfter
	result.DaysRemaining = int(time.Until(cert.NotAfter).Hours() / 24)
	result.Protocol = getSSLProtocolVersion(state.Version)

	if cert.IsCA {
		result.Issues = append(result.Issues, "Certificate is a CA certificate")
	}

	result.Issuer = cert.Issuer.CommonName
	result.Subject = cert.Subject.CommonName

	if time.Until(cert.NotAfter).Hours()/24 < 30 {
		result.Issues = append(result.Issues, "Certificate expires within 30 days")
	}

	if time.Until(cert.NotBefore).Hours() > 0 {
		result.Issues = append(result.Issues, "Certificate is not yet valid")
	}

	// Build intermediate certificate pool from the chain returned by the server
	intermediates := x509.NewCertPool()
	for _, ic := range state.PeerCertificates[1:] {
		intermediates.AddCert(ic)
	}

	opts := x509.VerifyOptions{
		DNSName:       host,
		Intermediates: intermediates,
	}
	_, err = cert.Verify(opts)
	if err != nil {
		result.Valid = false
		result.Issues = append(result.Issues, "Certificate verification failed: "+err.Error())
	}

	return result, nil
}

func getSSLProtocolVersion(version uint16) string {
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
