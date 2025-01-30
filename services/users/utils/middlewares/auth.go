package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/parrothacker1/Solvelt/users/config"
)

type jwtclaims struct {
  UserID string `json:"user_id"`
  TeamID string `json:"team_id"`
  Role string `json:"role"`
  jwt.RegisteredClaims
}

var AuthMiddleware http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  authHeader := r.Header.Get("Authorization")
  var ctx context.Context
  if authHeader == "" {
    ctx = context.WithValue(r.Context(),"stop",true)
    r.WithContext(ctx)
    http.Error(w, `{"status":"fail","message":"Missing Authorization header"}`, http.StatusUnauthorized)
		return
  }
	parts := strings.Split(authHeader, " ")
  if len(parts) != 2 || parts[0] != "Bearer" {
    ctx = context.WithValue(r.Context(),"stop",true)
    r.WithContext(ctx)
		http.Error(w, `{"status":"fail","message":"Invalid Authorization format"}`, http.StatusUnauthorized)
		return
  }
	tokenString := parts[1]
  token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
      return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return config.JWTSecret, nil
	})
  if err != nil || !token.Valid {
    ctx = context.WithValue(r.Context(),"stop",true)
    r.WithContext(ctx)
		http.Error(w, `{"status":"fail","message":"Invalid token"}`, http.StatusUnauthorized)
		return
  }
  claims,ok := token.Claims.(jwtclaims)
  if !ok {
    ctx = context.WithValue(r.Context(),"stop",true)
    r.WithContext(ctx)
    http.Error(w,`{"status":"fail","message":"Invalid token claims"}`,http.StatusUnauthorized)
    return
  }
  exp,err := claims.GetExpirationTime()
  if err != nil || exp.Time.Before(time.Now()) {
    ctx = context.WithValue(r.Context(),"stop",true)
    r.WithContext(ctx)
    http.Error(w,`{"status":"fail","message":"Token expired"}`,http.StatusUnauthorized)
    return
  }
  if issuer,err := claims.GetIssuer(); err != nil || issuer != "Solvelt" {
    ctx = context.WithValue(r.Context(),"stop",true)
    r.WithContext(ctx)
    http.Error(w,`{"status":"fail","message":"Invalid token"}`,http.StatusUnauthorized)
    return
  }
  ctx = context.WithValue(r.Context(),"user_id",claims.UserID)
  ctx = context.WithValue(ctx,"team_id",claims.TeamID)
  ctx = context.WithValue(ctx,"role",claims.Role)
  r.WithContext(ctx)
}
