package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/mark3labs/mcp-go/server"
	"ops-probe-mcp/tools"
)

func main() {
	port := flag.Int("port", 8080, "Port for SSE server")
	stdio := flag.Bool("stdio", false, "Use stdio transport")
	flag.Parse()

	// 创建 MCP 服务器
	s := server.NewMCPServer("ops-probe-mcp", "1.0.0")

	// 注册 4 个工具
	s.AddTool(tools.NewCheckDomainTool(), tools.CheckDomainHandler)
	s.AddTool(tools.NewCheckDomainsBatchTool(), tools.CheckDomainsBatchHandler)
	s.AddTool(tools.NewDNSLookupTool(), tools.DNSLookupHandler)
	s.AddTool(tools.NewSSLCheckTool(), tools.SSLCheckHandler)

	// 注册新增的 3 个运维工具
	s.AddTool(tools.NewPingCheckTool(), tools.PingCheckHandler)
	s.AddTool(tools.NewHealthCheckTool(), tools.HealthCheckHandler)
	s.AddTool(tools.NewWhoisCheckTool(), tools.WhoisCheckHandler)

	if *stdio {
		// Stdio 模式
		log.SetOutput(os.Stderr)
		log.Println("Starting ops-probe-mcp in stdio mode")
		transport := server.NewStdioServer(s)
		if err := transport.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
			log.Fatal(err)
		}
	} else {
		// Streamable HTTP 模式（新版，支持 Claude.ai）
		transport := server.NewStreamableHTTPServer(s)

		log.Printf("Starting ops-probe-mcp Streamable HTTP server on port %d", *port)

		if err := transport.Start(":" + strconv.Itoa(*port)); err != nil {
			log.Fatal(err)
		}
	}
}

