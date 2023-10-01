package roleService

import (
	"context"
	"errors"
	appConst "vgoadmin-go/constant/app"
	tipConst "vgoadmin-go/constant/tip"
	roleDto "vgoadmin-go/dto/role"
	gormHelper "vgoadmin-go/helper/gorm"
	dbModel "vgoadmin-go/model/db"
)

// TakeRole 查找某条角色信息
func TakeRole(takeRoleUri roleDto.TakeRoleDto, ctx context.Context) (dbModel.SQLSysRole, error) {
	var role dbModel.SQLSysRole
	role.ID = takeRoleUri.RoleId
	result := gormHelper.NewDBClient(ctx).Take(&role)
	if result.Error != nil {
		return role, result.Error
	}
	return role, nil
}

// FindAllRoles 查找所有角色信息
func FindAllRoles(ctx context.Context) ([]dbModel.SQLSysRole, error) {
	var roles []dbModel.SQLSysRole
	err := gormHelper.NewDBClient(ctx).Find(&roles).Error
	if err != nil {
		return roles, err
	}
	return roles, nil
}

// PutRole 修改一个角色
func PutRole(putRoleBody roleDto.PutRoleDto, ctx context.Context) (int64, error) {
	// 是否为超管 或者 是否存在
	if putRoleBody.RoleId == appConst.RootId {
		return 0, errors.New(tipConst.DisHandle)
	}
	result := gormHelper.NewDBClient(ctx).Take(&dbModel.SQLSysRole{ID: putRoleBody.RoleId})
	if result.Error != nil {
		return 0, result.Error
	}
	// 存在的话进行修改
	err := gormHelper.NewDBClient(ctx).Where(&dbModel.SQLSysRole{ID: putRoleBody.RoleId}).Updates(&dbModel.SQLSysRole{
		RoleName: putRoleBody.RoleName,
	}).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// DelRole  删除一个角色
func DelRole(delRoleBody roleDto.DelRoleDto, ctx context.Context) (int64, error) {
	// 是否为超管 或者 是否存在
	if delRoleBody.RoleId == appConst.RootId {
		return 0, errors.New(tipConst.DisHandle)
	}
	result := gormHelper.NewDBClient(ctx).Take(&dbModel.SQLSysRole{ID: delRoleBody.RoleId})
	if result.Error != nil {
		return 0, result.Error
	}
	// 有的话删除
	err := gormHelper.NewDBClient(ctx).Delete(&dbModel.SQLSysRole{ID: delRoleBody.RoleId}).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// AddRole 添加一个角色
func AddRole(addRoleBody roleDto.AddRoleDto, ctx context.Context) (int64, error) {
	var roles []dbModel.SQLSysRole
	// 是否已经存在
	result := gormHelper.NewDBClient(ctx).Where(&dbModel.SQLSysRole{
		RoleName: addRoleBody.RoleName,
	}).Find(&roles)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected >= 1 {
		return 0, errors.New("已存在该角色!")
	}
	// 不存在则添加
	role := dbModel.SQLSysRole{
		RoleName: addRoleBody.RoleName,
	}
	result2 := gormHelper.NewDBClient(ctx).Create(&role)
	if result2.Error != nil {
		return 0, result2.Error
	}
	return result2.RowsAffected, nil
}
