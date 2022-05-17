package redis

import (
	"encoding/json"
	"gin-gorm-practice/conf"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"time"
)

var RedisConn *redis.Pool

func Init() error {
	RedisConn = &redis.Pool{
		MaxIdle:     conf.RedisSetting.MaxIdle,     // 最大空闲连接数
		MaxActive:   conf.RedisSetting.MaxActive,   // 最大连接数
		IdleTimeout: conf.RedisSetting.IdleTimeout, // 最大空闲连接等待时间
		Dial: func() (redis.Conn, error) { // 连接建立函数
			c, err := redis.Dial("tcp", conf.RedisSetting.Host)
			if err != nil {
				logrus.Errorf("redis连接失败: %v", err)
				return nil, err
			}

			if conf.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", conf.RedisSetting.Password); err != nil {
					logrus.Errorf("redis认证失败: %v", err)
					defer func(c redis.Conn) {
						err := c.Close()
						if err != nil {
							logrus.Errorf("redis关闭失败: %v", err)
						}
					}(c)
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error { // 可选的测试连接函数
			_, err := c.Do("PING") // ping redis
			return err
		},
	}
	logrus.Info("redis连接成功")
	return nil
}

func getConn() redis.Conn {
	conn := RedisConn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			logrus.Error("redis关闭失败", err)
		}
	}(conn)
	return conn
}

func Set(key string, value interface{}, time int) error {
	conn := getConn()

	value, err := json.Marshal(value)
	if err != nil {
		logrus.Errorf("json序列化失败: %v", err)
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		logrus.Errorf("redis设置失败: %v", err)
		return err
	}

	_, err = conn.Do("EXPIRE", key, time) // 设置过期时间
	if err != nil {
		logrus.Errorf("redis设置过期时间失败: %v", err)
		return err
	}

	return nil
}

func Exists(key string) bool {
	conn := getConn()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		logrus.Errorf("redis查询失败: %v", err)
		return false
	}

	return exists
}

func Get(key string) ([]byte, error) {
	conn := getConn()
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		logrus.Errorf("redis查询失败: %v", err)
		return nil, err
	}
	return reply, nil
}

func Delete(key string) error {
	conn := getConn()
	_, err := conn.Do("DEL", key)
	if err != nil {
		logrus.Errorf("redis删除失败: %v", err)
		return err
	}
	return nil
}

func LikeKey(key string) ([]string, error) {
	conn := getConn()
	reply, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		logrus.Errorf("redis查询失败: %v", err)
		return nil, err
	}
	return reply, nil
}

func LikeDel(key string) error {
	conn := getConn()

	keys, err := LikeKey(key)
	if err != nil {
		logrus.Errorf("redis查询失败: %v", err)
		return err
	}

	for _, k := range keys {
		_, err := conn.Do("DEL", k)
		if err != nil {
			logrus.Errorf("redis删除失败: %v", err)
			return err
		}
	}
	return nil
}
