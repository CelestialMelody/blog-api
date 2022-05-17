package blogTag

import (
	"gin-gorm-practice/pkg/mysql"
	"gin-gorm-practice/temp_deprecated"
	"github.com/sirupsen/logrus"
)

type Tag struct {
	mysql.Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func Init() {
	if err := mysql.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Tag{}); err != nil {
		logrus.Panicf("blog_tag migrate failed, %v\n", err)
	}
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag, err error) {
	err = temp.DB.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	return tags, err
}

func GetTagCount(maps interface{}) (count int64, err error) {
	err = temp.DB.Model(&Tag{}).Where(maps).Count(&count).Error
	return count, err
}

func ExistTagByName(name string) error {
	var tag Tag
	err := temp.DB.Select("id").Where("name = ?", name).First(&tag).Error
	return err
}

func ExistTagByID(id int) error {
	var tag Tag
	err := temp.DB.Select("id").Where("id = ?", id).First(&tag).Error
	return err
}

func AddTag(name string, state int, createdBy string) error {
	tag := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}
	err := temp.DB.Model(&Tag{}).Create(&tag).Error
	return err
}

func EditTag(id int, data interface{}) error {
	err := temp.DB.Model(&Tag{}).Where("id = ?", id).Updates(data).Error
	return err
}

func DeleteTag(id int) error {
	err := temp.DB.Where("id = ?", id).Delete(&Tag{}).Error
	return err
}

func ClearAllTag() error {
	err := temp.DB.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{}).Error
	return err
}
