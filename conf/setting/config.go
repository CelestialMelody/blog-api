package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

// GlobalConfig 获取配置 默认全局配置
var GlobalConfig *Config

func Init() { // v0.2.2没有使用 所以报错
	GlobalConfig = &Config{
		Viper: viper.New(),
	}
	// v0.2.2 之前 文件名 搜索地址写错了
	GlobalConfig.SetConfigName("app")    // 读取yaml文件
	GlobalConfig.SetConfigType("yaml")   // 读取yaml文件
	GlobalConfig.AddConfigPath(".")      // 添加搜索路径
	GlobalConfig.AddConfigPath("./conf") // 搜索路径 yaml在当前目录
	err := GlobalConfig.ReadInConfig()
	if err != nil {
		// logrus.WithField  底层为logger.WithField 分配了一个新条目，并在其中添加了一个字段
		// Debug, Print, Info, Warn, Error, Fatal 或 Panic 必须被应用于这个新返回的条目
		// logrus.WithError  将一个错误作为单个字段（使用ErrorKey中定义的键）添加到条目中
		logrus.WithField("config", "GlobalConfig").WithError(err).Panicf("unable to read global config; 读取配置文件失败")
	}

	GlobalConfig.WatchConfig() // 监听配置文件变化 更新配置
	GlobalConfig.OnConfigChange(func(e fsnotify.Event) {
		err := GlobalConfig.ReadInConfig()
		if err != nil {
			// 底层 Entry.log; 签名 func (entry *Entry) Log(level Level, args ...interface{})
			logrus.WithField("config", "GlobalConfig").Info("config file update; change: ", e.Name) // 这里传不传e.Name?
		}
	})
}
