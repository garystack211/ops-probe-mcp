package checker

import (
	"os/exec"
	"strings"
	"time"
)

type WhoisResult struct {
	Domain         string    `json:"domain"`
	Registrar      string    `json:"registrar,omitempty"`
	CreationDate   string    `json:"creation_date,omitempty"`
	ExpirationDate string    `json:"expiration_date,omitempty"`
	NameServers    []string  `json:"name_servers,omitempty"`
	Status         []string  `json:"status,omitempty"`
	Available      bool      `json:"available"`
	Error          string    `json:"error,omitempty"`
}

func WhoisCheck(domain string, timeout time.Duration) (*WhoisResult, error) {
	result := &WhoisResult{
		Domain:      domain,
		Available:   false,
		NameServers: []string{},
		Status:      []string{},
	}

	if timeout == 0 {
		timeout = 10 * time.Second
	}

	cmd := exec.Command("whois", domain)
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = err.Error()
		return result, nil
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		lower := strings.ToLower(line)

		if strings.Contains(lower, "no match") || strings.Contains(lower, "not found") {
			result.Available = true
			continue
		}

		if strings.HasPrefix(lower, "registrar:") {
			result.Registrar = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		} else if strings.Contains(lower, "creation date") || strings.Contains(lower, "created") {
			if parts := strings.SplitN(line, ":", 2); len(parts) > 1 {
				result.CreationDate = strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(lower, "expir") {
			if parts := strings.SplitN(line, ":", 2); len(parts) > 1 {
				result.ExpirationDate = strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(lower, "name server") {
			if parts := strings.SplitN(line, ":", 2); len(parts) > 1 {
				ns := strings.TrimSpace(parts[1])
				if ns != "" {
					result.NameServers = append(result.NameServers, ns)
				}
			}
		} else if strings.Contains(lower, "status:") {
			if parts := strings.SplitN(line, ":", 2); len(parts) > 1 {
				status := strings.TrimSpace(parts[1])
				if status != "" {
					result.Status = append(result.Status, status)
				}
			}
		}
	}

	return result, nil
}
