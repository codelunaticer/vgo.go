package jwtHelper

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	appModel "vgoadmin-go/model/app"
	configShare "vgoadmin-go/share/config"
)

// GenerateJWT 生成JWT
func GenerateJWT(tokenStruct appModel.TokenStruct) (string, error) {
	secretKey := configShare.GetVgoConfig().JWTSecretKey
	hour := configShare.GetVgoConfig().JWTExpiresAt
	claims := appModel.JwtStruct{
		tokenStruct,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(hour * time.Hour)), // 过期时间24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                       // 生效时间
		},
	}
	// 使用HS256签名算法
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString([]byte(secretKey))

	return s, err
}

// ParseJwt 解析JWT
func ParseJwt(tokenString string) (*appModel.JwtStruct, error) {
	secretKey := configShare.GetVgoConfig().JWTSecretKey
	t, err := jwt.ParseWithClaims(tokenString, &appModel.JwtStruct{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if claims, ok := t.Claims.(*appModel.JwtStruct); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
