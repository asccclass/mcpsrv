package main

import (
   "fmt"
   "time"
   "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("justdrink@gmail.com")

// 用戶數據庫 (簡化版本，生產環境應使用真實數據庫)
var users = map[string]string{
   "admin":    "password123",
   "client1":  "secret456",
   "testuser": "test789",
}

// 生成 JWT Token
func generateToken(username string) (string, time.Time, error) {
   expirationTime := time.Now().Add(24 * time.Hour)
   claims := &Claims{
      Username: username,
      RegisteredClaims: jwt.RegisteredClaims{
         ExpiresAt: jwt.NewNumericDate(expirationTime),
         IssuedAt:  jwt.NewNumericDate(time.Now()),
         Issuer:    "mcp-server",
      },
   }
   token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
   tokenString, err := token.SignedString(jwtSecret)
   if err != nil {
      return "", time.Time{}, err
   }
   return tokenString, expirationTime, nil
}

// 驗證 JWT Token
func validateToken(tokenString string) (*Claims, error) {
   claims := &Claims{}
   token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
      if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
         return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
      }
      return jwtSecret, nil
   })
   if err != nil {
      return nil, err
   }
   if !token.Valid {
      return nil, fmt.Errorf("invalid token")
   }
   return claims, nil
}
