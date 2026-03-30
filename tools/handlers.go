package tools

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"ops-probe-mcp/handler"
)

func CheckDomainHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()
	args, _ := request.Params.Arguments.(map[string]interface{})
	domain, _ := args["domain"].(string)

	log.Printf("[CHECK_DOMAIN] Request - Domain: %s, Args: %+v", domain, args)

	options := make(map[string]interface{})
	for k, v := range args {
		if k != "domain" {
			options[k] = v
		}
	}

	result, err := handler.CheckDomain(domain, options)
	duration := time.Since(start)

	if err != nil {
		log.Printf("[CHECK_DOMAIN] Error - Domain: %s, Duration: %v, Error: %v", domain, duration, err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	jsonData, _ := json.Marshal(result)
	log.Printf("[CHECK_DOMAIN] Success - Domain: %s, Duration: %v, Result: %s", domain, duration, string(jsonData))
	return mcp.NewToolResultText(string(jsonData)), nil
}

func CheckDomainsBatchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()
	args, _ := request.Params.Arguments.(map[string]interface{})
	domainsRaw, _ := args["domains"].([]interface{})

	log.Printf("[CHECK_DOMAINS_BATCH] Request - Domains count: %d, Args: %+v", len(domainsRaw), args)

	options := make(map[string]interface{})
	for k, v := range args {
		if k != "domains" {
			options[k] = v
		}
	}

	results, err := handler.CheckDomainsBatch(domainsRaw, options)
	duration := time.Since(start)

	if err != nil {
		log.Printf("[CHECK_DOMAINS_BATCH] Error - Duration: %v, Error: %v", duration, err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	jsonData, _ := json.Marshal(map[string]interface{}{
		"results": results,
		"count":   len(results),
	})
	log.Printf("[CHECK_DOMAINS_BATCH] Success - Count: %d, Duration: %v", len(results), duration)
	return mcp.NewToolResultText(string(jsonData)), nil
}

func DNSLookupHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()
	args, _ := request.Params.Arguments.(map[string]interface{})
	domain, _ := args["domain"].(string)

	log.Printf("[DNS_LOOKUP] Request - Domain: %s", domain)

	result, err := handler.DNSLookup(domain, "A")
	duration := time.Since(start)

	if err != nil {
		log.Printf("[DNS_LOOKUP] Error - Domain: %s, Duration: %v, Error: %v", domain, duration, err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	jsonData, _ := json.Marshal(result)
	log.Printf("[DNS_LOOKUP] Success - Domain: %s, Duration: %v, Resolved: %v", domain, duration, result.Resolved)
	return mcp.NewToolResultText(string(jsonData)), nil
}

func SSLCheckHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()
	args, _ := request.Params.Arguments.(map[string]interface{})
	domain, _ := args["domain"].(string)

	log.Printf("[SSL_CHECK] Request - Domain: %s", domain)

	result, err := handler.SSLCertificateCheck(domain)
	duration := time.Since(start)

	if err != nil {
		log.Printf("[SSL_CHECK] Error - Domain: %s, Duration: %v, Error: %v", domain, duration, err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	jsonData, _ := json.Marshal(result)
	log.Printf("[SSL_CHECK] Success - Domain: %s, Duration: %v, Valid: %v, Days: %d", domain, duration, result.Valid, result.DaysRemaining)
	return mcp.NewToolResultText(string(jsonData)), nil
}

func PingCheckHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()
	args, _ := request.Params.Arguments.(map[string]interface{})
	host, _ := args["host"].(string)
	count := 4
	if c, ok := args["count"].(float64); ok {
		count = int(c)
	}
	timeout := 5 * time.Second
	if t, ok := args["timeout"].(float64); ok {
		timeout = time.Duration(t) * time.Second
	}

	log.Printf("[PING_CHECK] Request - Host: %s, Count: %d", host, count)

	result, err := handler.PingCheck(host, count, timeout)
	duration := time.Since(start)

	if err != nil {
		log.Printf("[PING_CHECK] Error - Host: %s, Duration: %v, Error: %v", host, duration, err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	jsonData, _ := json.Marshal(result)
	log.Printf("[PING_CHECK] Success - Host: %s, Duration: %v, Loss: %.1f%%", host, duration, result.PacketLoss)
	return mcp.NewToolResultText(string(jsonData)), nil
}

func HealthCheckHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()
	args, _ := request.Params.Arguments.(map[string]interface{})
	url, _ := args["url"].(string)
	expectedStatus := 0
	if s, ok := args["expected_status"].(float64); ok {
		expectedStatus = int(s)
	}
	expectedContent, _ := args["expected_content"].(string)
	timeout := 10 * time.Second
	if t, ok := args["timeout"].(float64); ok {
		timeout = time.Duration(t) * time.Second
	}

	log.Printf("[HEALTH_CHECK] Request - URL: %s", url)

	result, err := handler.HTTPHealthCheck(url, expectedStatus, expectedContent, timeout)
	duration := time.Since(start)

	if err != nil {
		log.Printf("[HEALTH_CHECK] Error - URL: %s, Duration: %v, Error: %v", url, duration, err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	jsonData, _ := json.Marshal(result)
	log.Printf("[HEALTH_CHECK] Success - URL: %s, Duration: %v, Healthy: %v, Status: %d", url, duration, result.Healthy, result.StatusCode)
	return mcp.NewToolResultText(string(jsonData)), nil
}

func WhoisCheckHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()
	args, _ := request.Params.Arguments.(map[string]interface{})
	domain, _ := args["domain"].(string)

	log.Printf("[WHOIS_CHECK] Request - Domain: %s", domain)

	result, err := handler.WhoisCheck(domain, 10*time.Second)
	duration := time.Since(start)

	if err != nil {
		log.Printf("[WHOIS_CHECK] Error - Domain: %s, Duration: %v, Error: %v", domain, duration, err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	jsonData, _ := json.Marshal(result)
	log.Printf("[WHOIS_CHECK] Success - Domain: %s, Duration: %v, Available: %v", domain, duration, result.Available)
	return mcp.NewToolResultText(string(jsonData)), nil
}
