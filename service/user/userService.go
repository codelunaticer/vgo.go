package userService

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	appConst "vgoadmin-go/constant/app"
	fileConst "vgoadmin-go/constant/file"
	tipConst "vgoadmin-go/constant/tip"
	userDto "vgoadmin-go/dto/user"
	gormHelper "vgoadmin-go/helper/gorm"
	hashHelper "vgoadmin-go/helper/hash"
	jwtHelper "vgoadmin-go/helper/jwt"
	redisLogic "vgoadmin-go/logic/redis"
	appModel "vgoadmin-go/model/app"
	dbModel "vgoadmin-go/model/db"
	userModel "vgoadmin-go/model/user"
)

// IsExist 账号是否存在
func IsExist(exist userModel.IsExist, ctx context.Context) ([]dbModel.SQLSysUser, error) {
	var whereUser dbModel.SQLSysUser
	var resultUsers []dbModel.SQLSysUser
	if exist.ID != 0 {
		whereUser.ID = exist.ID
	}
	if exist.UserName != "" {
		whereUser.UserName = exist.UserName
	}
	result := gormHelper.NewDBClient(ctx).Where(&whereUser).Find(&resultUsers)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, errors.New(tipConst.AccountNotExist) // 没有记录,账号不存在
	}
	return resultUsers, nil
}

// Login 登录
func Login(loginDto userDto.LoginDto, ctx context.Context) (interface{}, error) {
	users, err := IsExist(userModel.IsExist{UserName: loginDto.Username}, ctx)
	if err != nil {
		return nil, err
	}
	isDel := users[0].IsDel // 账号是否冻结
	if isDel >= 1 {
		return nil, errors.New(tipConst.AccountDel)
	}
	hashPass := users[0].Password
	_, err = hashHelper.CompareHash(loginDto.Password, hashPass)
	if err != nil {
		return nil, errors.New("密码错误!")
	}
	// 获取密码版本
	passVersion, err := redisLogic.GetUserVersion(users[0].ID, true)
	if err != nil {
		return nil, err
	}
	// 生成 token
	token, err2 := jwtHelper.GenerateJWT(appModel.TokenStruct{
		UserId:      users[0].ID,
		PassVersion: passVersion,
	})
	if err2 != nil {
		return nil, err2
	}
	return userModel.LoginRes{
		Token: token,
	}, nil
}

// SetAvatar 修改头像
func SetAvatar(tempFinderFileName string, userId uint, userNick string, oldAvatar string, ctx context.Context) (string, error) {
	accessFilePath := fileConst.AvatarFinder + tempFinderFileName
	err := os.Rename(fmt.Sprintf("%s%s", fileConst.TempFinder, tempFinderFileName), accessFilePath)
	if err != nil {
		return "", err
	}

	accessPath := fileConst.AvatarFinderAccess + "/" + tempFinderFileName
	_, err = SetUserInfo(userDto.SetUserInfoDto{
		NikeName: userNick,
		Avatar:   accessPath,
	}, userId, ctx)
	RemoveAvatar(oldAvatar)
	if err != nil {
		// 数据库修改失败,再将文件移动回去
		os.Rename(accessFilePath, fmt.Sprintf("%s%s", fileConst.TempFinder, tempFinderFileName))
		return "", err
	}
	return accessPath, nil
}

// RemoveAvatar 删除头像
func RemoveAvatar(oldAvatarPath string) {
	fileName := filepath.Base(oldAvatarPath)
	filePath := fileConst.AvatarFinder + fileName
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s文件不存在,不需要删除\n", filePath)
		}
	}
}

// Register 注册
func Register(registerBody userDto.RegisterDto, ctx context.Context) (interface{}, error) {
	hashStr, err := hashHelper.GenerateHash(registerBody.Password)
	if err != nil {
		return "", err
	}
	var user = dbModel.SQLSysUser{
		Password: hashStr,
		UserName: registerBody.Username,
		NikeName: "普通用户",
		IsDel:    0,
	}
	result := gormHelper.NewDBClient(ctx).Create(&user)
	if result.Error != nil && result.RowsAffected < 1 {
		if result.RowsAffected < 1 {
			err = errors.New("受影响行数为0")
		} else {
			err = result.Error
		}
		return "", err
	}
	// 注册成功后需要更新下密码版本
	_, err = redisLogic.UpdateUserVersion(user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// SetUserPass 修改用户密码
func SetUserPass(setUserPassBody userDto.SetPassDto, userId uint, ctx context.Context) (rowsAffected int64, err error) {
	var user dbModel.SQLSysUser
	result := gormHelper.NewDBClient(ctx).Take(&user, userId)
	if result.Error != nil {
		return 0, result.Error
	}
	// 旧密码校验
	_, err = hashHelper.CompareHash(setUserPassBody.OldPass, user.Password)
	if err != nil {
		return 0, errors.New(tipConst.PassError)
	}
	var userResult dbModel.SQLSysUser
	userResult.Password, err = hashHelper.GenerateHash(setUserPassBody.NewPass)
	if err != nil {
		return 0, err
	}
	result2 := gormHelper.NewDBClient(ctx).Where(&dbModel.SQLSysUser{
		ID: userId,
	}).Updates(&userResult)
	if result2.Error != nil {
		return 0, result2.Error
	}
	// 密码更新成功后更新版本号
	_, err = redisLogic.UpdateUserVersion(user.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

// FindAllUser 查看所有账号信息
func FindAllUser(ctx context.Context) ([]dbModel.SQLSysUser, error) {
	var users []dbModel.SQLSysUser
	err := gormHelper.NewDBClient(ctx).Find(&users).Error
	if err != nil {
		return users, err
	}
	return users, nil
}

// DisUser 冻结禁用某个账号
func DisUser(dto userDto.DisUserDto, ctx context.Context) (int64, error) {
	// 账号是否存在
	var user dbModel.SQLSysUser
	err := gormHelper.NewDBClient(ctx).Take(&user, dto.UserId).Error
	if err != nil {
		return 0, err
	}
	// 禁止操作root
	if user.UserName == appConst.RootName {
		return 0, errors.New(tipConst.DisHandle)
	}
	// 已经为冻结状态则不需要处理
	if user.IsDel == 1 {
		return 0, errors.New("该账号已被冻结,无需处理!")
	}
	// 修改冻结字段
	user.IsDel = 1
	err = gormHelper.NewDBClient(ctx).Save(&user).Error
	if err != nil {
		return 0, err
	}
	// 更新密码版本
	_, err = redisLogic.UpdateUserVersion(user.ID)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// OpenUser 解冻某个账号
func OpenUser(dto userDto.OpenUserDto, ctx context.Context) (int64, error) {
	// 账号是否存在
	var user dbModel.SQLSysUser
	err := gormHelper.NewDBClient(ctx).Take(&user, dto.UserId).Error
	if err != nil {
		return 0, err
	}
	// 禁止操作root
	if user.UserName == appConst.RootName {
		return 0, errors.New(tipConst.DisHandle)
	}
	// 已经为解冻状态则不需要处理
	if user.IsDel == 0 {
		return 0, errors.New("该账号未被冻结,无需处理!")
	}
	// 修改冻结字段
	user.IsDel = 0
	err = gormHelper.NewDBClient(ctx).Save(&user).Error
	if err != nil {
		return 0, err
	}
	// 更新密码版本
	_, err = redisLogic.UpdateUserVersion(user.ID)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// GetUserInfo 获取当前用户信息
func GetUserInfo(userId uint, ctx context.Context) (dbModel.SQLSysUser, error) {
	var user dbModel.SQLSysUser
	user.ID = userId
	err := gormHelper.NewDBClient(ctx).Take(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// SetUserInfo 修改当前用户信息
func SetUserInfo(setUserInfoBody userDto.SetUserInfoDto, userId uint, ctx context.Context) (int64, error) {
	var user dbModel.SQLSysUser
	user.ID = userId
	user.NikeName = setUserInfoBody.NikeName
	if setUserInfoBody.NikeName == "" {
		user.NikeName = "普通用户"
	}
	user.Avatar = setUserInfoBody.Avatar
	user.QQ = setUserInfoBody.QQ
	user.Wechat = setUserInfoBody.Wechat
	user.Email = setUserInfoBody.Email
	user.Github = setUserInfoBody.Github
	// 修改
	err := gormHelper.NewDBClient(ctx).Updates(&user).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// ResetUserPass 重置某账号密码
func ResetUserPass(dto userDto.ResetUserPassDto, ctx context.Context) (ResetUserPassRes, error) {
	var result ResetUserPassRes
	result.NewPass = "123456"
	hashPass, err := hashHelper.GenerateHash(result.NewPass)
	if err != nil {
		return result, err
	}
	result2 := gormHelper.NewDBClient(ctx).Where(&dbModel.SQLSysUser{
		ID: dto.UserId,
	}).Updates(&dbModel.SQLSysUser{Password: hashPass})
	if result2.Error != nil {
		return result, result2.Error
	}
	if result2.RowsAffected < 1 {
		return result, errors.New("账号不在,重置失败!")
	}
	// 记得更新密码版本
	_, err = redisLogic.UpdateUserVersion(dto.UserId)
	if err != nil {
		return result, err
	}
	return result, nil
}

type ResetUserPassRes struct {
	NewPass string `json:"new_pass"`
}
