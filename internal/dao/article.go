package dao

import (
	model "blog-api/internal/models"
	"blog-api/pkg/mysql"
)

func ExistArticleByID(id int) error {
	var article model.Article
	err := mysql.DB.Select("id").Where("id = ?", id).First(&article).Error
	return err
}

func GetArticleTotalCount(maps interface{}) (count int64, err error) {
	err = mysql.DB.Model(&model.Article{}).Where(maps).Count(&count).Error
	return count, err
}

func GetArticleLists(pageNum int, pageSize int, maps interface{}) (articles []*model.Article, err error) {
	// mysql.Preload 查询关联表; 大写Tag -> blogArticle.Tag
	err = mysql.DB.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	return articles, err
}

func GetArticle(id int) (*model.Article, error) {
	var article model.Article
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
	mysql.DB.Create(&model.Article{
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
	err := mysql.DB.Model(&model.Article{}).Where("id = ?", id).Updates(data).Error
	return err
}

func DeleteArticle(id int) error {
	err := mysql.DB.Where("id = ?", id).Delete(&model.Article{}).Error
	return err
}

func CleanAllArticle() error {
	err := mysql.DB.Unscoped().Where("deleted_at != ?", 0).Delete(&model.Article{}).Error
	return err
}
