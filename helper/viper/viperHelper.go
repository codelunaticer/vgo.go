package viperHelper

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"vgoadmin-go/model/config"
	"vgoadmin-go/share/config"
)

func _getYamlContent(v *viper.Viper) (*viper.Viper, error) {
	err := v.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("找不到配置文件:%v", err)
		}
		return nil, fmt.Errorf("解析配置文件出错:%v", err)
	}
	return v, nil
}

// 初始化mysql配置数据
func _initSqlConfig() {
	v := viper.New()
	v.SetConfigFile("config/sql.yaml")
	sqlViper, err := _getYamlContent(v)
	if err != nil {
		panic("config/sql.yaml 配置文件读取失败" + err.Error())
	}
	configShare.UpdateMysqlConfig(configModel.SQLConfig{
		Addr:     sqlViper.GetString("SQL.addr"),
		Port:     sqlViper.GetInt("SQL.port"),
		Pass:     sqlViper.GetString("SQL.pass"),
		User:     sqlViper.GetString("SQL.user"),
		Database: sqlViper.GetString("SQL.database"),
		MaxIdle:  sqlViper.GetInt("SQL.maxIdle"),
		MaxOpen:  sqlViper.GetInt("SQL.maxOpen"),
		MaxLife:  sqlViper.GetDuration("SQL.maxLife"),
	})
	configShare.UpdateRedisConfig(configModel.RedisConfig{
		Addr:        sqlViper.GetString("Redis.addr"),
		Password:    sqlViper.GetString("Redis.password"),
		DB:          sqlViper.GetInt("Redis.db"),
		DialTimeout: sqlViper.GetDuration("Redis.dialTimeout"),
	})
}

// 初始化vgo配置数据
func _initVgoConfig() {
	v := viper.New()
	v.SetConfigFile("config/vgo.yaml")
	vgoViper, err := _getYamlContent(v)
	if err != nil {
		panic("config/vgo.yaml 配置文件读取失败" + err.Error())
	}
	configShare.UpdateVgoConfig(configModel.VgoConfig{
		JWTExpiresAt: vgoViper.GetDuration("JWT.ExpiresAt"),
		JWTSecretKey: vgoViper.GetString("JWT.SecretKey"),
		APPPort:      vgoViper.GetString("APP.Port"),
	})
}

func InitConfigData() {
	_initSqlConfig()
	_initVgoConfig()
}
