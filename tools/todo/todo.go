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
   prompt := fmt.Sprintf(`你是一個待辦事項助手。請分析以下使用者輸入，判斷是否與待辦事項相關。

請依據使用者輸入，回應一個 JSON 格式，包含以下欄位：
- "is_related": true/false (是否與待辦事項相關)
- "action": "get_all_todos" | "get_todo_by_id" | "create_todo" | "update_todo" | "delete_todo" | "general_chat" (動作類型)
- "parameters": {} (相關參數，如果有的話)

如果是待辦事項相關：
- 查看/列出待辦事項 -> action: "get_all_todos"
- 查看特定待辦事項 -> action: "get_todo_by_id", parameters: {"id": "string"}
- 新增/建立待辦事項 -> action: "create_todo", parameters: {"context": "內容", "user": "負責者", "isFinished": "0"}
- 修改/更新待辦事項 -> action: "update_todo", parameters: {"id": "string", 其他要更新的欄位...}
- 刪除待辦事項 -> action: "delete_todo", parameters: {"id": "string"}
- 若欄位名稱為status/狀態時，請將欄位名稱替換為isFinish欄位：
   - 當status/狀態為「未處理」、「未分類」、「未知」或類似狀態時，isFinish應設為 "0"
   - 當status/狀態為「進行中」、「待處理」、「審核中」、「處理中/open/process/in progress」或類似狀態時，isFinish應設為 "1"
   - 當status/狀態為「完成/completed/done」、「已結束/closed」、「已處理/finished」或類似狀態時，isFinish應設為 "2"
   - 當status/狀態為「擱置/rejected」、「不處理」、「暫存/pause/suspended/pending」或類似狀態時，isFinish應設為 "3"
- 當使用者輸入的文字包含「時間相關語句」（如「明天上午十點」、「下週三下午三點」等），請根據「台北當前日期」進行解析，並將該時間轉換為 YYYY-MM-DD HH:MM:SS的時間格式。
   - 例如使用者輸入「明天上午十點開會」，判斷當前時間(如：2025-06-14）為基準，加一天後時間為 2025-06-15 10:00:00，並將此值存為欄位名稱：duetime。
   - 若時間資訊模糊（如「晚上」、「中午」），請估算合理時間（如「晚上」可設為 20:00:00）。
   - 若沒有指定年份、月份、日期，則年月日為當前西元年月日。
   - 輸出要求：
      - 時間應為完整的 UTC 格式時間字串（YYYY-MM-DD HH:MM:SS）。
      - 資料應儲存在欄位 duetime。

如果不是待辦事項相關 -> action: "general_chat"
請只回應 JSON，不要其他文字。
`)
   capz := &SryMCPServer.HostCapabilities {
      Version: "",
      ServerID: "0",
   }

   mcphost := &SryMCPServer.MCPHost {
      ID: "todo001",
      Name: "TODO",
      IsRelatedPrompt: prompt,
      ProcessPrompt: "",
   }
   // get_all_todos)註冊取得所有待辦事項的工具
   capz.Tools = append(capz.Tools, SryMCPServer.HostTool {
      Name: "get_all_todos",
      Description: "取得所有待辦事項",
      Parameters: make(map[string]string),
   })
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
   capz.Tools = append(capz.Tools, SryMCPServer.HostTool {
      Name: "get_todo_by_id",
      Description: "根據ID取得特定的待辦事項",
      Parameters: map[string]string{
         "id": "number",
      },
   })
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
   capz.Tools = append(capz.Tools, SryMCPServer.HostTool {
      Name: "create_todo",
      Description: "建立新的待辦事項",
      Parameters: map[string]string{
         "context": "string",
	 "user": "string",
	 "duetime": "string",
         "isFinish": "string",
      },
   })
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
	       Default: "",
            },
            "isFinish": {
               Type:        "string",
               Description: "事件狀態，是否完成 (可選，預設為 0)",
	       Default: "0",
            },
         },
         Required: []string{"context", "user"},
      },
      Handler: app.createTodo,
   })
   // 註冊更新待辦事項的工具
   capz.Tools = append(capz.Tools, SryMCPServer.HostTool {
      Name: "update_todo",
      Description: "更新待辦事項",
      Parameters: map[string]string{
         "id": "string",
         "context": "string",
         "user": "string",
         "duetime": "string",
         "isFinish": "string",
      },
   })
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
   capz.Tools = append(capz.Tools, SryMCPServer.HostTool {
      Name: "delete_todo",
      Description: "刪除待辦事項",
      Parameters: map[string]string{
         "id": "int",
      },
   })
   s.AddTool(SryMCPServer.Tool{
      Name:        "delete_todo",
      Description: "刪除待辦事項",
      InputSchema: SryMCPServer.InputSchema {
         Type: "object",
         Properties: map[string]SryMCPServer.PropertySchema{
            "id": {
               Type:        "int",
               Description: "要刪除的待辦事項ID",
            },
         },
         Required: []string{"id"},
      },
      Handler: app.deleteTodo,
   })
   mcphost.Capabilities = *capz
   s.Tools["todo"] = mcphost
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
