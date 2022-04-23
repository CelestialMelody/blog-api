package setting

import (
	"gin-gorm-practice/pkg/logging"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"time"
)

// use go-ini

type App struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	Port        int
	TablePrefix string
}

var DatabaseSetting = &Database{}

func SetUp() {
	Cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		logging.LoggoZap.Panic("Fail to parse 'conf/app.ini': %v", zap.Any("err", err))
	}

	if err = Cfg.Section("app").MapTo(AppSetting); err != nil {
		logging.LoggoZap.Panic("Fail to map 'conf/app.ini': %v", zap.Any("err", err))
	}

	if err = Cfg.Section("server").MapTo(ServerSetting); err != nil {
		logging.LoggoZap.Panic("Fail to map 'conf/app.ini': %v", zap.Any("err", err))
	}

	if err = Cfg.Section("database").MapTo(DatabaseSetting); err != nil {
		logging.LoggoZap.Panic("Fail to map 'conf/app.ini': %v", zap.Any("err", err))
	}

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	logging.LoggoZap.Info("init conf success")
}

//var (
//	Cfg          *ini.File
//	RunMode      string
//	HTTPPort     int
//	ReadTimeout  int
//	WriteTimeout int
//	PageSize     int
//	JwtSecret    string
//)

//func init() {
//	var err error
//	Cfg, err = ini.Load("conf/app.ini")
//	if err != nil {
//		panic(err)
//	}
//	LoadBase()
//	LoadServer()
//	LoadApp()
//}
//
//func LoadBase() {
//	RunMode = Cfg.Section("").Key("run_mode").MustString("debug")
//}
//
//func LoadServer() {
//	sec, err := Cfg.GetSection("server")
//	if err != nil {
//		panic(err)
//	}
//	HTTPPort = sec.Key("http_port").MustInt(8088)
//	ReadTimeout = sec.Key("read_timeout").MustInt(60)
//	WriteTimeout = sec.Key("write_timeout").MustInt(60)
//}
//
//func LoadApp() {
//	sec, err := Cfg.GetSection("app")
//	if err != nil {
//		panic(err)
//	}
//	JwtSecret = sec.Key("jwt_secret").MustString("(❁´◡`❁)")
//	PageSize = sec.Key("page_size").MustInt(10)
//}
