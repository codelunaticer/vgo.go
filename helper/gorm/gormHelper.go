package gormHelper

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
	configShare "vgoadmin-go/share/config"
)

var _db *gorm.DB

func InitSqlConnect() {
	sqlConfig := configShare.GetMysqlConfig()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s",
		sqlConfig.User,     // 用户名
		sqlConfig.Pass,     // 密码
		sqlConfig.Addr,     // 主机地址
		sqlConfig.Port,     // 端口号
		sqlConfig.Database, // 数据库名称
	)
	logFile, err := os.Create("sql.log")
	if err != nil {
		panic("无法创建日志文件")
	}
	ormLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Printf("无法建立数据库连接:%v\n", err)
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(sqlConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(sqlConfig.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Second * sqlConfig.MaxLife)
	_db = db

	fmt.Println("数据库连接建立成功!")
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
