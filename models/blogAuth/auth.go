package blogAuth

import (
	"gin-gorm-practice/pkg/mysql"
	"github.com/sirupsen/logrus"
)

type Auth struct {
	Id          int    `json:"id" gorm:"primary_key"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func Init() {
	if err := mysql.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Auth{}); err != nil {
		logrus.Panicf("blog_auth migrate failed, %v", err)
	}
}

func CheckAuth(username, password string) error {
	var auth Auth
	err := mysql.DB.Select("id").Where("username = ? and password = ?",
		username, password).First(&auth).Error
	return err
}

func Register(username, password string) error {
	var auth Auth
	auth.Username = username
	auth.Password = password
	err := mysql.DB.Create(&auth).Error
	return err
}
