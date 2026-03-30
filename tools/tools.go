package tools

import "github.com/mark3labs/mcp-go/mcp"

func NewCheckDomainTool() mcp.Tool {
	return mcp.NewTool("check_domain",
		mcp.WithDescription("Check a single domain for HTTP, DNS, SSL, and port availability"),
		mcp.WithString("domain",
			mcp.Required(),
			mcp.Description("The domain name to check (e.g., example.com)")),
		mcp.WithBoolean("check_http",
			mcp.Description("Enable HTTP/HTTPS check")),
		mcp.WithBoolean("check_dns",
			mcp.Description("Enable DNS lookup")),
		mcp.WithBoolean("check_ssl",
			mcp.Description("Enable SSL certificate check")),
		mcp.WithBoolean("check_ports",
			mcp.Description("Enable port scanning")),
		mcp.WithNumber("timeout",
			mcp.Description("Request timeout in seconds")),
	)
}

func NewCheckDomainsBatchTool() mcp.Tool {
	return mcp.NewTool("check_domains_batch",
		mcp.WithDescription("Check multiple domains at once for HTTP, DNS, SSL availability"),
		mcp.WithString("domains",
			mcp.Required(),
			mcp.Description("List of domain names to check")),
		mcp.WithBoolean("check_http",
			mcp.Description("Enable HTTP/HTTPS check")),
		mcp.WithBoolean("check_dns",
			mcp.Description("Enable DNS lookup")),
		mcp.WithBoolean("check_ssl",
			mcp.Description("Enable SSL certificate check")),
		mcp.WithNumber("timeout",
			mcp.Description("Request timeout in seconds")),
	)
}

func NewDNSLookupTool() mcp.Tool {
	return mcp.NewTool("dns_lookup",
		mcp.WithDescription("Perform DNS lookup for a domain and return all DNS records"),
		mcp.WithString("domain",
			mcp.Required(),
			mcp.Description("The domain name to lookup")),
	)
}

func NewSSLCheckTool() mcp.Tool {
	return mcp.NewTool("ssl_check",
		mcp.WithDescription("Check SSL certificate information for a domain"),
		mcp.WithString("domain",
			mcp.Required(),
			mcp.Description("The domain name to check SSL certificate")),
	)
}

func NewPingCheckTool() mcp.Tool {
	return mcp.NewTool("ping_check",
		mcp.WithDescription("Ping a host to check network connectivity and latency"),
		mcp.WithString("host",
			mcp.Required(),
			mcp.Description("The host to ping (IP or domain)")),
		mcp.WithNumber("count",
			mcp.Description("Number of ping packets (default: 4)")),
		mcp.WithNumber("timeout",
			mcp.Description("Timeout in seconds (default: 5)")),
	)
}

func NewHealthCheckTool() mcp.Tool {
	return mcp.NewTool("http_health_check",
		mcp.WithDescription("Check HTTP endpoint health status. Supports connecting to a specific IP to bypass DNS (useful for testing origin servers behind CDN)."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL to check")),
		mcp.WithNumber("expected_status",
			mcp.Description("Expected HTTP status code (default: 200-399)")),
		mcp.WithString("expected_content",
			mcp.Description("Expected content in response body")),
		mcp.WithNumber("timeout",
			mcp.Description("Timeout in seconds (default: 10)")),
		mcp.WithString("resolve_ip",
			mcp.Description("Force TCP connection to this IP address instead of resolving the hostname via DNS. The Host header and TLS SNI still use the original domain. Useful for testing origin servers behind a CDN (e.g., '1.2.3.4' or '1.2.3.4:8080').")),
	)
}

func NewWhoisCheckTool() mcp.Tool {
	return mcp.NewTool("whois_check",
		mcp.WithDescription("Query domain WHOIS information"),
		mcp.WithString("domain",
			mcp.Required(),
			mcp.Description("The domain name to query")),
	)
}
