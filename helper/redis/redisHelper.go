package redisHelper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	configShare "vgoadmin-go/share/config"
)

var _redisClient *redis.Client

func InitRedisClient() {
	redisConfig := configShare.GetRedisConfig()
	_redisClient = redis.NewClient(&redis.Options{
		Addr:        redisConfig.Addr,
		Password:    redisConfig.Password,
		DB:          redisConfig.DB,                        // 数据库编号
		DialTimeout: redisConfig.DialTimeout * time.Second, // 超时时间
	})
	// 测试连接
	_, err := _redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic("redis连接失败," + err.Error())
	}
	fmt.Println("redis连接建立成功!")
}

// RedisSetJson RedisSet 设置键值对
func RedisSetJson(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return _redisClient.Set(context.Background(), key, jsonData, expiration).Err()
}

// RedisGetJson 获取某键值
func RedisGetJson(key string, result interface{}) error {
	jsonData, err := _redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return err
	}
	return nil
}

// RedisDelete 删除指定键
func RedisDelete(key string) error {
	return _redisClient.Del(context.Background(), key).Err()
}

// RedisExist 检测该键是否存在
func RedisExist(key string) (bool, error) {
	exist, err := _redisClient.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	if exist == 1 {
		return true, nil
	}
	return false, nil
}
