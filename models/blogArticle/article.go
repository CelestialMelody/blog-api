package blogArticle

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

// test: article+s 通过AutoMigrate()创建表 db.set()
type blogArticles struct { // gorm 字段设置; comment 注释 add comment for field when migration
	TagID int64  `json:"tag_id" gorm:"column:tag_id;type:bigint(20) unsigned;not null;default:0;comment:'标签ID'" binding:"required"`
	Title string `json:"title" gorm:"column:title;type:varchar(100);not null;default:'';comment:'文章标题'" binding:"required"`

	Desc       string `json:"desc" gorm:"column:desc;type:varchar(255);not null;default:'';comment:'简述'" binding:"required"`
	Content    string `json:"content" gorm:"column:content;type:text;not null;default:'';comment:'内容'" binding:"required"`
	CreatedBy  string `json:"created_by" gorm:"column:created_by;type:varchar(100);not null;default:'';comment:'创建人'" binding:"required"`
	ModifiedBy string `json:"modified_by" gorm:"column:modified_by;type:varchar(255);not null;default:'';comment:'修改人'" binding:"required"`
	DeletedOn  int64  `json:"deleted_on" gorm:"column:deleted_on;type:varchar(100);not null;default:'';comment:'删除时间'" binding:"required"`
	State      int    `json:"state" gorm:"column:state;type:tinyint(3);not null;default:1;comment:'状态 0为禁用1为启用'" binding:"required"`
}

// BeforeCreate 建议抽象为接口
func (article *blogArticles) BeforeCreate(db *gorm.DB) error {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	hour := time.Now().Hour()
	minute := time.Now().Minute()
	second := time.Now().Second()
	db.Statement.SetColumn("created_on", fmt.Sprintf("%d-%d-%d %d:%d:%d", year, month, day, hour, minute, second))
	return nil
}

func (article *blogArticles) BeforeUpdate(db *gorm.DB) error {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	hour := time.Now().Hour()
	minute := time.Now().Minute()
	second := time.Now().Second()
	db.Statement.SetColumn("modified_on", fmt.Sprintf("%d-%d-%d %d:%d:%d", year, month, day, hour, minute, second))
	return nil
}
