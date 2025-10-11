package SryMCPServer

import(
   "fmt"
   "net/http"
   "encoding/json"
)

// MCP Host 提供的工具定義
type HostTool struct {
   Name        string            `json:"name"`
   Description string            `json:"description"`
   Parameters  map[string]string `json:"parameters,omitempty"`
}

// MCP Host 的能力描述
type HostCapabilities struct {
   Version  string	`json:"version"`
   ServerID string	`json:"server_id"`
   Tools    []HostTool	`json:"tools"`
}

// MCP Host結構
type MCPHost struct {
   ID           string			`json:"id"` // Server ID
   Name         string			`json:"name"` // Server名稱
   Capabilities HostCapabilities	`json:"capabilities"` // Server能力描述
   Endpoint     string			`json:"endpoint",oomitempty` // Server的API端點
   IsRelatedPrompt string		`json:"isRelatedPrompt",omitempty` // 是否與ID服務事項相關
   ProcessPrompt string			`json:"processPrompt",omitempty` // 處理ID服務事項的提示，若是則需要做何處理
}

// 處理工具列表請求
func(s *MCPServer) ToolsListFromWeb(w http.ResponseWriter, r *http.Request) {
   // 設定 Response Header
   w.Header().Set("Content-Type", "application/json")
   name := r.PathValue("toolName")
   tool, ok := s.Tools[name]
   if !ok {
      fmt.Println(name + " 並未在 Tools 中")
      return
   }
   // 直接將 struct 編碼並寫入 response
   json.NewEncoder(w).Encode(tool)
   /*
   return MCPMessage{
      JSONRPC: "2.0",
      ID:      msg.ID,
      Result: map[string]interface{}{
         "tools": []map[string]interface{}{
         // 列出所有工具資訊
         },
      },
   }
   */
}
