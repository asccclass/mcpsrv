package main

import (
   "time"
   "net/http"
   "encoding/json"
   "github.com/golang-jwt/jwt/v5"
)

// 認證相關結構
type AuthRequest struct {
   Username string `json:"username"`
   Password string `json:"password"`
}

type AuthResponse struct {
   Token     string    `json:"token"`
   ExpiresAt time.Time `json:"expires_at"`
}

type Claims struct {
   Username string `json:"username"`
   jwt.RegisteredClaims
}

// 認證端點
func authHandler(w http.ResponseWriter, r *http.Request) {
   var authReq AuthRequest
   if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
      http.Error(w, "Invalid request body", http.StatusBadRequest)
      return
   }
   // 驗證用戶憑證
   if storedPassword, exists := users[authReq.Username]; !exists || storedPassword != authReq.Password {
      http.Error(w, "Invalid credentials", http.StatusUnauthorized)
      return
   }
   // 生成 JWT Token
   token, expiresAt, err := generateToken(authReq.Username)
   if err != nil {
      http.Error(w, "Failed to generate token", http.StatusInternalServerError)
      return
   }
   response := AuthResponse{
      Token:     token,
      ExpiresAt: expiresAt,
   }
   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(response)
}
