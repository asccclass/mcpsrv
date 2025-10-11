package SryMCPServer

import(
   "os"
   "fmt"
   "strings"
   "net/http"
   "io/ioutil"
)

// 取得Web資料
func(app *MCPServer) GetDataFromWeb(r *http.Request)([]byte, error) {
   var b []byte
   var err error
   if err = r.ParseForm(); err != nil {
      return b, err
   }
   b, err = ioutil.ReadAll(r.Body)
   defer r.Body.Close()
   if err != nil {
      return b, err
   }
   return b, nil
}

// 直接用傳過來的內容複寫檔案 /update/overwrite/{type}
func(app *MCPServer) OverWriteDataFromWeb(w http.ResponseWriter, r *http.Request) {
   fmt.Println("write file to data")
   b, err := app.GetDataFromWeb(r)
   if err != nil {
      fmt.Fprintf(w, "{\"status\": \"failure\", \"message\":\"GetDataFromWeb error:" + err.Error() + "\"}")
      return
   }
   typez := strings.ToLower(r.PathValue("type")) // stock, stocks, found
   if typez == "" {
      fmt.Fprintf(w, "{\"status\": \"failure\", \"message\":\"GetDataFromWeb error:" + err.Error() + "\"}")
      return
   }
   fileName := os.Getenv("DataRoot") + typez + ".json"
   fmt.Println(fileName)
   if err := os.WriteFile(fileName, []byte(b), 0644); err != nil {
      fmt.Fprintf(w, "{\"status\": \"failure\", \"message\":\"" + err.Error() + "\"}")
      return
   }
   fmt.Fprintf(w, "{\"status\": \"ok\", \"message\":\"update jobs finished\"}")
   return
}
