package model

import "blog-api/pkg/mysql"

type User struct {
	mysql.Model

	//ID          int    `json:"id" gorm:"primaryKey"`
	Username    string `json:"username" gorm:"not null;unique;index"`
	Password    string `json:"password" gorm:"not null"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
