package roleDto

type AddRoleDto struct {
	RoleName string `json:"role_name" binding:"required"` // 角色名称
}

type DelRoleDto struct {
	RoleId uint `json:"role_id" binding:"required"` // 角色id
}

type PutRoleDto struct {
	RoleId   uint   `json:"role_id" binding:"required"`   // 角色id
	RoleName string `json:"role_name" binding:"required"` // 角色名称
}

type TakeRoleDto struct {
	RoleId uint `json:"role_id" uri:"role_id"` // 角色id
}

// UserAddRoleDto 为某账户添加某个角色
type UserAddRoleDto struct {
	UserId uint `json:"user_id" binding:"required"` // 账户id
	RoleId uint `json:"role_id" binding:"required"` // 角色id
}

// DelUserRoleDto 删除某账号的某个角色
type DelUserRoleDto struct {
	UserId uint `json:"user_id" binding:"required"` // 账户id
	RoleId uint `json:"role_id" binding:"required"` // 角色id
}

// FindUserRoleDto 查找某用户拥有的角色
type FindUserRoleDto struct {
	UserId uint `json:"user_id" binding:"required"` // 账户id
}
