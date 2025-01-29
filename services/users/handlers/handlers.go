package handlers

import (
	"github.com/golang-jwt/jwt/v5"
)

type response struct {
  Status string `json:"status"`
  Message string `json:"message"`
}

type jwtclaims struct {
  UserID string `json:"user_id"`
  TeamID string `json:"team_id"`
  Role string `json:"role"`
  jwt.RegisteredClaims
}
