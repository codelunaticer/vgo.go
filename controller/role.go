package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	roleDto "vgoadmin-go/dto/role"
	httpLogic "vgoadmin-go/logic/http"
	"vgoadmin-go/logic/middleware"
	menuListService "vgoadmin-go/service/menuList"
	roleService "vgoadmin-go/service/role"
	roleListService "vgoadmin-go/service/roleList"
)

// 删除一个角色
func _delRole(c *gin.Context) {
	var delRoleBody roleDto.DelRoleDto
	err := c.ShouldBindJSON(&delRoleBody)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rows, err := roleService.DelRole(delRoleBody, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rows, nil)
}

// 添加一个角色
func _addRole(c *gin.Context) {
	var addRoleBody roleDto.AddRoleDto
	err := c.ShouldBindJSON(&addRoleBody)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := roleService.AddRole(addRoleBody, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 修改一个角色
func _putRole(c *gin.Context) {
	var putRoleBody roleDto.PutRoleDto
	err := c.ShouldBindJSON(&putRoleBody)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	_, err = roleService.PutRole(putRoleBody, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, 1, nil)
}

// 查找单个角色信息
func _takeRole(c *gin.Context) {
	var takeRoleUri roleDto.TakeRoleDto
	err := c.ShouldBindUri(&takeRoleUri)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	role, err := roleService.TakeRole(takeRoleUri, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, role, nil)
}

// 查找所有的角色信息
func _findAllRole(c *gin.Context) {
	roles, err := roleService.FindAllRoles(c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, roles, nil)
}

// 获取当前用户的角色信息
func _FindCurrentRole(c *gin.Context) {
	userId, _ := c.Get("userId")
	roles, err := roleListService.FindUserRoles(userId.(uint), c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	httpLogic.OKResponse(c, roles, nil)
}

// 获取当前用户的菜单信息
func _FindCurrentMenu(c *gin.Context) {
	userId, _ := c.Get("userId")
	menList, err := menuListService.FindCurrentUserMenu(userId.(uint), c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, menList, nil)
}

// 为某个账户添加某个角色
func _userAddRole(c *gin.Context) {
	var dto roleDto.UserAddRoleDto
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := roleListService.UserAddRole(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 删除某账号的某个角色
func _delUserRole(c *gin.Context) {
	var dto roleDto.DelUserRoleDto
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := roleListService.DelUserRole(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 查找某用户拥有的角色
func _findUserRole(c *gin.Context) {
	var dto roleDto.FindUserRoleDto
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	roles, err := roleListService.FindUserRoles(dto.UserId, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, roles, nil)
}

func UseRoleRouter(router *gin.Engine) {
	roleRouter := router.Group("role")
	roleRouter.Use(middleware.JWTAuthMiddleware())
	{
		roleRouterAuth := roleRouter.Group("auth")
		{
			roleRouterAuth.GET("takeCRole", _FindCurrentRole) // 获取当前用户的角色信息
			roleRouterAuth.GET("takeCMenu", _FindCurrentMenu) // 获取当前用户的菜单信息
		}
		roleRouterRoot := roleRouter.Group("root").Use(middleware.RootAccountMiddleware())
		{
			roleRouterRoot.POST("addRole", _addRole)           // 添加一个角色
			roleRouterRoot.DELETE("delRole", _delRole)         // 删除一个角色
			roleRouterRoot.PUT("putRole", _putRole)            // 修改一个角色
			roleRouterRoot.GET("takeRole/:role_id", _takeRole) // 查找一个角色
			roleRouterRoot.GET("findAllRole", _findAllRole)    // 查找所有角色
			roleRouterRoot.POST("userAddRole", _userAddRole)   // 为某个账户添加某个角色
			roleRouterRoot.POST("delUserRole", _delUserRole)   // 删除某账号的某个角色
			roleRouterRoot.POST("findUserRole", _findUserRole) // 查找某用户拥有的角色
		}
	}
}
