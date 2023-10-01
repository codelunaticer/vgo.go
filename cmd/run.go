package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
	fileConst "vgoadmin-go/constant/file"
	"vgoadmin-go/controller"
	gormHelper "vgoadmin-go/helper/gorm"
	redisHelper "vgoadmin-go/helper/redis"
	validatorHelper "vgoadmin-go/helper/validator"
	"vgoadmin-go/helper/viper"
	"vgoadmin-go/logic/middleware"
	mosLogic "vgoadmin-go/logic/mos"
	configShare "vgoadmin-go/share/config"
)

func main() {
	// 初始化翻译环境
	validatorHelper.InitTrans("zh")
	// 初始化项目配置文件数据
	viperHelper.InitConfigData()
	// 初始化静态目录
	mosLogic.InitAppStaticFinder()
	// 初始数据库
	gormHelper.InitSqlConnect()
	// 初始redis
	redisHelper.InitRedisClient()
	// 日志
	logfile, _ := os.Create("app.log")
	gin.DefaultWriter = io.MultiWriter(logfile, os.Stdout)
	// 路由
	router := gin.Default()
	router.Use(middleware.GlobalErrorMiddleware())
	router.Static(fileConst.AvatarFinderAccess, fileConst.AvatarFinder)
	controller.UseAppRouter(router)
	controller.UseUserRouter(router)
	controller.UseRoleRouter(router)
	controller.UserMenuRouter(router)
	// 服务
	err := router.Run(":" + configShare.GetVgoConfig().APPPort)
	if err != nil {
		return
	}
}
