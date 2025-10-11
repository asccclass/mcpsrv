# MCPsrv - MCP Server

[![Go Report Card](https://goreportcard.com/badge/github.com/asccclass/mcpsrv)](https://goreportcard.com/report/github.com/asccclass/mcpsrv)
[![License](https://img.shields.io/github/license/asccclass/mcpsrv)](https://github.com/asccclass/mcpsrv/blob/main/LICENSE)

MCPsrv is a lightweight MCP (Management Control Protocol) server implementation written in Go. It provides a flexible framework for handling MCP commands and managing connections.

## Features

- Lightweight and efficient MCP server implementation
- Built-in support for common MCP commands
- Support for custom command handlers
- Simple configuration through environment variables
- Connection management with timeouts

## Installation

### 目錄架構
```
.gitignore: Git 忽略檔案規則。
Dockerfile: 定義如何構建 Docker 映像檔。
README.md: 專案的說明文件。
auth.go: 與認證相關的功能。
clean.sh: 清理腳本。
envfile.example: 環境變數設定範例檔案。
go.mod: Go modules 設定。
go.sum: Go modules 的相依性摘要。
libs/: 自定義函式庫的目錄。
makefile: 編譯和自動化相關的指令。
router.go: 負責路由的功能。
server.go: 伺服器主程式。
tls.go: 與 TLS 加密協定相關的功能。
token.go: 與 Token 處理相關的功能。
tools/: 工具程式的目錄。
www/: 靜態網頁檔案的目錄。
```

### Prerequisites

- Go 1.16 or higher

### Using Go Get

```bash
go get github.com/asccclass/mcpsrv
```

### Clone and Build

```bash
git clone https://github.com/asccclass/mcpsrv.git
cd mcpsrv
go build
```

## Usage

### Starting the Server

```go
package main

import (
    "github.com/asccclass/mcpsrv/server"
)

func main() {
    // Create a new server instance
    srv := server.NewServer()
    
    // Start the server
    srv.Start()
}
```

### Environment Configuration

MCPsrv can be configured using environment variables:

- `MCP_HOST`: Host address to bind (default: "localhost")
- `MCP_PORT`: Port to listen on (default: "8080")
- `MCP_TIMEOUT`: Connection timeout in seconds (default: 60)

### Adding Custom Command Handlers

You can extend MCPsrv with your own command handlers:

```go
package main

import (
    "github.com/asccclass/mcpsrv/server"
    "github.com/asccclass/mcpsrv/handler"
)

func main() {
    srv := server.NewServer()
    
    // Register a custom command handler
    srv.RegisterHandler("CUSTOM_CMD", func(cmd *handler.Command) *handler.Response {
        return handler.NewResponse(200, "Custom command executed")
    })
    
    srv.Start()
}
```

## API Documentation

### Server

The `server` package provides the main server functionality:

- `NewServer()`: Creates a new MCP server instance
- `Start()`: Starts the server and listens for incoming connections
- `RegisterHandler(cmd string, handler HandlerFunc)`: Registers a custom command handler

### Handler

The `handler` package contains command handling utilities:

- `Command`: Represents an MCP command with parameters
- `Response`: Represents a response to an MCP command
- `NewResponse(code int, message string)`: Creates a new response

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Contact
Project maintainer: [asccclass](https://github.com/asccclass)
---
## 參考構思

```
// 設定 agent 參考資訊
agent.create_guideline(
   condition=”Customer asks about refunds”,
   action=”Check order status first to see if eligible”,
   tools=[check_order_status],
)
```
