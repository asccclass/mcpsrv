package todoMCPServer

import(
   "fmt"
   "time"
   "strconv"
   "encoding/json"
   "github.com/mark3labs/mcp-go/mcp"
)

// Todo 代表 todo 資料表中的一個項目
type Todo struct {
   ID         int       `json:"id"`
   CreateDate time.Time `json:"createdate"`
   Context    string    `json:"context"`
   User       string    `json:"user"`
   DueTime    string    `json:"duetime"`
   IsFinish   string    `json:"isFinish"`
}

// getAllTodos 取得所有待辦事項
func(s *TodoMCPServer) getAllTodos(args map[string]interface{}) (*mcp.CallToolResult, error) {
   url := fmt.Sprintf("%s/todos", s.API)
   resp, err := s.makeAPIRequest("GET", url, nil)
   if err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to get todos: %v", err)), nil
   }

   var todos []Todo
   if err := json.Unmarshal(resp, &todos); err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to parse response: %v", err)), nil
   }
   // 格式化輸出
   result := "# 現有待辦事項如下:\n"
   for _, todo := range todos {
      status := "❌ 未完成"
      if todo.IsFinish == "1" {
         status = "處理中"
      } else if todo.IsFinish == "2" {
         status = "✅ 完成"
      } else if todo.IsFinish == "3" {
         status = "擱置"
      }
      result += fmt.Sprintf("## ID: %d\n* 負責者: %s\n* 事件狀態: %s\n* 到期時間: %s\n", todo.ID, todo.User, status, todo.DueTime)
      result += fmt.Sprintf("* 內容: %s（建立時間: %s）", todo.Context, todo.CreateDate.Format("2006-01-02 15:04:05"))
   }
   return mcp.NewToolResultText(result), nil
}

// getTodoByID 取得特定ID的待辦事項
func(s *TodoMCPServer) getTodoByID(args map[string]interface{}) (*mcp.CallToolResult, error) {
   idArg, ok := args["id"]
   if !ok {
      return mcp.NewToolResultError("Missing required parameter: id"), nil
   }
   id, err := strconv.Atoi(fmt.Sprintf("%v", idArg))
   if err != nil {
      return mcp.NewToolResultError("Invalid id format"), nil
   }
   url := fmt.Sprintf("%s/todo/%d", s.API, id)
   resp, err := s.makeAPIRequest("GET", url, nil)
   if err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to get todo: %v", err)), nil
   }
   var todo Todo
   if err := json.Unmarshal(resp, &todo); err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to parse response: %v", err)), nil
   }
   status := "❌ 未完成"
   if todo.IsFinish == "1" {
      status = "處理中"
   } else if todo.IsFinish == "2" {
      status = "✅ 完成"
   } else if todo.IsFinish == "3" {
      status = "擱置"
   }
   result := fmt.Sprintf("# 📝 待辦事項詳情:\n\n")
   result += fmt.Sprintf("* ID: %d\n", todo.ID)
   result += fmt.Sprintf("* 使用者: %s\n", todo.User)
   result += fmt.Sprintf("* 內容: %s\n", todo.Context)
   result += fmt.Sprintf("* 到期時間: %s\n", todo.DueTime)
   result += fmt.Sprintf("* 狀態: %s\n", status)
   result += fmt.Sprintf("* 建立時間: %s\n", todo.CreateDate.Format("2006-01-02 15:04:05"))
   return mcp.NewToolResultText(result), nil
}

// createTodo 新增待辦事項
func(s *TodoMCPServer) createTodo(args map[string]interface{}) (*mcp.CallToolResult, error) {
   contextArg, ok := args["context"]
   if !ok {
      return mcp.NewToolResultError("Missing required parameter: context"), nil
   }
   userArg, ok := args["user"]
   if !ok {
      return mcp.NewToolResultError("Missing required parameter: user"), nil
   }
   todo := Todo{
      Context:    fmt.Sprintf("%v", contextArg),
      User:       fmt.Sprintf("%v", userArg),
      CreateDate: time.Now(),
      IsFinish:   "false",
   }
   // 可選參數
   if dueTimeArg, ok := args["duetime"]; ok {
      todo.DueTime = fmt.Sprintf("%v", dueTimeArg)
   }
   if isFinishArg, ok := args["isFinish"]; ok {
      todo.IsFinish = fmt.Sprintf("%v", isFinishArg)
   }
   url := fmt.Sprintf("%s/todo", s.API)
   resp, err := s.makeAPIRequest("POST", url, todo)
   if err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to create todo: %v", err)), nil
   }

   var createdTodo Todo
   if err := json.Unmarshal(resp, &createdTodo); err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to parse response: %v", err)), nil
   }
   result := fmt.Sprintf("# ✅ 待辦事項建立成功!\n\n")
   result += fmt.Sprintf("* ID: %d\n", createdTodo.ID)
   result += fmt.Sprintf("* 使用者: %s\n", createdTodo.User)
   result += fmt.Sprintf("* 內容: %s\n", createdTodo.Context)
   result += fmt.Sprintf("* 到期時間: %s\n", createdTodo.DueTime)
   result += fmt.Sprintf("* 建立時間: %s\n", createdTodo.CreateDate.Format("2006-01-02 15:04:05"))
   return mcp.NewToolResultText(result), nil
}

// updateTodo 修改待辦事項
func(s *TodoMCPServer) updateTodo(args map[string]interface{}) (*mcp.CallToolResult, error) {
   idArg, ok := args["id"].(string)
   if !ok {
      return mcp.NewToolResultError("Missing required parameter: id"), nil
   }
   id, err := strconv.Atoi(fmt.Sprintf("%v", idArg))
   if err != nil {
      return mcp.NewToolResultError("Invalid id format"), nil
   }
   // 建立更新用的 todo 物件
   todo := Todo{
      ID: id,
   }
   // 可選參數
   if contextArg, ok := args["context"]; ok {
      todo.Context = fmt.Sprintf("%v", contextArg)
   }
   if userArg, ok := args["user"]; ok {
      todo.User = fmt.Sprintf("%v", userArg)
   }
   if dueTimeArg, ok := args["duetime"]; ok {
      todo.DueTime = fmt.Sprintf("%v", dueTimeArg)
   }
   if isFinishArg, ok := args["isFinish"]; ok {
      todo.IsFinish = fmt.Sprintf("%v", isFinishArg)
   }
   url := fmt.Sprintf("%s/todo/%d", s.API, id)
   fmt.Println(url, todo.User)
   resp, err := s.makeAPIRequest("PUT", url, todo)
   if err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to update todo: %v", err)), nil
   }
   var updatedTodo Todo
   if err := json.Unmarshal(resp, &updatedTodo); err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to parse response: %v", err)), nil
   }
   status := "❌ 未完成"
   if todo.IsFinish == "1" {
      status = "處理中"
   } else if todo.IsFinish == "2" {
      status = "✅ 完成"
   } else if todo.IsFinish == "3" {
      status = "擱置"
   }
   result := fmt.Sprintf("✅ 待辦事項更新成功!\n\n")
   result += fmt.Sprintf("ID: %d\n", updatedTodo.ID)
   result += fmt.Sprintf("內容: %s\n", updatedTodo.Context)
   result += fmt.Sprintf("使用者: %s\n", updatedTodo.User)
   result += fmt.Sprintf("到期時間: %s\n", updatedTodo.DueTime)
   result += fmt.Sprintf("狀態: %s\n", status)
   return mcp.NewToolResultText(result), nil
}

// deleteTodo 刪除待辦事項
func(s *TodoMCPServer) deleteTodo(args map[string]interface{}) (*mcp.CallToolResult, error) {
   idArg, ok := args["id"].(string)
   if !ok {
      return mcp.NewToolResultError("Missing required parameter: id"), nil
   }
   id, err := strconv.Atoi(fmt.Sprintf("%v", idArg))
   if err != nil {
      return mcp.NewToolResultError("Invalid id format"), nil
   }
   url := fmt.Sprintf("%s/todo/%d", s.API, id)
   _, err = s.makeAPIRequest("DELETE", url, nil)
   if err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Failed to delete todo: %v", err)), nil
   }
   result := fmt.Sprintf("🗑️ 待辦事項 ID %d 已成功刪除!", id)
   return mcp.NewToolResultText(result), nil
}
