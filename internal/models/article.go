package model

import (
	"blog-api/pkg/mysql"
)

type Article struct {
	mysql.Model

	TagID         int    `json:"tag_id" gorm:"index;not null"`
	Tag           Tag    `json:"blogTag" gorm:"not null"` // 关联; 数据表中并无
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	CoverImageUrl string `json:"cover_image_url"`
}
