package utils

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/icosmos-space/iadmin/server/global"
	appReq "github.com/icosmos-space/iadmin/server/model/appclient/request"
)

const appAudience = "IADMIN_APP"

type AppJWT struct {
	SigningKey []byte
}

func NewAppJWT() *AppJWT {
	return &AppJWT{SigningKey: []byte(global.IADMIN_CONFIG.AppJWT.SigningKey)}
}

func (j *AppJWT) CreateClaims(userID uint, username string, uid uuid.UUID) appReq.AppClaims {
	ep, _ := ParseDuration(global.IADMIN_CONFIG.AppJWT.ExpiresTime)
	return appReq.AppClaims{
		UserID:   userID,
		Username: username,
		UUID:     uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{appAudience},
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),
			Issuer:    global.IADMIN_CONFIG.AppJWT.Issuer,
		},
	}
}

func (j *AppJWT) CreateToken(claims appReq.AppClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.SigningKey)
}

func (j *AppJWT) ParseToken(tokenString string) (*appReq.AppClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &appReq.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, TokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}
	if claims, ok := token.Claims.(*appReq.AppClaims); ok && token.Valid {
		okAud := false
		for _, a := range claims.Audience {
			if a == appAudience {
				okAud = true
				break
			}
		}
		if !okAud {
			return nil, TokenInvalid
		}
		return claims, nil
	}
	return nil, TokenValid
}
