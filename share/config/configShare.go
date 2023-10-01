package configShare

import "vgoadmin-go/model/config"

var _allConfigContent configModel.AllConfig

func UpdateMysqlConfig(mysqlConfig configModel.SQLConfig) {
	_allConfigContent.SQLConfig = mysqlConfig
}

func GetMysqlConfig() configModel.SQLConfig {
	return _allConfigContent.SQLConfig
}

func UpdateVgoConfig(vgoConfig configModel.VgoConfig) {
	_allConfigContent.VgoConfig = vgoConfig
}

func GetVgoConfig() configModel.VgoConfig {
	return _allConfigContent.VgoConfig
}

func UpdateRedisConfig(redisConfig configModel.RedisConfig) {
	_allConfigContent.RedisConfig = redisConfig
}

func GetRedisConfig() configModel.RedisConfig {
	return _allConfigContent.RedisConfig
}
