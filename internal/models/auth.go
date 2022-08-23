package model

type Auth struct {
	Id          int    `json:"id" gorm:"primary_key"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
