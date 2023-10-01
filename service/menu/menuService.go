package menuService

import (
	"context"
	"errors"
	menuDto "vgoadmin-go/dto/menu"
	gormHelper "vgoadmin-go/helper/gorm"
	dbModel "vgoadmin-go/model/db"
)

// AddMenu 添加一条菜单
func AddMenu(dto menuDto.AddMenuDto, ctx context.Context) (int64, error) {
	var menu dbModel.SQLSysMenu
	menu.MenuName = dto.MenuName
	err := gormHelper.NewDBClient(ctx).Where(&menu).First(&menu).Error
	if err == nil {
		return 0, errors.New("此菜单已经存在!")
	}
	err = gormHelper.NewDBClient(ctx).Create(&menu).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// DelMenu 删除一条菜单
func DelMenu(dto menuDto.DelMenuDto, ctx context.Context) (int64, error) {
	err := gormHelper.NewDBClient(ctx).Delete(&dbModel.SQLSysMenu{ID: dto.MenuId}).Error
	if err != nil {
		return 0, err
	}
	return 1, err
}

// PutMenu 修改一条菜单
func PutMenu(dto menuDto.PutMenuDto, ctx context.Context) (int64, error) {
	result := gormHelper.NewDBClient(ctx).Where(&dbModel.SQLSysMenu{
		ID: dto.MenuId,
	}).Updates(&dbModel.SQLSysMenu{
		MenuName: dto.MenuName,
	})
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected < 1 {
		return 0, errors.New("不存在此菜单!")
	}
	return 1, nil
}

// TakeMenu 查看一条菜单
func TakeMenu(dto menuDto.TakeMenuDto, ctx context.Context) (dbModel.SQLSysMenu, error) {
	var menu dbModel.SQLSysMenu
	menu.ID = dto.MenuId
	err := gormHelper.NewDBClient(ctx).Take(&menu).Error
	if err != nil {
		return menu, errors.New("查找失败,可能不在!")
	}
	return menu, nil
}

// FindMenu 查找所有菜单
func FindMenu(ctx context.Context) ([]dbModel.SQLSysMenu, error) {
	var menus []dbModel.SQLSysMenu
	err := gormHelper.NewDBClient(ctx).Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}
