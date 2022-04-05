package blogAuth

import (
	"gin-gorm-practice/models"
	"go.uber.org/zap"
)

var (
	instanceDB *models.DBList
	logger     = zap.NewExample()
)

//type Auth struct {
//	Id       int64  `json:"id" gorm:"primary_key"`
//	Username string `json:"username"`
//	Password string `json:"password"`
//}

type Auth struct {
	Id       int    `json:"id" gorm:"column:id;type:int(10) unsigned;not null;primary_key;AUTO_INCREMENT"`
	Username string `json:"username" gorm:"column:username;type:varchar(50);not null;default:''"`
	Password string `json:"password" gorm:"column:password;type:varchar(50);not null;default:''"`
}

func init() {
	instanceDB = models.InitDB()
	if err := instanceDB.MysqlDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Auth{}); err != nil {
		logger.Error("models.InitDB", zap.Error(err))
	}
}

func CheckAuth(username, password string) bool {
	var auth Auth
	instanceDB.MysqlDB.Select("id").Where("username = ? and password = ?",
		username, password).First(&auth)
	if auth.Id > 0 {
		return true
	}
	return false
}

//CREATE TABLE `blog_auth` (
//`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
//`username` varchar(50) DEFAULT '' COMMENT '账号',
//`password` varchar(50) DEFAULT '' COMMENT '密码',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8;
//
//INSERT INTO `blog`.`blog_auth` (`id`, `username`, `password`) VALUES (null, 'test', 'test123456');
