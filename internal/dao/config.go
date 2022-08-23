package dao

import (
	models2 "gin-gorm-practice/internal/models"
	"gin-gorm-practice/pkg/mysql"
	"github.com/sirupsen/logrus"
)

func Init() {
	err := mysql.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&models2.Article{},
		&models2.Auth{},
		&models2.Tag{},
	)
	if err != nil {
		logrus.Panicf("blog_article migrate failed, %v\n", err)
	}
}
