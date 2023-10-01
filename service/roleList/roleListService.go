package roleListService

import (
	"context"
	"errors"
	appConst "vgoadmin-go/constant/app"
	tipConst "vgoadmin-go/constant/tip"
	roleDto "vgoadmin-go/dto/role"
	userDto "vgoadmin-go/dto/user"
	gormHelper "vgoadmin-go/helper/gorm"
	redisLogic "vgoadmin-go/logic/redis"
	dbModel "vgoadmin-go/model/db"
	roleModel "vgoadmin-go/model/role"
	userModel "vgoadmin-go/model/user"
	userService "vgoadmin-go/service/user"
)

// DelRootAccount 撤销某个账号的超管角色
func DelRootAccount(delRootAccountBody userDto.DelRootAccountDto, ctx context.Context) (int64, error) {
	// 查找用户信息
	var user dbModel.SQLSysUser
	user.ID = delRootAccountBody.UserId
	err := gormHelper.NewDBClient(ctx).Take(&user).Error
	if err != nil {
		return 0, err
	}
	// 禁止操作root账号
	if user.UserName == appConst.RootName {
		return 0, errors.New(tipConst.DisHandle)
	}
	// 删除这条记录
	err = gormHelper.NewDBClient(ctx).Where(&dbModel.SQLSysRoleList{
		UserID: delRootAccountBody.UserId,
		RoleID: appConst.RootId,
	}).Delete(&dbModel.SQLSysRoleList{}).Error
	if err != nil {
		return 0, err
	}
	// 删除完更新密码版本
	_, err = redisLogic.UpdateUserVersion(delRootAccountBody.UserId)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// AddRootAccount 为某个账户添加超管角色
func AddRootAccount(addRootAccountBody userDto.AddRootAccountDto, ctx context.Context) (int64, error) {
	var roleListItems []dbModel.SQLSysRoleList
	var roleListItem dbModel.SQLSysRoleList
	roleId := appConst.RootId
	// 用户是否存在
	users, err := userService.IsExist(userModel.IsExist{ID: addRootAccountBody.UserId}, ctx)
	if err != nil {
		return 0, err
	}
	// 禁止操作root账号
	if users[0].UserName == appConst.RootName {
		return 0, errors.New(tipConst.DisHandle)
	}
	// 是否已经存在
	result := gormHelper.NewDBClient(ctx).Where(&dbModel.SQLSysRoleList{RoleID: uint(roleId), UserID: addRootAccountBody.UserId}).Find(&roleListItems)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected >= 1 {
		return 0, errors.New("该账号已经是超管角色")
	}
	// 不存在进行添加
	roleListItem.UserID = addRootAccountBody.UserId
	roleListItem.RoleID = appConst.RootId
	err = gormHelper.NewDBClient(ctx).Create(&roleListItem).Error
	if err != nil {
		return 0, err
	}
	// 更新密码版本
	_, err = redisLogic.UpdateUserVersion(addRootAccountBody.UserId)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// FindUserRoles 查找某用户的角色信息
func FindUserRoles(userId uint, ctx context.Context) ([]roleModel.FindUserRolesRes, error) {
	var roles []roleModel.FindUserRolesRes
	err := gormHelper.NewDBClient(ctx).Where(&dbModel.SQLSysRoleList{
		UserID: userId,
	}).Preload("Role").Find(&roles).Error
	if err != nil {
		return roles, err
	}
	return roles, nil
}

// UserAddRole 为某个账户添加某个角色
func UserAddRole(dto roleDto.UserAddRoleDto, ctx context.Context) (int64, error) {
	err := gormHelper.NewDBClient(ctx).Take(&dbModel.SQLSysUser{ID: dto.UserId}).Error
	if err != nil {
		return 0, errors.New("该账号查找失败!")
	}
	err = gormHelper.NewDBClient(ctx).Take(&dbModel.SQLSysRole{ID: dto.RoleId}).Error
	if err != nil {
		return 0, errors.New("该角色查找失败!")
	}
	err = gormHelper.NewDBClient(ctx).Where(&dbModel.SQLSysRoleList{
		UserID: dto.UserId,
		RoleID: dto.RoleId,
	}).Take(&dbModel.SQLSysRoleList{}).Error
	if err == nil {
		return 0, errors.New("该账号已有此角色!")
	}
	err = gormHelper.NewDBClient(ctx).Create(&dbModel.SQLSysRoleList{
		RoleID: dto.RoleId,
		UserID: dto.UserId,
	}).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// DelUserRole 删除某用户的某个角色
func DelUserRole(dto roleDto.DelUserRoleDto, ctx context.Context) (int64, error) {
	err := gormHelper.NewDBClient(ctx).Where(&dbModel.SQLSysRoleList{
		UserID: dto.UserId,
		RoleID: dto.RoleId,
	}).Delete(&dbModel.SQLSysRoleList{}).Error
	if err != nil {
		return 0, errors.New("该用户不存在该角色!")
	}
	return 1, nil
}
