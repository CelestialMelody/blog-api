package blogTag

import (
	"fmt"
	"gin-gorm-practice/models"
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	models.Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// BeforeCreate 新建前; gorm v1 *gorm.Scope.SetColumn("CreatedOn", time.Now()); v2 *gorm.DB.SetColumn("CreatedOn", time.Now())
func (t *Tag) BeforeCreate(db *gorm.DB) error {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	hour := time.Now().Hour()
	minute := time.Now().Minute()
	second := time.Now().Second()
	db.Statement.SetColumn("created_on", fmt.Sprintf("%d-%d-%d %d:%d:%d",
		year, month, day, hour, minute, second))
	//db.Statement.SetColumn("created_on", time.Now().Unix()) // string(time.Now().Unix())
	return nil
}

func (t *Tag) BeforeUpdate(db *gorm.DB) error {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	hour := time.Now().Hour()
	minute := time.Now().Minute()
	second := time.Now().Second()
	db.Statement.SetColumn("modified_on", fmt.Sprintf("%d-%d-%d %d:%d:%d",
		year, month, day, hour, minute, second))
	//db.Statement.SetColumn("modified_on", time.Now().Unix())  // string(time.Now().Unix())
	return nil
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	models.DB.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	return
}

func GetTagTotal(maps interface{}) (count int64) {
	models.DB.Model(&Tag{}).Where(maps).Count(&count)
	return
}

func ExistTagByName(name string) bool {
	var tag Tag
	models.DB.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

func ExistTagByID(id int) bool {
	var tag Tag
	models.DB.Select("id").Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

func AddTag(name string, state int, createdBy string) bool {
	tag := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}
	models.DB.Model(&Tag{}).Create(&tag)
	return true
}

func EditTag(id int, data interface{}) bool {
	models.DB.Model(&Tag{}).Where("id = ?", id).Updates(data)
	return true
}

func DeleteTag(id int) bool {
	models.DB.Where("id = ?", id).Delete(&Tag{})
	return true
}

//type BlogTag struct { // 后面把时间改string了
//	Id         int    `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`
//	Name       string `json:"name" gorm:"column:name;type:varchar(100);DEFAULT:'';COMMENT:'标签名称'"`
//	CreatedOn  int    `json:"created_on" gorm:"column:created_on;type:int(10);DEFAULT:0;COMMENT:'创建时间'"`
//	CreatedBy  string `json:"created_by" gorm:"column:created_by;type:varchar(100);DEFAULT:'';COMMENT:'创建人'"`
//	ModifiedOn int    `json:"modified_on" gorm:"column:modified_on;type:int(10);DEFAULT:0;COMMENT:'修改时间'"`
//	ModifiedBy string `json:"modified_by" gorm:"column:modified_by;type:varchar(100);DEFAULT:'';COMMENT:'修改人'"`
//	DeletedOn  int    `json:"deleted_on" gorm:"column:deleted_on;type:int(10);DEFAULT:0"`
//	State      int    `json:"state" gorm:"column:state;type:tinyint(3);DEFAULT:1;COMMENT:'状态 0为禁用、1为启用'"`
//}

//CREATE TABLE `blog_tag` (
//`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
//`name` varchar(100) DEFAULT '' COMMENT '标签名称',
//`created_on` varchar(100) DEFAULT '' COMMENT '创建时间',
//`created_by` varchar(100) DEFAULT '' COMMENT '创建人',
//`modified_on` varchar(100) DEFAULT '' COMMENT '修改时间',
//`modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
//`deleted_on` varchar(100) DEFAULT '' COMMENT '删除时间',
//`state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章标签管理';
