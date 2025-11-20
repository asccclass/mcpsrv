package weatherMCPServer

import(
   "fmt"
   "time"
   "net/http"
   "github.com/asccclass/mcpsrv/libs/mcpserver"
)

// TodoMCPServer MCP Server 結構
type WeatherMCPServer struct {
   httpClient *http.Client
   API	string
   WeatherData []WeatherData
}

// 註冊工具
func(app *WeatherMCPServer) AddTools(s *SryMCPServer.MCPServer) {
   prompt := fmt.Sprintf(`你是世界頂尖的語意分類專家，專門負責判斷「自然語言句子是否為天氣相關的詢問」。
你將根據語句內容與語氣特徵，精確分類其是否與「天氣詢問」相關。

### 你的任務是：

- **判斷每個輸入句子是否為一個「與天氣相關的詢問」**
- 回答格式必須為 "是" 或 "否"
- 若回答為「是」，則句子**必須同時**符合以下兩個條件：
  1. **主題明確與天氣有關**（如氣溫、下雨、晴天、颱風、濕度等）
  2. **句子語氣為詢問句**（含問句形式或帶有詢問意圖）

### Chain of Thoughts 思維鏈條如下：

1. **理解句子語意**：
   - 讀取句子，辨識句中是否存在「詢問性語氣」（如疑問詞：「會不會」「怎麼樣」「是否」「幾度」「為什麼」等）

2. **判斷主題是否為天氣**：
   - 檢查是否包含與氣象有關的名詞或動詞（如：下雨、晴天、氣溫、濕度、天氣、颱風、氣象、風速等）

3. **同時滿足兩個條件時回答「是」**，否則回答「否」

4. **考慮邊界案例**：
   - 若是與天氣有關但非詢問（如「今天下雨。」），回答「否」
   - 若是詢問但非天氣主題（如「你今天去哪裡？」），回答「否」

### 如果是待辦事項相關：

- 查看/列出氣候狀況 -> action: "get_weather"
- 查看特定地點天氣 -> action: "get_weather_by_city", parameters: {"city": "string"}
- 查看特定時間與地點天氣 -> action: "get_weather_by_city", parameters: {"city": "string", "date": "YYYY-MM-DD HH:MM:SS"}

### 回答格式：

請依據使用者輸入，回應一個 JSON 格式，包含以下欄位：
- "is_related": true/false (問題是否與天氣相關)
- "action": "get_weather" | "get_weather_by_city" | "general_chat" (動作類型)
- "parameters": {} (相關參數，如果有的話)
如果不是待辦事項相關 -> action: "general_chat"
請只回應 JSON，不要其他文字。

### What Not To Do：

- **切記：**
  - **不要**回應除了「是」或「否」以外的文字（如「我不確定」「看起來是」）
  - **不要**因為出現天氣詞彙但不是問句就回答「是」
    - 錯誤示例：「今天晴天。」→ 應為「否」
  - **不要**將「假設性語氣」或「未明確詢問」當成詢問句處理
    - 錯誤示例：「如果下雨怎麼辦。」→ 不夠明確，應回答「否」
  - **不要**因為出現問號就直接認定為詢問（需檢查主題）
`)
   capz := &SryMCPServer.HostCapabilities {
      Version: "",
      ServerID: "0",
   }

   mcphost := &SryMCPServer.MCPHost {
      ID: "Weather001",
      Name: "WEATHER",
      IsRelatedPrompt: prompt,
      ProcessPrompt: "",
   }
   // 取得目前天氣
   capz.Tools = append(capz.Tools, SryMCPServer.HostTool {
      Name: "get_weather",
      Description: "取得天氣狀態",
      Parameters: make(map[string]string),
   })
   s.AddTool(SryMCPServer.Tool{
      Name:        "get_weather",
      Description: "取得天氣狀態",
      InputSchema: SryMCPServer.InputSchema{
         Type:       "object",
         Properties: map[string]SryMCPServer.PropertySchema{},
      },
      Handler: app.getWeather,
   })
   // 取得特定天候
   capz.Tools = append(capz.Tools, SryMCPServer.HostTool {
      Name: "get_weather_by_city",
      Description: "取得特定城市的天氣",
      Parameters: map[string]string{
         "city": "string",
	 "date": "datetime",
      },
   })
   s.AddTool(SryMCPServer.Tool{
      Name:        "get_weather_by_city",
      Description: "取得特定城市的天氣",
      InputSchema: SryMCPServer.InputSchema{
         Type: "object",
         Properties: map[string]SryMCPServer.PropertySchema {
            "city": {
               Type:        "string",
               Description: "城市名稱",
            },
            "date": {
               Type:        "datetime",
               Description: "特定時間",
            },
         },
         Required: []string{"city"},
      },
      Handler: app.getWeatherByCity,
   })
   mcphost.Capabilities = *capz
   s.Tools["weather"] = mcphost
}

// NewTodoMCPServer 建立新的 Weather MCP Server
func NewWeatherMCPServer(apiurl string)(*WeatherMCPServer, error) {
   if apiurl == "" {
      return nil, fmt.Errorf("API url endpoint is empty")
   }
   w := &WeatherMCPServer{
      httpClient: &http.Client{
         Timeout: 30 * time.Second,
      },
      API: apiurl,
   }
   if err := w.ReadWeatherJson(); err != nil {
      fmt.Println("Read weather error:%s", err.Error())
   }
   return w, nil
}
