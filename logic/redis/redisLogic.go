package redisLogic

import (
	"errors"
	"strconv"
	redisConst "vgoadmin-go/constant/redis"
	redisHelper "vgoadmin-go/helper/redis"
	appModel "vgoadmin-go/model/app"
)

// UpdateUserVersion 更新用户密码版本
func UpdateUserVersion(id uint) (uint, error) {
	userId := strconv.Itoa(int(id))
	var userVersionData = make(appModel.UserVersionStruct)
	//  是否存在
	exist, err := redisHelper.RedisExist(redisConst.PassVersionKEY)
	if err != nil {
		return 0, err
	}
	if !exist {
		// 不存在需要先创建
		err := redisHelper.RedisSetJson(redisConst.PassVersionKEY, userVersionData, 0)
		if err != nil {
			return 0, err
		}
	}
	// 检查数据
	err = redisHelper.RedisGetJson(redisConst.PassVersionKEY, &userVersionData)
	if err != nil {
		return 0, err
	}
	// 查看是否存在该用户的版本
	currentVersion, ok := userVersionData[userId]
	if ok {
		// 存在则更新
		userVersionData[userId] = currentVersion + 1
	} else {
		// 不存在初始化从1开始
		userVersionData[userId] = 1
	}
	// 将新的版本存储进去
	err = redisHelper.RedisSetJson(redisConst.PassVersionKEY, userVersionData, 0)
	if err != nil {
		return 0, err
	}
	// 返回当前新版本
	return userVersionData[userId], err
}

// GetUserVersion 获取当前密码版本
func GetUserVersion(id uint, isGnerate bool) (uint, error) {
	userId := strconv.Itoa(int(id))
	var userVersionData = make(appModel.UserVersionStruct)
	//  是否存在
	exist, err := redisHelper.RedisExist(redisConst.PassVersionKEY)
	if err != nil {
		return 0, err
	}
	if !exist {
		// 不存在需要先创建
		err := redisHelper.RedisSetJson(redisConst.PassVersionKEY, userVersionData, 0)
		if err != nil {
			return 0, err
		}
	}
	// 检查数据
	err = redisHelper.RedisGetJson(redisConst.PassVersionKEY, &userVersionData)
	if err != nil {
		return 0, err
	}
	// 查看是否存在该用户的版本
	_, ok := userVersionData[userId]
	if !ok {
		if isGnerate { // 鉴权时是不需要初始化的
			// 不存在初始化从1开始
			userVersionData[userId] = 1
			// 将初始的存入
			err = redisHelper.RedisSetJson(redisConst.PassVersionKEY, userVersionData, 0)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, errors.New("没有该账号的版本记录")
		}
	}
	// 返回当前版本
	return userVersionData[userId], err
}
