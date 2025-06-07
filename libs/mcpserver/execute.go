package SryMCPServer

import (
   "fmt"
   "github.com/mark3labs/mcp-go/mcp"
)

// 執行工具
func(s *MCPServer) executeTool(toolName string, args map[string]interface{}) (*mcp.CallToolResult, error) {
   tool, exists := s.ToolKits[toolName]
   if !exists {
      return mcp.NewToolResultText(""), fmt.Errorf("tool not found: %s", toolName)
   }
   // 設置預設值
   for propName, propSchema := range tool.InputSchema.Properties {
      if _, exists := args[propName]; !exists && propSchema.Default != nil {
         args[propName] = propSchema.Default
      }
   }
   // 執行工具處理函數
   if tool.Handler != nil {
      return tool.Handler(args)
   }
   return mcp.NewToolResultText(""), fmt.Errorf("no handler for tool: %s", toolName)
}

