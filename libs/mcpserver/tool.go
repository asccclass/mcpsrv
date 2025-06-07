package SryMCPServer

import(
   "github.com/mark3labs/mcp-go/mcp"
)

type PropertySchema struct {
   Type        string      `json:"type"`
   Description string      `json:"description"`
   Enum        []string    `json:"enum,omitempty"`
   Default     interface{} `json:"default,omitempty"`
}

type InputSchema struct {
   Type       string                     `json:"type"`
   Properties map[string]PropertySchema  `json:"properties"`
   Required   []string                   `json:"required"`
}

// 工具相關結構體
type Tool struct {
   Name        string             `json:"name"`
   Description string             `json:"description"`
   InputSchema InputSchema        `json:"inputSchema"`
   Handler     func(map[string]interface{}) (*mcp.CallToolResult, error) `json:"-"`
}
