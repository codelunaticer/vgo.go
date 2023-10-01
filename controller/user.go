package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"time"
	fileConst "vgoadmin-go/constant/file"
	userDto "vgoadmin-go/dto/user"
	httpLogic "vgoadmin-go/logic/http"
	"vgoadmin-go/logic/middleware"
	utilsLogic "vgoadmin-go/logic/utils"
	userModel "vgoadmin-go/model/user"
	roleListService "vgoadmin-go/service/roleList"
	userService "vgoadmin-go/service/user"
)

// 注册
func _register(c *gin.Context) {
	var registerBody userDto.RegisterDto
	err := c.ShouldBindJSON(&registerBody)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	newUserInfo, err2 := userService.Register(registerBody, c.Request.Context())
	if err2 != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err2.Error())
		return
	}
	httpLogic.OKResponse(c, newUserInfo, "注册成功")
}

// 登录
func _login(c *gin.Context) {
	var loginBody userDto.LoginDto
	err := c.ShouldBindJSON(&loginBody)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	result, err2 := userService.Login(loginBody, c.Request.Context())
	if err2 != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err2.Error())
		return
	}
	httpLogic.OKResponse(c, result.(userModel.LoginRes), nil)
}

func _setPass(c *gin.Context) {
	var setPassBody userDto.SetPassDto
	err := c.ShouldBindJSON(&setPassBody)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	userId, _ := c.Get("userId")
	rowsAffected, err := userService.SetUserPass(setPassBody, userId.(uint), c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 添加一个超管身份
func _addRootAccount(c *gin.Context) {
	var addRootAccountBody userDto.AddRootAccountDto
	err := c.ShouldBindJSON(&addRootAccountBody)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := roleListService.AddRootAccount(addRootAccountBody, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 撤销一个超管身份
func _delRootAccount(c *gin.Context) {
	var delRootAccountBody userDto.DelRootAccountDto
	err := c.ShouldBindJSON(&delRootAccountBody)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := roleListService.DelRootAccount(delRootAccountBody, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 查看所有账号信息
func _findAllUser(c *gin.Context) {
	users, err := userService.FindAllUser(c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, users, nil)
}

// 冻结禁用某个账号
func _disUser(c *gin.Context) {
	var dto userDto.DisUserDto
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := userService.DisUser(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 解冻某个账号
func _openUser(c *gin.Context) {
	var dto userDto.OpenUserDto
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := userService.OpenUser(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 获取当前用户信息
func _getUserInfo(c *gin.Context) {
	userId, _ := c.Get("userId")
	user, err := userService.GetUserInfo(userId.(uint), c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, user, nil)
}

// 修改当前用户信息
func _serUserInfo(c *gin.Context) {
	userId, _ := c.Get("userId")
	var setUserInfoBody userDto.SetUserInfoDto
	err := c.ShouldBindJSON(&setUserInfoBody)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := userService.SetUserInfo(setUserInfoBody, userId.(uint), c.Request.Context())
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 修改某个用户的密码
func _resetUserPass(c *gin.Context) {
	var dto userDto.ResetUserPassDto
	err := c.ShouldBindUri(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	result, err := userService.ResetUserPass(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, result, nil)
}

// 修改头像
func _setUserAvatar(c *gin.Context) {
	extList := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"} // 格式限制
	maxFileSize := int64(10 << 20)                                // 限制文件大小为 10MB
	file, err := c.FormFile("avatar")
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	ext := filepath.Ext(file.Filename) // 获取文件扩展名
	if !utilsLogic.MContains(extList, ext) {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, "不支持的文件格式!")
		return
	}
	if file.Size > maxFileSize {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, "文件大小不能超出5M!")
		return
	}
	fileNewName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	err = c.SaveUploadedFile(file, fileConst.TempFinder+fileNewName)
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	nickName, _ := c.Get("nikeName")
	userId, _ := c.Get("userId")
	avatar, _ := c.Get("avatar")
	avatarAccessPath, err := userService.SetAvatar(fileNewName, userId.(uint), nickName.(string), avatar.(string), c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, avatarAccessPath, nil)
}

func UseUserRouter(router *gin.Engine) {
	userRouter := router.Group("user")
	{
		userRouter.POST("register", _register) // 注册
		userRouter.POST("login", _login)       // 登录
		userRouteAuth := userRouter.Group("auth").Use(middleware.JWTAuthMiddleware())
		{
			userRouteAuth.PUT("setPass", _setPass)         // 修改密码
			userRouteAuth.GET("getUserInfo", _getUserInfo) // 获取当前用户信息
			userRouteAuth.PUT("setUserInfo", _serUserInfo) // 修改当前用户信息
			userRouteAuth.PUT("setAvatar", _setUserAvatar) // 修改当前用户头像
		}
		userRouteRoot := userRouter.Group("root").Use(middleware.JWTAuthMiddleware(), middleware.RootAccountMiddleware())
		{
			userRouteRoot.GET("findAllUser", _findAllUser) // 查看所有账号信息
			userRouteRoot.POST("disUser", _disUser)        // 冻结禁用某个账号
			userRouteRoot.POST("openUser", _openUser)      // 解冻某个账号
		}
		userRouterSSS := userRouter.Group("SSS").Use(middleware.JWTAuthMiddleware(), middleware.SSSAccountMiddleware())
		{
			userRouterSSS.POST("addRootAccount", _addRootAccount)       // 添加一个超管身份
			userRouterSSS.DELETE("delRootAccount", _delRootAccount)     // 撤销一个超管身份
			userRouterSSS.PUT("resetUserPass/:user_id", _resetUserPass) // 重置某个用户的密码
		}
	}
}
