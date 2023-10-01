package configModel

import "time"

// SQLConfig 数据库配置
type SQLConfig struct {
	Addr     string
	Port     int
	User     string
	Pass     string
	Database string
	MaxIdle  int
	MaxOpen  int
	MaxLife  time.Duration
}

// VgoConfig vgo配置
type VgoConfig struct {
	JWTExpiresAt time.Duration // JWT有效时长
	JWTSecretKey string        // JWT密钥
	APPPort      string        // 程序端口
}

// RedisConfig redis配置
type RedisConfig struct {
	Addr        string        // 连接地址
	Password    string        // 密码
	DB          int           //数据库编号
	DialTimeout time.Duration // 连接超时时间
}

// AllConfig 所有配置
type AllConfig struct {
	SQLConfig   SQLConfig
	VgoConfig   VgoConfig
	RedisConfig RedisConfig
}
