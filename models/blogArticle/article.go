package blogArticle

import (
	"gin-gorm-practice/models/blogTag"
	"gin-gorm-practice/pkg/mysql"
	"github.com/sirupsen/logrus"
)

type Article struct {
	mysql.Model

	TagID         int         `json:"tag_id" gorm:"index" validate:"min=1"`
	Tag           blogTag.Tag `json:"blogTag"` // 关联; 数据表中并无
	Title         string      `json:"title" validate:"min=1,max=100"`
	Desc          string      `json:"desc" validate:"min=1,max=100"`
	Content       string      `json:"content" validate:"min=1"`
	CreatedBy     string      `json:"created_by" validate:"min=1,max=100"`
	ModifiedBy    string      `json:"modified_by" validate:"min=1,max=100"`
	CoverImageUrl string      `json:"cover_image_url" validate:"min=1,max=255"`
	State         int         `json:"state" validate:"oneof=0 1"`
}

func Init() {
	if err := mysql.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Article{}); err != nil {
		logrus.Panicf("blog_article migrate failed, %v\n", err)
	}
}

func ExistArticleByID(id int) error {
	var article Article
	err := mysql.DB.Select("id").Where("id = ?", id).First(&article).Error
	return err
}

func GetArticleTotalCount(maps interface{}) (count int64, err error) {
	err = mysql.DB.Model(&Article{}).Where(maps).Count(&count).Error
	return count, err
}

func GetArticleLists(pageNum int, pageSize int, maps interface{}) (articles []*Article, err error) {
	// mysql.Preload 查询关联表; 大写Tag -> blogArticle.Tag
	err = mysql.DB.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	return articles, err
}

func GetArticle(id int) (*Article, error) {
	var article Article
	err := mysql.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	// mysql.Association 关联查询; 小写tag -> column is 'blogTag'
	err = mysql.DB.Model(&article).Association("blogTag").Find(&article.Tag)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func AddArticle(data map[string]interface{}) error {
	mysql.DB.Create(&Article{
		//map[string]interface{}.(type) 接口类型断言
		TagID:         data["tag_id"].(int),
		Title:         data["title"].(string),
		Desc:          data["desc"].(string),
		Content:       data["content"].(string),
		CreatedBy:     data["created_by"].(string),
		State:         data["state"].(int),
		CoverImageUrl: data["cover_image_url"].(string),
	})
	return nil
}

func EditArticle(id int, data interface{}) error {
	err := mysql.DB.Model(&Article{}).Where("id = ?", id).Updates(data).Error
	return err
}

func DeleteArticle(id int) error {
	err := mysql.DB.Where("id = ?", id).Delete(&Article{}).Error
	return err
}

func CleanAllArticle() error {
	err := mysql.DB.Unscoped().Where("deleted_at != ?", 0).Delete(&Article{}).Error
	return err
}
