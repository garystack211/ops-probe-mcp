# ops-probe-mcp

一个基于 [MCP (Model Context Protocol)](https://modelcontextprotocol.io/) 的运维探测服务器，使用 Go 编写。提供域名、DNS、SSL、HTTP 健康检查、Ping、WHOIS 等运维探测工具，可直接集成到 Claude Code、Claude Desktop 等 AI 助手中。

## 功能特性

| 工具 | 说明 |
|------|------|
| `check_domain` | 对单个域名进行 HTTP/DNS/SSL/端口全面检测 |
| `check_domains_batch` | 批量检测多个域名的可用性 |
| `dns_lookup` | 查询域名的 A/AAAA/CNAME/MX/TXT/NS 等 DNS 记录 |
| `ssl_check` | 检查 SSL 证书有效性、颁发机构、过期时间 |
| `ping_check` | Ping 主机检测网络连通性与延迟 |
| `http_health_check` | 检查 HTTP 端点健康状态（支持状态码和内容匹配） |
| `whois_check` | 查询域名 WHOIS 注册信息 |

## 快速开始

### 前置要求

- Go 1.23+
- （可选）`ping` 命令（系统自带）
- （可选）`whois` 命令（`brew install whois` 或 `apt install whois`）

### 编译运行

```bash
# 克隆项目
git clone https://github.com/garystack211/ops-probe-mcp.git
cd ops-probe-mcp

# 编译
go build -o ops-probe-mcp

# 运行（HTTP 模式，默认端口 8080）
./ops-probe-mcp -port 8080

# 运行（stdio 模式，用于本地 MCP 客户端）
./ops-probe-mcp -stdio
```

### 交叉编译（Linux）

```bash
chmod +x build.sh
./build.sh
# 生成 ops-probe-mcp-linux-amd64
```

## 集成到 Claude Code

### 远程 HTTP 模式

在 `.claude/settings.json` 或项目的 MCP 配置中添加：

```json
{
  "mcpServers": {
    "ops-probe": {
      "type": "http",
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

### 本地 stdio 模式

```json
{
  "mcpServers": {
    "ops-probe": {
      "type": "stdio",
      "command": "/path/to/ops-probe-mcp",
      "args": ["-stdio"]
    }
  }
}
```

## 工具使用示例

### check_domain — 单域名全面检测

```json
{
  "name": "check_domain",
  "arguments": {
    "domain": "example.com",
    "check_http": true,
    "check_dns": true,
    "check_ssl": true,
    "check_ports": false,
    "timeout": 10
  }
}
```

### check_domains_batch — 批量检测

```json
{
  "name": "check_domains_batch",
  "arguments": {
    "domains": "example.com,google.com,github.com",
    "check_http": true,
    "check_ssl": true,
    "timeout": 10
  }
}
```

### dns_lookup — DNS 查询

```json
{
  "name": "dns_lookup",
  "arguments": {
    "domain": "example.com"
  }
}
```

### ssl_check — SSL 证书检查

```json
{
  "name": "ssl_check",
  "arguments": {
    "domain": "example.com"
  }
}
```

### ping_check — Ping 检测

```json
{
  "name": "ping_check",
  "arguments": {
    "host": "8.8.8.8",
    "count": 4,
    "timeout": 5
  }
}
```

### http_health_check — HTTP 健康检查

```json
{
  "name": "http_health_check",
  "arguments": {
    "url": "https://example.com/health",
    "expected_status": 200,
    "expected_content": "ok",
    "timeout": 10
  }
}
```

### whois_check — WHOIS 查询

```json
{
  "name": "whois_check",
  "arguments": {
    "domain": "example.com"
  }
}
```

## 项目结构

```
ops-probe-mcp/
├── main.go              # 入口：注册工具、启动服务器（HTTP/stdio 双模式）
├── tools/
│   ├── tools.go         # MCP 工具定义（参数 Schema）
│   └── handlers.go      # 工具请求路由处理
├── checker/
│   ├── dns_checker.go   # DNS 查询实现
│   ├── ssl_checker.go   # SSL 证书检查实现
│   ├── http_checker.go  # HTTP 检查实现
│   ├── health_checker.go# HTTP 健康检查实现
│   ├── ping_checker.go  # Ping 检查实现
│   ├── port_checker.go  # 端口扫描实现
│   └── whois_checker.go # WHOIS 查询实现
├── handler/
│   ├── domain_handler.go# 域名相关工具处理逻辑
│   └── ops_handler.go   # 运维工具处理逻辑
├── build.sh             # Linux 交叉编译脚本
├── go.mod
└── go.sum
```

## 技术栈

- **[mcp-go](https://github.com/mark3labs/mcp-go)** — MCP 协议 Go 实现
- **[miekg/dns](https://github.com/miekg/dns)** — DNS 查询库
- 支持 **Streamable HTTP** 和 **stdio** 两种传输模式

## License

MIT
