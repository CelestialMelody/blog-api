package models

import (
	"fmt"
	"gin-gorm-practice/conf/setting"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DBList struct {
	MysqlDB *gorm.DB
}

type model struct {
	ID         int64 `json:"id" gorm:"primary_key;column:id;type:bigint(20) unsigned;not null;default:0;comment:'主键'" binding:"required"`
	CreatedOn  int64 `json:"created_on" gorm:"column:created_on;type:varchar(100);not null;default:'';comment:'创建时间'" binding:"required"`
	ModifiedOn int64 `json:"modified_on" gorm:"column:modified_on;type:varchar(100);not null;default:'';comment:'修改时间'" binding:"required"`
}

type BeforeBD interface {
	BeforeCreate(db *gorm.DB) error
	BeforeUpdate(db *gorm.DB) error
}

func initDB() *DBList {
	//dbList := &DBList{}
	dbList := new(DBList)
	db, err := CreateDB(struct {
		Addr           string
		User           string
		Pass           string
		DB             string
		ConnectTimeout uint
	}{
		Addr:           setting.GlobalConfig.GetString("db.mysql.Addr"),
		User:           setting.GlobalConfig.GetString("db.mysql.User"),
		Pass:           setting.GlobalConfig.GetString("db.mysql.Pass"),
		DB:             setting.GlobalConfig.GetString("db.mysql.DB"),
		ConnectTimeout: 10,
	})
	if err != nil {
		//logger.Error("connect DB error: ", err.Error()) // 打印错误日志 hduhelp/server对logger 重新写了, 使用logrus.WithFields
		logrus.Panic("connect DB error: ", err.Error())
	}
	dbList.MysqlDB = db
	return dbList
}

func CreateDSN(dbInfo struct {
	Addr           string
	User           string
	Pass           string
	DB             string
	ConnectTimeout uint
}) string {
	//user:password@/dbname?charset=utf8&parseTime=True&loc=Local
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbInfo.User, dbInfo.Pass, dbInfo.Addr, dbInfo.DB)
}

func CreateDB(dbInfo struct {
	Addr           string
	User           string
	Pass           string
	DB             string
	ConnectTimeout uint
}) (*gorm.DB, error) {
	cfg := struct {
		Addr           string
		User           string
		Pass           string
		DB             string
		ConnectTimeout uint
	}{
		Addr:           dbInfo.Addr,
		User:           dbInfo.User,
		Pass:           dbInfo.Pass,
		DB:             dbInfo.DB,
		ConnectTimeout: dbInfo.ConnectTimeout,
	}
	tablePrefix := setting.GlobalConfig.GetString("db.mysql.TablePrefix")
	DB, err := gorm.Open(mysql.Open(CreateDSN(cfg)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix, // 数据库表前缀
			SingularTable: true,        // 使用单数表名
		},
		PrepareStmt: true,           // 预处理语句
		Logger:      logger.Default, // 日志级别
	})
	return DB, err
}
