package main

import(
   "os"
   "fmt"
   "net/http"
   "github.com/asccclass/serverstatus"
   "github.com/asccclass/sherryserver"
   "github.com/asccclass/mcpsrv/tools/todo"
   "github.com/asccclass/mcpsrv/tools/weather"
   "github.com/asccclass/mcpsrv/libs/mcpserver"
)

// Create your Router function
func NewRouter(srv *SherryServer.Server, documentRoot string)(*http.ServeMux) {
   router := http.NewServeMux()

   // Static File server
   staticfileserver := SherryServer.StaticFileServer{documentRoot, "index.html"}
   staticfileserver.AddRouter(router)

   // 啟動 MCP Server
   mcpserver := SryMCPServer.NewMCPServer()
   mcpserver.AddRouter(router)

   // 註冊 TODO 工具
   todoSrv, err := todoMCPServer.NewTodoMCPServer(os.Getenv("ToDoAPIEndpoint"))
   if err == nil {
      todoSrv.AddTools(mcpserver)
   } else {
      fmt.Println("service todo is not initial:", err.Error())
   }

   // 註冊天候
   wSrv, err := weatherMCPServer.NewWeatherMCPServer(os.Getenv("weatherAPIEndpoint"))
   if err == nil {
      wSrv.AddTools(mcpserver)
   } else {
      fmt.Println("service weather is not initial:", err.Error())
   }

   // health check
   m := serverstatus.NewServerStatus(os.Getenv("SystemName"))
   router.HandleFunc("GET /healthz", m.Healthz)

   /*
   // CORS 中間件
   router.Use(func(next http.Handler)(http.Handler) {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
         w.Header().Set("Access-Control-Allow-Origin", "*")
         w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
         w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
         if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
         }
         next.ServeHTTP(w, r)
      })
   })
   */
   return router
}
