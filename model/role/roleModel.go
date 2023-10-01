package roleModel

import dbModel "vgoadmin-go/model/db"

// FindUserRolesRes 同SQLSysRoleList,但是排除了不必要的User
type FindUserRolesRes struct {
	dbModel.SQLSysRoleList
	User struct{} `gorm:"-" json:"user"`
}
