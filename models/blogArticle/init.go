package blogArticle

type blogArticle struct { // gorm 字段设置; comment 注释 add comment for field when migration
	Id         int64  `json:"id" gorm:"primary_key;column:id;type:bigint(20) unsigned;not null;default:0;comment:'主键'"`
	TagId      int64  `json:"tag_id" gorm:"column:tag_id;type:bigint(20) unsigned;not null;default:0;comment:'标签ID'"`
	Title      string `json:"title" gorm:"column:title;type:varchar(100);not null;default:'';comment:'文章标题'"`
	Desc       string `json:"desc" gorm:"column:desc;type:varchar(255);not null;default:'';comment:'简述'"`
	Content    string `json:"content" gorm:"column:content;type:text;not null;default:'';comment:'内容'"`
	CreatedOn  int64  `json:"created_on" gorm:"column:created_on;type:int(11);not null;default:null;comment:'创建时间'"`
	CreatedBy  string `json:"created_by" gorm:"column:created_by;type:varchar(100);not null;default:'';comment:'创建人'"`
	ModifiedOn int64  `json:"modified_on" gorm:"column:modified_on;type:int(10) unsigned;not null;default:0;comment:'修改时间'"`
	ModifiedBy string `json:"modified_by" gorm:"column:modified_by;type:varchar(255);not null;default:'';comment:'修改人'"`
	DeletedOn  int64  `json:"deleted_on" gorm:"column:deleted_on;type:int(10) unsigned;not null;default:0;comment:'删除时间'"`
	State      int    `json:"state" gorm:"column:state;type:tinyint(3);not null;default:1;comment:'状态 0为禁用1为启用'"`
}

// sql 配置
//CREATE TABLE `blog_article` (
//`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
//`tag_id` int(10) unsigned DEFAULT '0' COMMENT '标签ID',
//`title` varchar(100) DEFAULT '' COMMENT '文章标题',
//`desc` varchar(255) DEFAULT '' COMMENT '简述',
//`content` text,
//`created_on` int(11) DEFAULT NULL,
//`created_by` varchar(100) DEFAULT '' COMMENT '创建人',
//`modified_on` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
//`modified_by` varchar(255) DEFAULT '' COMMENT '修改人',
//`deleted_on` int(10) unsigned DEFAULT '0',
//`state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用1为启用',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章管理';
