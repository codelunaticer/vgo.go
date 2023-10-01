package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	menuDto "vgoadmin-go/dto/menu"
	httpLogic "vgoadmin-go/logic/http"
	"vgoadmin-go/logic/middleware"
	menuService "vgoadmin-go/service/menu"
	menuListService "vgoadmin-go/service/menuList"
)

// 添加一条菜单
func _addMenu(c *gin.Context) {
	var dto menuDto.AddMenuDto
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := menuService.AddMenu(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 删除一条菜单
func _delMenu(c *gin.Context) {
	var dto menuDto.DelMenuDto
	err := c.ShouldBindUri(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := menuService.DelMenu(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 修改一条菜单
func _putMenu(c *gin.Context) {
	var dto menuDto.PutMenuDto
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := menuService.PutMenu(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 查看一条菜单
func _takeMenu(c *gin.Context) {
	var dto menuDto.TakeMenuDto
	err := c.ShouldBindUri(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	menu, err := menuService.TakeMenu(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, menu, nil)
}

// 查找所有菜单
func _findAllMenu(c *gin.Context) {
	menus, err := menuService.FindMenu(c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, menus, nil)
}

// 为某角色添加某菜单
func _roleAddMenu(c *gin.Context) {
	var dto menuDto.AddRoleMenuDto
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := menuListService.AddRoleMenu(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 为某角色删除某菜单
func _delRoleMenu(c *gin.Context) {
	var dto menuDto.DelRoleMenuDto
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	rowsAffected, err := menuListService.DelRoleMenu(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, rowsAffected, nil)
}

// 查找某角色拥有的菜单
func _findRoleMenus(c *gin.Context) {
	var dto menuDto.FindRoleMenusDto
	err := c.ShouldBindUri(&dto)
	if err != nil {
		httpLogic.BadErrorResponse(c, httpLogic.GetBindErrorTranslate(err))
		return
	}
	menus, err := menuListService.FindRoleMenus(dto, c.Request.Context())
	if err != nil {
		httpLogic.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	httpLogic.OKResponse(c, menus, nil)
}

func UserMenuRouter(router *gin.Engine) {
	menuRouter := router.Group("menu")
	menuRouter.Use(middleware.JWTAuthMiddleware())
	{
		menuRootRouter := menuRouter.Group("root").Use(middleware.RootAccountMiddleware())
		menuRootRouter.POST("add", _addMenu)                         // 添加一条菜单
		menuRootRouter.DELETE("del/:menu_id", _delMenu)              //删除一条菜单
		menuRootRouter.PUT("put", _putMenu)                          //修改一条菜单
		menuRootRouter.GET("take/:menu_id", _takeMenu)               // 查看一条菜单
		menuRootRouter.GET("findAll", _findAllMenu)                  // 查找所有菜单
		menuRootRouter.POST("roleAddMenu", _roleAddMenu)             // 为某角色添加某菜单
		menuRootRouter.POST("delRoleMenu", _delRoleMenu)             // 为某角色删除某菜单
		menuRootRouter.GET("findRoleMenus/:role_id", _findRoleMenus) // 查找某角色拥有的菜单
	}
}
