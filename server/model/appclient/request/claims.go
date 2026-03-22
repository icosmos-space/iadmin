package request

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AppClaims C 端 JWT，与后台 CustomClaims 分离，避免混用 token
type AppClaims struct {
	UserID   uint      `json:"userId"`
	Username string    `json:"username"`
	UUID     uuid.UUID `json:"uuid"`
	jwt.RegisteredClaims
}
