package appModel

import "github.com/golang-jwt/jwt/v5"

type TokenStruct struct {
	UserId      uint `json:"userId"`      // 用户id
	PassVersion uint `json:"passVersion"` // 密码版本
}

type JwtStruct struct {
	TokenStruct
	jwt.RegisteredClaims // v5版本新加的方法
}

type UserVersionStruct = map[string]uint // redis中用户版本记录存放的数据结构
