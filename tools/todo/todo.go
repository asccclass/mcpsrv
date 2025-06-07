package todoMCPServer

import(
   "fmt"
   "time"
   "net/http"
   "github.com/asccclass/mcpsrv/libs/mcpserver"
)

// TodoMCPServer MCP Server 結構
type TodoMCPServer struct {
   httpClient *http.Client
   API	string
}

// 註冊工具
func(app *TodoMCPServer) AddTools(s *SryMCPServer.MCPServer) {
   // get_all_todos)註冊取得所有待辦事項的工具
   s.AddTool(SryMCPServer.Tool{
      Name:        "get_all_todos",
      Description: "取得所有待辦事項",
      InputSchema: SryMCPServer.InputSchema{
         Type:       "object",
         Properties: map[string]SryMCPServer.PropertySchema{},
      },
      Handler: app.getAllTodos,
   })
   // 註冊根據ID取得待辦事項的工具
   s.AddTool(SryMCPServer.Tool{
      Name:        "get_todo_by_id",
      Description: "根據ID取得特定的待辦事項",
      InputSchema: SryMCPServer.InputSchema{
         Type: "object",
         Properties: map[string]SryMCPServer.PropertySchema {
            "id": {
               Type: "number",
               Description: "待辦事項的ID",
            },
         },
         Required: []string{"id"},
      },
      Handler: app.getTodoByID,
   })
   // 註冊建立待辦事項的工具
   s.AddTool(SryMCPServer.Tool{
      Name:        "create_todo",
      Description: "建立新的待辦事項",
      InputSchema: SryMCPServer.InputSchema{
         Type: "object",
         Properties: map[string]SryMCPServer.PropertySchema {
            "context": {
               Type:        "string",
               Description: "待辦事項的內容",
            },
            "user": {
               Type:        "string",
               Description: "使用者名稱",
            },
            "duetime": {
               Type:        "string",
               Description: "到期時間 (可選)",
            },
            "isFinish": {
               Type:        "string",
               Description: "是否完成 (可選，預設為 false)",
            },
         },
         Required: []string{"context", "user"},
      },
      Handler: app.createTodo,
   })
   // 註冊更新待辦事項的工具
   s.AddTool(SryMCPServer.Tool{
      Name:        "update_todo",
      Description: "更新待辦事項",
      InputSchema: SryMCPServer.InputSchema { 
         Type: "object",
         Properties: map[string]SryMCPServer.PropertySchema {
            "id": {
               Type:        "string",
               Description: "待辦事項的ID",
            },
            "context": {
               Type:        "string",
               Description: "待辦事項的內容 (可選)",
            },
            "user": {
               Type:        "string",
               Description: "使用者名稱 (可選)",
            },
            "duetime": {
               Type:        "string",
               Description: "到期時間 (可選)",
            },
            "isFinish": {
               Type:        "string",
               Description: "是否完成 (可選)",
            },
         },
         Required: []string{"id"},
     },
     Handler: app.updateTodo,
   })
   // 註冊刪除待辦事項的工具
   s.AddTool(SryMCPServer.Tool{
      Name:        "delete_todo",
      Description: "刪除待辦事項",
      InputSchema: SryMCPServer.InputSchema {
         Type: "object",
         Properties: map[string]SryMCPServer.PropertySchema{
            "id": {
               Type:        "number",
               Description: "要刪除的待辦事項ID",
            },
         },
         Required: []string{"id"},
      },
      Handler: app.deleteTodo,
   })
}

// NewTodoMCPServer 建立新的 Todo MCP Server
func NewTodoMCPServer(todoapiurl string)(*TodoMCPServer, error) {
   if todoapiurl == "" {
      return nil, fmt.Errorf("todo api url endpoint is empty")
   }
   return &TodoMCPServer{
      httpClient: &http.Client{
         Timeout: 30 * time.Second,
      },
      API: todoapiurl,
   }, nil
}
