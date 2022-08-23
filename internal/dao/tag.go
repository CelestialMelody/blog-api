package dao

import (
	model "blog-api/internal/models"
	"blog-api/pkg/mysql"
)

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []model.Tag, err error) {
	err = mysql.DB.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	return tags, err
}

func GetTagCount(maps interface{}) (count int64, err error) {
	err = mysql.DB.Model(&model.Tag{}).Where(maps).Count(&count).Error
	return count, err
}

func ExistTagByName(name string) error {
	var tag model.Tag
	err := mysql.DB.Select("id").Where("name = ?", name).First(&tag).Error
	return err
}

func ExistTagByID(id int) error {
	var tag model.Tag
	err := mysql.DB.Select("id").Where("id = ?", id).First(&tag).Error
	return err
}

func AddTag(name string, state int, createdBy string) error {
	tag := model.Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}
	err := mysql.DB.Model(&model.Tag{}).Create(&tag).Error
	return err
}

func EditTag(id int, data interface{}) error {
	err := mysql.DB.Model(&model.Tag{}).Where("id = ?", id).Updates(data).Error
	return err
}

func DeleteTag(id int) error {
	err := mysql.DB.Where("id = ?", id).Delete(&model.Tag{}).Error
	return err
}

func ClearAllTag() error {
	err := mysql.DB.Unscoped().Where("deleted_on != ?", 0).Delete(&model.Tag{}).Error
	return err
}
