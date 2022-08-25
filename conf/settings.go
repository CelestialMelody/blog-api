package conf

import (
	"github.com/sirupsen/logrus"
	"time"
)

type App struct {
	Port            string
	Host            string
	JwtSecret       string
	RuntimeRootPath string
	PageSize        int
	MaxHeaderBytes  int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	Release         bool
	RunMode         string
}

type MySqlDB struct {
	Address     string
	User        string
	Password    string
	Database    string
	MaxIdle     int
	MaxOpen     int
	MaxLifetime time.Duration
	TablePrefix string
}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type Image struct {
	ImagePrefixUrl string
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExt  []string
}

type Log struct {
	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var (
	AppConfig   App
	DBConfig    MySqlDB
	RedisConfig Redis
	ImageConfig Image
	LogConfig   Log
)

func Init() {
	// 初始化配置
	var err error

	if err = GlobalConfig.UnmarshalKey("app", &AppConfig); err != nil {
		logrus.Panicf("解析配置文件失败 app: %v", err)
	}

	if err = GlobalConfig.UnmarshalKey("db.mysql", &DBConfig); err != nil {
		logrus.Panicf("解析配置文件失败 database: %v", err)
	}

	if err = GlobalConfig.UnmarshalKey("db.redis", &RedisConfig); err != nil {
		logrus.Panicf("解析配置文件失败 redis: %v", err)
	}

	if err = GlobalConfig.UnmarshalKey("image", &ImageConfig); err != nil {
		logrus.Panicf("解析配置文件失败 image: %v", err)
	}

	if err = GlobalConfig.UnmarshalKey("log", &LogConfig); err != nil {
		logrus.Panicf("解析配置文件失败 log: %v", err)
	}

	AppConfig.ReadTimeout = AppConfig.ReadTimeout * time.Second
	AppConfig.WriteTimeout = AppConfig.WriteTimeout * time.Second
	AppConfig.MaxHeaderBytes = AppConfig.MaxHeaderBytes << 20
	ImageConfig.ImageMaxSize = ImageConfig.ImageMaxSize * 1024 * 1024

	//fmt.Printf("%+v\n", DBConfig)
}
