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

type Module struct {
	ID         int    `json:"id" gorm:"primary_key;column:id;type:int(10) unsigned;not null;default:0;comment:'主键'" binding:"required"`
	CreatedOn  string `json:"created_on" gorm:"column:created_on;type:varchar(100);not null;default:'';comment:'创建时间'" binding:"required"` // v0.2.2 前写错了类型
	ModifiedOn string `json:"modified_on" gorm:"column:modified_on;type:varchar(100);not null;default:'';comment:'修改时间'" binding:"required"`
}

type BeforeBD interface {
	BeforeCreate(db *gorm.DB) error
	BeforeUpdate(db *gorm.DB) error
}

func InitDB() *DBList {
	//dbList := &DBList{} // v2.2 忘记config Init
	setting.Init()
	dbList := new(DBList)
	db, err := CreateDB(struct {
		Addr           string
		User           string
		Pass           string
		DB             string
		ConnectTimeout uint
	}{
		// 查看配置文件 ;v2.2 配置写错了
		Addr:           setting.GlobalConfig.GetString("db.mysql.Address"),
		User:           setting.GlobalConfig.GetString("db.mysql.User"),
		Pass:           setting.GlobalConfig.GetString("db.mysql.Password"),
		DB:             setting.GlobalConfig.GetString("db.mysql.Database"),
		ConnectTimeout: 10,
	})
	if err != nil {
		//logger.Error("connect DB error: ", err.Error()) // 打印错误日志 hduhelp/server对logger 重新写了, 使用logrus.WithFields
		logrus.Panic("connect DB error: ", err.Error())
	}
	dbList.MysqlDB = db

	// auto migrate ; 不应该放这里 这里的module 我在article中有用到 循环依赖

	logrus.Infof("Connected DB success")

	database, err := dbList.MysqlDB.DB()
	if err != nil {
		logrus.Println("get DB error: ", err.Error())
	}
	// 实际上我不太了解
	database.SetMaxIdleConns(10)     // 设置最大空闲连接数
	database.SetMaxOpenConns(100)    // 设置最大连接数
	database.SetConnMaxLifetime(100) // 设置连接最大存活时间

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
