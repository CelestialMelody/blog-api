package blogArticle

import (
	"fmt"
	"gin-gorm-practice/models"
	"gin-gorm-practice/models/blogTag"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

var instanceDB *models.DBList

func init() {
	instanceDB = models.InitDB()
}

// BlogArticles 通过AutoMigrate()创建表 db.set()
type BlogArticles struct { // gorm 字段设置; comment 注释 add comment for field when migration
	models.Module

	TagID int64       `json:"tag_id" gorm:"column:tag_id;type:bigint(20) unsigned;not null;default:0;comment:'标签ID'" binding:"required"`
	Tag   blogTag.Tag `json:"tag"`
	Title string      `json:"title" gorm:"column:title;type:varchar(100);not null;default:'';comment:'文章标题'" binding:"required"`

	Desc       string `json:"desc" gorm:"column:desc;type:varchar(255);not null;default:'';comment:'简述'" binding:"required"`
	Content    string `json:"content" gorm:"column:content;type:text;not null;default:'';comment:'内容'" binding:"required"`
	CreatedBy  string `json:"created_by" gorm:"column:created_by;type:varchar(100);not null;default:'';comment:'创建人'" binding:"required"`
	ModifiedBy string `json:"modified_by" gorm:"column:modified_by;type:varchar(255);not null;default:'';comment:'修改人'" binding:"required"`
	// 实际上 这是硬删除 没法用
	DeletedOn string `json:"deleted_on" gorm:"column:deleted_on;type:varchar(100);not null;default:'';comment:'删除时间'" binding:"required"`
	State     int    `json:"state" gorm:"column:state;type:tinyint(3);not null;default:1;comment:'状态 0为禁用1为启用'" binding:"required"`
}

// BeforeCreate 建议抽象为接口
func (article *BlogArticles) BeforeCreate(db *gorm.DB) error {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	hour := time.Now().Hour()
	minute := time.Now().Minute()
	second := time.Now().Second()
	db.Statement.SetColumn("created_on", fmt.Sprintf("%d-%d-%d %d:%d:%d", year, month, day, hour, minute, second))
	return nil
}

func (article *BlogArticles) BeforeUpdate(db *gorm.DB) error { // 大写
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	hour := time.Now().Hour()
	minute := time.Now().Minute()
	second := time.Now().Second()
	db.Statement.SetColumn("modified_on", fmt.Sprintf("%d-%d-%d %d:%d:%d", year, month, day, hour, minute, second))
	return nil
}

// ExistArticleByID 根据ID查询文章是否存在
func ExistArticleByID(id int64) (bool, error) {
	var article BlogArticles
	err := instanceDB.MysqlDB.Select("id").Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if article.TagID > 0 {
		return true, nil
	}
	return false, nil
}

// GetArticleTotal 查询文章总数
func GetArticleTotal(maps interface{}) (count int64) {
	instanceDB.MysqlDB.Model(&BlogArticles{}).Where(maps).Count(&count)
	return
}

// GetArticles 获取文章列表
func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []BlogArticles) {
	instanceDB.MysqlDB.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
	return
}

// GetArticle 获取文章
func GetArticle(id int64) (article BlogArticles) {
	instanceDB.MysqlDB.Where("id = ?", id).First(&article)
	err := instanceDB.MysqlDB.Model(&article).Association("tag").Find(&article.Tag)
	if err != nil {
		logrus.Debugln("Can't Find Article", err)
		return BlogArticles{}
	}
	return
}

// EditArticle 编辑文章; 只会返回true 不太好
func EditArticle(id int64, data interface{}) bool {
	instanceDB.MysqlDB.Model(&BlogArticles{}).Where("id = ?", id).Updates(data)
	return true
}

// DeleteArticle 删除文章 只会返回true 不太好
func DeleteArticle(id int64) bool {
	instanceDB.MysqlDB.Where("id = ?", id).Delete(&BlogArticles{})
	return true
}
