package model

import (
	"blog-api/pkg/mysql"
)

type Tag struct {
	mysql.Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}
