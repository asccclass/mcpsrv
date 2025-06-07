package SryMCPServer

import (
   "fmt"
   "time"
   "sync"
   "net/http"
   "encoding/json"
)

type MCPError struct {
   Code    int    `json:"code"`
   Message string `json:"message"`
   Data    interface{}  `json:"data,omitempty"`
}

// MCP Message 結構
type MCPMessage struct {
   JSONRPC string      `json:"jsonrpc"`
   Method  string      `json:"method,omitempty"`
   Params  interface{} `json:"params,omitempty"`
   ID      *string     `json:"id,omitempty"`
   Result  interface{} `json:"result,omitempty"`
   Error   *MCPError   `json:"error,omitempty"`
}

// SSE 客戶端連接
type SSEClient struct {
   Username    string
   Writer      http.ResponseWriter
   Flusher     http.Flusher
   Done        chan struct{}
   MessageChan chan MCPMessage
}

// MCP Server 結構
type MCPServer struct {
   clients    map[string]*SSEClient // username -> SSEClient
   clientsMu  sync.RWMutex
   requests   map[string]chan MCPMessage // 請求ID -> 響應通道
   requestsMu sync.RWMutex
   ToolKits   map[string]Tool
}

// 發送 SSE 資訊
func(s *MCPServer) sendSSEMessage(client *SSEClient, msg MCPMessage) (error) {
   data, err := json.Marshal(msg)
   if err != nil {
      return err
   }
   _, err = fmt.Fprintf(client.Writer, "data: %s\n\n", data)
   if err != nil {
      return err
   }
   client.Flusher.Flush()
   return nil
}

// SSE 連接處理
func(s *MCPServer) sseHandler(w http.ResponseWriter, r *http.Request) {
	/*
   // 從 header 或 query parameter 獲取 JWT token
   var tokenString string
   authHeader := r.Header.Get("Authorization")
   if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
      tokenString = strings.TrimPrefix(authHeader, "Bearer ")
   } else {
      tokenString, _, _ = generateToken("admin")   // r.URL.Query().Get("token")
      fmt.Println(tokenString)
   }
      tokenString, _, _ = generateToken("admin")   // r.URL.Query().Get("token")
   if tokenString == "" {
      http.Error(w, "Missing authentication token", http.StatusUnauthorized)
      return
   }

   // 驗證 token
   claims, err := validateToken(tokenString)
   if err != nil {
      fmt.Println("Unauthorized token")
      http.Error(w, "Invalid token", http.StatusUnauthorized)
      return
   }
*/
   // 使用匿名結構test
    claims := struct {
        Username string
    }{
        Username: "andyliu", // 初始化屬性
    }

   // 設置 SSE headers
   w.Header().Set("Content-Type", "text/event-stream")
   w.Header().Set("Cache-Control", "no-cache")
   w.Header().Set("Connection", "keep-alive")
   w.Header().Set("Access-Control-Allow-Origin", "*")
   w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
   flusher, ok := w.(http.Flusher)
   if !ok {
      http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
      return
   }
   // 創建 SSE 客戶端
   client := &SSEClient{
      Username:    claims.Username,
      Writer:      w,
      Flusher:     flusher,
      Done:        make(chan struct{}),
      MessageChan: make(chan MCPMessage, 100),
   }
   // 註冊客戶端
   s.clientsMu.Lock()
   s.clients[claims.Username] = client
   s.clientsMu.Unlock()
   // 清理函數
   defer func() {
      s.clientsMu.Lock()
      delete(s.clients, claims.Username)
      s.clientsMu.Unlock()
      close(client.Done)
      close(client.MessageChan)
      fmt.Printf("Client %s disconnected", claims.Username)
   }()
   fmt.Printf("Client %s connected via SSE", claims.Username)
   // 發送歡迎消息
   welcomeMsg := MCPMessage{
      JSONRPC: "2.0",
      Method:  "notification/welcome",
      Params: map[string]interface{}{
         "message": fmt.Sprintf("Welcome, %s!", claims.Username),
         "server":  "MCP Server v1.0",
      },
   }
   s.sendSSEMessage(client, welcomeMsg)
   // 監聽客戶端斷開連接
   notify := w.(http.CloseNotifier).CloseNotify()
   // 消息發送循環
   for {
      select {
         case <-notify:
            return
         case <-client.Done:
            return
         case msg := <-client.MessageChan:
            if err := s.sendSSEMessage(client, msg); err != nil {
               fmt.Printf("Failed to send message to %s: %v", claims.Username, err)
               return
            }
         case <-time.After(30 * time.Second): // 發送心跳
            heartbeat := MCPMessage{
               JSONRPC: "2.0",
               Method:  "notification/heartbeat",
               Params: map[string]interface{}{
                  "timestamp": time.Now().Unix(),
               },
            }
            if err := s.sendSSEMessage(client, heartbeat); err != nil {
               return
            }
      }
   }
}

// HTTP POST 請求處理 (用於客戶端發送請求)
func(s *MCPServer) requestHandler(w http.ResponseWriter, r *http.Request) {
/*
   var tokenString string  // 驗證 token
   tokenString, _, _ = generateToken("admin")   // r.URL.Query().Get("token")
   authHeader := r.Header.Get("Authorization")
   if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
      tokenString = strings.TrimPrefix(authHeader, "Bearer ")
   }
   if tokenString == "" {
      http.Error(w, "Missing authentication token", http.StatusUnauthorized)
      return
   }
   claims, err := validateToken(tokenString)
   if err != nil {
      http.Error(w, "Invalid token", http.StatusUnauthorized)
      return
   }
*/
   // 使用匿名結構
    claims := struct {
        Username string
    }{
        Username: "andyliu", // 初始化屬性
    }
   // 解析請求
   var msg MCPMessage
   if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
      fmt.Println("Invalid request body")
      http.Error(w, "Invalid request body", http.StatusBadRequest)
      return
   }
   fmt.Printf("Received request from %s: %+v", claims.Username, msg)
   // 處理請求
   response := s.handleMessage(claims.Username, &msg)
   // 返回響應
   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(response)
}

// 處理 ping
func(s *MCPServer) handlePing(msg *MCPMessage)(MCPMessage) {
   return MCPMessage{
      JSONRPC: "2.0",
      ID:      msg.ID,
      Result:  "pong",
   }
}

// 廣播消息到所有客戶端
func(s *MCPServer) broadcastToAll(msg MCPMessage) {
   s.clientsMu.RLock()
   defer s.clientsMu.RUnlock()

   for username, client := range s.clients {
      select {
         case client.MessageChan <- msg:
            fmt.Printf("Broadcasted message to %s", username)
         default:
            fmt.Printf("Failed to broadcast to %s - channel full", username)
      }
   }
}

// 創建錯誤響應
func(s *MCPServer) createErrorResponse(id *string, code int, message string)(MCPMessage) {
   return MCPMessage{
      JSONRPC: "2.0",
      ID:      id,
      Error: &MCPError{
         Code:    code,
         Message: message,
      },
   }
}

// 處理初始化
func(s *MCPServer) handleInitialize(msg *MCPMessage) MCPMessage {
   return MCPMessage{
      JSONRPC: "2.0",
      ID:      msg.ID,
      Result: map[string]interface{}{
         "protocolVersion": "2024-11-05",
         "capabilities": map[string]interface{}{
            "tools": map[string]interface{}{
            "listChanged": true,
            },
         },
         "serverInfo": map[string]interface{}{
            "name":    "MCP Demo Server",
            "version": "1.0.0",
         },
      },
   }
}

// 處理工具列表請求
func(s *MCPServer) handleToolsList(msg *MCPMessage)(MCPMessage) {
   return MCPMessage{
      JSONRPC: "2.0",
      ID:      msg.ID,
      Result: map[string]interface{}{
         "tools": []map[string]interface{}{
         // 列出所有工具資訊
         },
      },
   }
}

// 處理工具調用
func(s *MCPServer) handleToolsCall(msg *MCPMessage)(MCPMessage) {
   params, ok := msg.Params.(map[string]interface{})
   if !ok {
      return s.createErrorResponse(msg.ID, -32602, "Invalid params")
   }
   toolName, ok := params["name"].(string)
   if !ok {
      return s.createErrorResponse(msg.ID, -32602, "Missing tool name")
   }

   // 解析工具參數
   var toolArgs map[string]interface{}
   if args, exists := params["arguments"]; exists {
      if argsMap, ok := args.(map[string]interface{}); ok {
         toolArgs = argsMap
      } else {
         return s.createErrorResponse(msg.ID, -32602, "Invalid tool arguments")
      }
   } else {
      toolArgs = make(map[string]interface{})
   }
   // 檢查工具是否存在
   tool, exists := s.ToolKits[toolName]
   if !exists {
      return s.createErrorResponse(msg.ID, -32601, fmt.Sprintf("Unknown tool: %s", toolName))
   }
   // 驗證必需的參數
   if err := s.validateToolArguments(tool, toolArgs); err != nil {
      return s.createErrorResponse(msg.ID, -32602, fmt.Sprintf("Invalid arguments: %s", err.Error()))
   }
   // 執行工具
   result, err := s.executeTool(toolName, toolArgs)
   if err != nil {
      return s.createErrorResponse(msg.ID, -32603, fmt.Sprintf("Tool execution failed: %s", err.Error()))
   }
   // 返回成功響應
   return MCPMessage{
      JSONRPC: "2.0",
      ID:      msg.ID,
      Result: map[string]interface{}{
         "content": result.Content,
         "isError": false,
      },
   }
}

// 處理 MCP 消息
func(s *MCPServer) handleMessage(username string, msg *MCPMessage) (MCPMessage) {
   switch msg.Method {
      case "initialize":
         return s.handleInitialize(msg)
      case "tools/list":
         return s.handleToolsList(msg)
      case "tools/call":
         return s.handleToolsCall(msg)
      case "ping":
         return s.handlePing(msg)
      default: // 未知方法
         return MCPMessage{
            JSONRPC: "2.0",
            ID:      msg.ID,
            Error: &MCPError{
               Code:    -32601,
               Message: "Method not found",
            },
         }
   }
}

// 註冊工具
func(app *MCPServer) RegisterTool(tool Tool) { 
   app.ToolKits[tool.Name] = tool
}
// alias
func(app *MCPServer) AddTool(tool Tool) {
   app.RegisterTool(tool)
}

func(app *MCPServer) AddRouter(router *http.ServeMux) {
   // router.Handle("POST /auth", http.HandlerFunc(app.authHandler))
   router.Handle("GET /sse", http.HandlerFunc(app.sseHandler))
   router.Handle("POST /request", http.HandlerFunc(app.requestHandler))
}

func NewMCPServer()(*MCPServer) {
   return &MCPServer{
      clients:  make(map[string]*SSEClient),
      requests: make(map[string]chan MCPMessage),
      ToolKits: make(map[string]Tool),
   }
}
