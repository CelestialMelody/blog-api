package redis

import (
	"blog-api/conf"
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
)

var RDB *redis.Client

func Init() error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     conf.RedisConfig.Host,
		Password: conf.RedisConfig.Password,
		DB:       0, // use default DB
	})
	if _, err := RDB.Ping(context.Background()).Result(); err != nil {
		logrus.Panic("connect redis failed: %v", err)
		return err
	}
	logrus.Info("Connect redis succeeded")
	return nil
}
