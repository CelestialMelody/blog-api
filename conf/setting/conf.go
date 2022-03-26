package setting

import "gopkg.in/ini.v1"

// use go-ini

var (
	Cfg          *ini.File
	RunMode      string
	HTTPPort     int
	ReadTimeout  int
	WriteTimeout int
	PageSize     int
	JwtSecret    string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		panic(err)
	}
	LoadBase()
	LoadServer()
	LoadApp()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("run_mode").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		panic(err)
	}
	HTTPPort = sec.Key("http_port").MustInt(8088)
	ReadTimeout = sec.Key("read_timeout").MustInt(60)
	WriteTimeout = sec.Key("write_timeout").MustInt(60)
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		panic(err)
	}
	JwtSecret = sec.Key("jwt_secret").MustString("(❁´◡`❁)")
	PageSize = sec.Key("page_size").MustInt(10)
}

//var Conf = new(AppConfig)
//
//// AppConfig is the configuration of application
//type AppConfig struct {
//	ReleaseMode  bool `ini:"release_mode"`
//	Port         int  `ini:"port"`
//	*MySQLConfig `ini:"mysql"`
//}
//
//// MySQLConfig is the configuration of mysql
//type MySQLConfig struct {
//	Host     string `ini:"host"`
//	Port     int    `ini:"port"`
//	User     string `ini:"user"`
//	Password string `ini:"password"`
//	Database string `ini:"database"`
//}
//
//func Init(file string) error { // 同一个包下只能有一个Init()函数
//	// MapTo maps data sources to given struct.
//	return ini.MapTo(Conf, file)
//}
