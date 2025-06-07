package main

import(
   "fmt"
   "crypto/tls"
)

// 在生產環境中，應該使用有效的 TLS 證書
func generateTLSConfig()(*tls.Config) {
   cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
   if err != nil {
      fmt.Printf("Failed to load TLS certificate, using default config: %s", err.Error())
      return &tls.Config{}
   }

   return &tls.Config{
      Certificates: []tls.Certificate{cert},
   }
}
