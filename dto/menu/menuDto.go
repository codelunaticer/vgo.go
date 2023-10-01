package menuDto

// AddMenuDto 添加一条菜单
type AddMenuDto struct {
	MenuName string `json:"menu_name" binding:"required"`
}

// DelMenuDto 删除一条菜单
type DelMenuDto struct {
	MenuId uint `json:"menu_id" uri:"menu_id" binding:"required"`
}

// PutMenuDto 修改一条菜单
type PutMenuDto struct {
	MenuId   uint   `json:"menu_id" binding:"required"`   // 菜单id
	MenuName string `json:"menu_name" binding:"required"` // 菜单名称
}

// TakeMenuDto 查看一条菜单
type TakeMenuDto struct {
	MenuId uint `json:"menu_id" uri:"menu_id" binding:"required"`
}

// AddRoleMenuDto 为某角色添加某菜单
type AddRoleMenuDto struct {
	RoleId uint `json:"role_id" binding:"required"` // 角色id
	MenuId uint `json:"menu_id" binding:"required"` // 菜单id
}

// DelRoleMenuDto 删除某角色的某菜单
type DelRoleMenuDto struct {
	RoleId uint `json:"role_id" binding:"required"` // 角色id
	MenuId uint `json:"menu_id" binding:"required"` // 菜单id
}

// FindRoleMenusDto 查找某角色拥有的菜单
type FindRoleMenusDto struct {
	RoleId uint `json:"role_id" uri:"role_id" binding:"required" // 角色id`
}
