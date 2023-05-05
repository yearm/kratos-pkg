package xredis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yearm/kratos-pkg/config/env"
	"time"
)

// RedisDB ...
type RedisDB *redis.Client

// NewDefaultRedisDB ...
func NewDefaultRedisDB() RedisDB {
	return NewRedisDB(env.GetRedisDefaultConfigPath())
}

// NewRedisDB ...
func NewRedisDB(configPath string) RedisDB {
	return newRedisClient(configPath)
}

// NewRedisClient ...
func NewRedisClient(configPath string) *redis.Client {
	return newRedisClient(configPath)
}

// newRedisClient ...
func newRedisClient(configPath string) *redis.Client {
	addr := viper.GetString(fmt.Sprintf("%s.addr", configPath))
	password := viper.GetString(fmt.Sprintf("%s.password", configPath))
	database := viper.GetInt(fmt.Sprintf("%s.database", configPath))
	maxActive := viper.GetInt(fmt.Sprintf("%s.maxActive", configPath))
	idleTimeout := time.Duration(viper.GetInt(fmt.Sprintf("%s.idleTimeout", configPath))) * time.Second

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password,
		DB:          database,
		MaxRetries:  3,
		IdleTimeout: idleTimeout,
		PoolSize:    maxActive,
	})
	_, err := client.Ping().Result()
	if err != nil {
		logrus.Panicln(fmt.Sprintf("failed to connect redis:[%s], error:%s", addr, err))
	}
	return client
}
