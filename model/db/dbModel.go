package dbModel

import (
	"gorm.io/gorm"
	"time"
)

// SQLSysUser 系统用户表
type SQLSysUser struct {
	ID         uint             `gorm:"primarykey" json:"id"`                               //用户id
	UserName   string           `gorm:"type:VARCHAR(191);not null;unique" json:"user_name"` //用户名
	Password   string           `gorm:"type:VARCHAR(191);not null" json:"-"`                //密码
	NikeName   string           `gorm:"type:VARCHAR(191);default:null" json:"nike_name"`    // 别名
	Avatar     string           `gorm:"type:VARCHAR(191);default:null" json:"avatar"`       // 头像
	QQ         string           `gorm:"type:VARCHAR(191);default:null" json:"qq"`           // qq号
	Wechat     string           `gorm:"type:VARCHAR(191);default:null" json:"wechat"`       // 微信号
	Email      string           `gorm:"type:VARCHAR(191);default:null" json:"email"`        // 邮箱号
	Github     string           `gorm:"type:VARCHAR(191);default:null" json:"github"`       // github号
	IsDel      int              `gorm:"type:TINYINT;not null;default:0" json:"is_del"`      // 是否已删除
	UpdateTime string           `gorm:"type:VARCHAR(191);not null" json:"update_time"`      // 更新时间
	RoleList   []SQLSysRoleList `gorm:"foreignKey:UserID;references:ID" json:"role_list"`
}

func (user *SQLSysUser) TableName() string {
	return "sys_user"
}
func (user *SQLSysUser) BeforeCreate(tx *gorm.DB) (err error) {
	user.UpdateTime = time.Now().String()
	return
}

func (user *SQLSysUser) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdateTime = time.Now().String()
	return
}

// SQLSysRole 系统角色表
type SQLSysRole struct {
	ID       uint   `json:"id" gorm:"primary_key;autoIncrement"`
	RoleName string `json:"role_name" gorm:"not null;unique"`
}

func (SQLSysRole) TableName() string {
	return "sys_role"
}

// SQLSysRoleList 系统角色列表
type SQLSysRoleList struct {
	ID     uint       `gorm:"primary_key;autoIncrement" json:"id"`         // 分配记录id
	RoleID uint       `gorm:"not null" json:"role_id"`                     // 角色id
	UserID uint       `gorm:"not null" json:"user_id"`                     // 用户id
	Role   SQLSysRole `gorm:"foreignKey:RoleID;references:ID" json:"role"` // 关联的父表
	User   SQLSysUser `gorm:"foreignKey:UserID;references:ID" json:"user"` // 关联的父表
}

func (roleList *SQLSysRoleList) TableName() string {

	return "sys_role_list" // 角色分配表
}

// SQLSysMenu 系统菜单表
type SQLSysMenu struct {
	ID       uint             `gorm:"primary_key;autoIncrement" json:"id"`                // 菜单id
	MenuName string           `gorm:"type:VARCHAR(191);not null;unique" json:"menu_name"` // 菜单名称
	MenuList []SQLSysMenuList `gorm:"foreignKey:MenuID;references:ID" json:"menu_list"`
}

func (SQLSysMenu) TableName() string {
	return "sys_menu"
}

// SQLSysMenuList 系统菜单列表
type SQLSysMenuList struct {
	ID     uint       `gorm:"primary_key;autoIncrement" json:"id"`         // id
	MenuID uint       `gorm:"not null" json:"menu_id"`                     // 菜单id
	RoleID uint       `gorm:"not null" json:"role_id"`                     // 角色id
	Role   SQLSysRole `gorm:"foreignKey:RoleID;references:ID" json:"role"` // 关联的父表
	Menu   SQLSysMenu `gorm:"foreignKey:MenuID;references:ID" json:"menu"` // 关联的父表
}

func (SQLSysMenuList) TableName() string {
	return "sys_menu_list"
}
