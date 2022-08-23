package model

import (
	"blog-api/pkg/mysql"
)

type Article struct {
	mysql.Model

	TagID         int    `json:"tag_id" gorm:"index" validate:"min=1"`
	Tag           Tag    `json:"blogTag"` // 关联; 数据表中并无
	Title         string `json:"title" validate:"min=1,max=100"`
	Desc          string `json:"desc" validate:"min=1,max=100"`
	Content       string `json:"content" validate:"min=1"`
	CreatedBy     string `json:"created_by" validate:"min=1,max=100"`
	ModifiedBy    string `json:"modified_by" validate:"min=1,max=100"`
	CoverImageUrl string `json:"cover_image_url" validate:"min=1,max=255"`
	State         int    `json:"state" validate:"oneof=0 1"`
}
