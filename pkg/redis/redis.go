package redis

import (
	"encoding/json"
	"gin-gorm-practice/conf/setting"
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisConn *redis.Pool

func Init() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,     // 最大空闲连接数
		MaxActive:   setting.RedisSetting.MaxActive,   // 最大连接数
		IdleTimeout: setting.RedisSetting.IdleTimeout, // 最大空闲连接等待时间
		Dial: func() (redis.Conn, error) { // 连接建立函数
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				setting.Logger.Error("redis连接失败", err)
				return nil, err
			}

			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					setting.Logger.Error("redis认证失败", err)
					// 关闭连接
					//c.Close()
					defer func(c redis.Conn) {
						err := c.Close()
						if err != nil {
							setting.Logger.Error("redis关闭失败", err)
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
	setting.Logger.Debug("redis连接成功")
	return nil
}

func getConn() redis.Conn {
	conn := RedisConn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			setting.Logger.Error("redis关闭失败", err)
		}
	}(conn)
	return conn
}

func Set(key string, value interface{}, time int) error {
	conn := getConn()

	value, err := json.Marshal(value)
	if err != nil {
		setting.Logger.Error("json序列化失败", err)
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		setting.Logger.Error("redis设置失败", err)
		return err
	}

	_, err = conn.Do("EXPIRE", key, time) // 设置过期时间
	if err != nil {
		setting.Logger.Error("redis设置过期时间失败", err)
		return err
	}

	return nil
}

func Exists(key string) bool {
	conn := getConn()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		setting.Logger.Error("redis查询失败", err)
		return false
	}

	return exists
}

func Get(key string) ([]byte, error) {
	conn := getConn()
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		setting.Logger.Error("redis查询失败", err)
		return nil, err
	}
	return reply, nil
}

func Delete(key string) error {
	conn := getConn()
	_, err := conn.Do("DEL", key)
	if err != nil {
		setting.Logger.Error("redis删除失败", err)
		return err
	}
	return nil
}

func LikeKey(key string) ([]string, error) {
	conn := getConn()
	reply, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		setting.Logger.Error("redis查询失败", err)
		return nil, err
	}
	return reply, nil
}

func LikeDel(key string) error {
	conn := getConn()

	keys, err := LikeKey(key)
	if err != nil {
		setting.Logger.Error("redis查询失败", err)
		return err
	}

	for _, k := range keys {
		_, err := conn.Do("DEL", k)
		if err != nil {
			setting.Logger.Error("redis删除失败", err)
			return err
		}
	}
	return nil
}
