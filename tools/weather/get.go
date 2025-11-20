package weatherMCPServer

import(
   "io"
   "fmt"
   "bytes"
   "net/http"
   "encoding/json"
)

// makeAPIRequest 通用的 API 請求方法
func(s *WeatherMCPServer) makeAPIRequest(method, url string, body interface{}) ([]byte, error) {
   var reqBody io.Reader
   if body != nil {
      jsonData, err := json.Marshal(body)
      if err != nil {
         return nil, fmt.Errorf("marshal request body: %w", err)
      }
      reqBody = bytes.NewBuffer(jsonData)
   }
   req, err := http.NewRequest(method, url, reqBody)
   if err != nil {
      return nil, fmt.Errorf("create request: %w", err)
   }
   if body != nil {
      req.Header.Set("Content-Type", "application/json")
   }
   resp, err := s.httpClient.Do(req)
   if err != nil {
      return nil, fmt.Errorf("make request: %w", err)
   }
   defer resp.Body.Close()
   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("read response: %w", err)
   }
   if resp.StatusCode >= 400 {
      return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
   }
   return respBody, nil
}
