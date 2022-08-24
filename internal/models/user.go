package model

type User struct {
	ID          int    `json:"id" gorm:"primary_key"`
	Username    string `json:"username" gorm:"not null;unique;index"`
	Password    string `json:"password" gorm:"not null"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
