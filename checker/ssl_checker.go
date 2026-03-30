package checker

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
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

func SSLCheck(targetURL string, timeout time.Duration) (*SSLResult, error) {
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

	conn, err := tls.Dial("tcp", host+":"+port, &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS10,
	})

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

	if len(state.PeerCertificates) > 1 {
		result.Issues = append(result.Issues, "Certificate chain has multiple certificates")
	}

	chainVerified := verifyCertificateChain(state.PeerCertificates)
	if !chainVerified {
		result.Issues = append(result.Issues, "Certificate chain verification failed")
	}

	certIssuer := ""
	for _, name := range cert.Issuer.Organization {
		certIssuer = name
		break
	}
	result.Issuer = certIssuer

	certSubject := ""
	for _, name := range cert.Subject.Organization {
		certSubject = name
		break
	}
	if certSubject == "" {
		certSubject = cert.Subject.CommonName
	}
	result.Subject = certSubject

	if time.Until(cert.NotAfter).Hours()/24 < 30 {
		result.Issues = append(result.Issues, "Certificate expires within 30 days")
	}

	if time.Until(cert.NotBefore).Hours() > 0 {
		result.Issues = append(result.Issues, "Certificate is not yet valid")
	}

	opts := x509.VerifyOptions{
		DNSName: host,
	}
	_, err = cert.Verify(opts)
	if err != nil {
		result.Issues = append(result.Issues, "Certificate hostname verification failed")
	}

	return result, nil
}

func verifyCertificateChain(certs []*x509.Certificate) bool {
	if len(certs) < 2 {
		return true
	}

	for i := 0; i < len(certs)-1; i++ {
		if certs[i].Subject.CommonName != certs[i+1].Issuer.CommonName {
			return false
		}
	}

	return true
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
