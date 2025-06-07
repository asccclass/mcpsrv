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
