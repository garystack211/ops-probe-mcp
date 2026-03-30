package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/mark3labs/mcp-go/server"
	"ops-probe-mcp/tools"
)

func main() {
	port := flag.Int("port", 8080, "Port for SSE server")
	flag.Parse()

	// 创建 MCP 服务器
	s := server.NewMCPServer("ops-probe-mcp", "1.0.0")

	// 注册 4 个工具
	s.AddTool(tools.NewCheckDomainTool(), tools.CheckDomainHandler)
	s.AddTool(tools.NewCheckDomainsBatchTool(), tools.CheckDomainsBatchHandler)
	s.AddTool(tools.NewDNSLookupTool(), tools.DNSLookupHandler)
	s.AddTool(tools.NewSSLCheckTool(), tools.SSLCheckHandler)
	s.AddTool(tools.NewSSLCertDetailTool(), tools.SSLCertDetailHandler)

	// 注册新增的运维工具
	s.AddTool(tools.NewPingCheckTool(), tools.PingCheckHandler)
	s.AddTool(tools.NewHealthCheckTool(), tools.HealthCheckHandler)
	s.AddTool(tools.NewWhoisCheckTool(), tools.WhoisCheckHandler)

	// Streamable HTTP 模式（支持 Claude.ai）
	transport := server.NewStreamableHTTPServer(s)

	log.Printf("Starting ops-probe-mcp Streamable HTTP server on port %d", *port)

	if err := transport.Start(":" + strconv.Itoa(*port)); err != nil {
		log.Fatal(err)
	}
}
