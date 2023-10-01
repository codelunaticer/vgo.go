package menuModel

import dbModel "vgoadmin-go/model/db"

// FindRoleMenusRes 查找某角色拥有的菜单返回值结构体
type FindRoleMenusRes struct {
	dbModel.SQLSysMenuList
	Role struct{} `gorm:"-" json:"role"`
}
