package dao

import (
	model "blog-api/internal/models"
	"blog-api/pkg/mysql"
)

func CheckAuth(username, password string) error {
	var auth model.Auth
	err := mysql.DB.Select("id").Where("username = ? and password = ?",
		username, password).First(&auth).Error
	return err
}

func Register(username, password string) error {
	var auth model.Auth
	auth.Username = username
	auth.Password = password
	err := mysql.DB.Create(&auth).Error
	return err
}
