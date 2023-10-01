package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	appConst "vgoadmin-go/constant/app"
	tipConst "vgoadmin-go/constant/tip"
	gormHelper "vgoadmin-go/helper/gorm"
	jwtHelper "vgoadmin-go/helper/jwt"
	httpLogic "vgoadmin-go/logic/http"
	redisLogic "vgoadmin-go/logic/redis"
	dbModel "vgoadmin-go/model/db"
	userModel "vgoadmin-go/model/user"
	userService "vgoadmin-go/service/user"
)

// RootAccountMiddleware 超管账号验证中间件
func RootAccountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := c.Get("userId")
		userName, _ := c.Get("userName")
		if userName != appConst.RootName {
			var roleListItem dbModel.SQLSysRoleList
			result := gormHelper.NewDBClient(c.Request.Context()).Where(&dbModel.SQLSysRoleList{UserID: userId.(uint), RoleID: appConst.RootId}).Take(&roleListItem)
			if result.Error != nil {
				httpLogic.NoAuthResponse(c, http.StatusUnauthorized, tipConst.NoAuth)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// SSSAccountMiddleware SSS账号验证中间件
func SSSAccountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userName, _ := c.Get("userName")
		if userName != appConst.RootName {
			httpLogic.NoAuthResponse(c, http.StatusUnauthorized, tipConst.NoAuth)
			c.Abort()
			return
		}
		c.Next()
	}
}

// JWTAuthMiddleware jwt鉴权中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 是否传入了token
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			httpLogic.NoAuthResponse(c, http.StatusUnauthorized, tipConst.NoAuth)
			c.Abort()
			return
		}
		tokenStrArray := strings.Split(tokenStr, "Bearer ")
		if len(tokenStrArray) < 2 {
			httpLogic.ExeErrorResponse(c, errors.New("token格式不正确").Error())
			c.Abort()
			return
		}
		tokenStr = tokenStrArray[1]
		tokenObj, err := jwtHelper.ParseJwt(tokenStr)
		if err != nil {
			httpLogic.NoAuthResponse(c, http.StatusUnauthorized, tipConst.NoAuth)
			c.Abort()
			return
		}
		// 版本校验
		redisPassVersion, err := redisLogic.GetUserVersion(tokenObj.UserId, false)
		if err != nil || redisPassVersion != tokenObj.TokenStruct.PassVersion { // 无法获取版本信息或者版本不一致
			httpLogic.NoAuthResponse(c, http.StatusUnauthorized, tipConst.NoAuth)
			c.Abort()
			return
		}
		users, err := userService.IsExist(userModel.IsExist{
			ID: tokenObj.UserId,
		}, c.Request.Context())
		if err != nil {
			httpLogic.NoAuthResponse(c, http.StatusUnauthorized, tipConst.NoAuth)
			c.Abort()
			return
		}
		if users[0].IsDel >= 1 {
			httpLogic.NoAuthResponse(c, http.StatusUnauthorized, tipConst.AccountDel)
			c.Abort()
			return
		}
		c.Set("userId", tokenObj.TokenStruct.UserId)
		c.Set("userName", users[0].UserName)
		c.Set("nikeName", users[0].NikeName)
		c.Set("avatar", users[0].Avatar)
		c.Set("passVersion", tokenObj.TokenStruct.PassVersion)
		c.Next()
	}
}

// GlobalErrorMiddleware 全局错误捕获
func GlobalErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("程序错误:", r)
				if err, ok := r.(error); ok {
					httpLogic.ExeErrorResponse(c, err.Error())
				} else {
					httpLogic.ExeErrorResponse(c, "")
				}
			}
		}()
		c.Next()
	}
}
