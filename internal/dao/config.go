package dao

import (
	model "blog-api/internal/models"
	"blog-api/pkg/mysql"
	"github.com/sirupsen/logrus"
)

func Init() {
	err := mysql.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.Article{},
		&model.Auth{},
		&model.Tag{},
	)
	if err != nil {
		logrus.Panicf("blog_article migrate failed, %v\n", err)
	}
}
