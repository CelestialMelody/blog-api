DROP TABLE IF EXISTS `blog_article`;

CREATE TABLE `blog_article`
(
    `id`              int(10) unsigned NOT NULL AUTO_INCREMENT,
    `tag_id`          int(10) unsigned    DEFAULT '0' COMMENT '标签ID',
    `title`           varchar(100)        DEFAULT '' COMMENT '文章标题',
    `desc`            varchar(255)        DEFAULT '' COMMENT '简述',
    `content`         text,
    `cover_image_url` varchar(255)        DEFAULT '' COMMENT '封面图片地址',
    `created_on`      varchar(100)        DEFAULT '' COMMENT '创建时间',
    `created_by`      varchar(100)        DEFAULT '' COMMENT '创建人',
    `modified_on`     varchar(100)        DEFAULT '' COMMENT '修改时间',
    `modified_by`     varchar(100)        DEFAULT '' COMMENT '修改人',
    `deleted_on`      varchar(100)        DEFAULT '' COMMENT '删除时间',
    `state`           tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用1为启用',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='文章管理';

DROP TABLE IF EXISTS `blog_auth`;

CREATE TABLE `blog_auth`
(
    `id`       int(10) unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(50) DEFAULT '' COMMENT '账号',
    `password` varchar(50) DEFAULT '' COMMENT '密码',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

INSERT INTO `blog`.`blog_auth` (`id`, `username`, `password`)
VALUES (1, 'melody', '20011717');

DROP TABLE IF EXISTS `blog_tag`;

CREATE TABLE `blog_tag`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(100)        DEFAULT '' COMMENT '标签名称',
    `created_on`  varchar(100)        DEFAULT '' COMMENT '创建时间',
    `created_by`  varchar(100)        DEFAULT '' COMMENT '创建人',
    `modified_on` varchar(100)        DEFAULT '' COMMENT '修改时间',
    `modified_by` varchar(100)        DEFAULT '' COMMENT '修改人',
    `deleted_on`  varchar(100)        DEFAULT '0' COMMENT '删除时间',
    `state`       tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='文章标签管理';

# // Article 通过AutoMigrate()创建表 mysql.set() 如果要用到外键 貌似不好创建; 放弃
# //type Article struct { // gorm 字段设置; comment 注释 add comment for field when migration
# //	models.Module
# //
# //	// Cannot add foreign key constraint;
# //	TagID uint        `json:"tag_id" gorm:"column:tag_id;type:int(10) unsigned;not null;default:0;comment:'标签ID'" binding:"required"`
# //	Tag   blogTag.Tag `json:"blogTag"`
# //	Title string      `json:"title" gorm:"column:title;type:varchar(100);not null;default:'';comment:'文章标题'" binding:"required"`
# //
# //	Desc       string `json:"desc" gorm:"column:desc;type:varchar(255);not null;default:'';comment:'简述'" binding:"required"`
#     //	Content    string `json:"content" gorm:"column:content;type:text;comment:'内容'" binding:"required"`
# //	CreatedBy  string `json:"created_by" gorm:"column:created_by;type:varchar(100);not null;default:'';comment:'创建人'" binding:"required"`
# //	ModifiedBy string `json:"modified_by" gorm:"column:modified_by;type:varchar(255);not null;default:'';comment:'修改人'" binding:"required"`
# //	// 实际上 这是硬删除 没法用
# //	DeletedOn string `json:"deleted_on" gorm:"column:deleted_on;type:varchar(100);not null;default:'';comment:'删除时间'" binding:"required"`
# //	State     int    `json:"state" gorm:"column:state;type:tinyint(3);not null;default:1;comment:'状态 0为禁用1为启用'" binding:"required"`
# //}


# // BeforeCreate 新建前; gorm v1 *gorm.Scope.SetColumn("CreatedOn", time.Now()); v2 *gorm.mysql.SetColumn("CreatedOn", time.Now())
# func (t *Tag) BeforeCreate(db *gorm.DB) error {
# 	year := time.Now().Year()
# 	month := time.Now().Month()
# 	day := time.Now().Day()
# 	hour := time.Now().Hour()
# 	minute := time.Now().Minute()
# 	second := time.Now().Second()
# 	db.Statement.SetColumn("created_on", fmt.Sprintf("%d-%d-%d %d:%d:%d",
# 		year, month, day, hour, minute, second))
# 	return nil
# }


# func (t *Tag) BeforeUpdate(db *gorm.DB) error {
# 	year := time.Now().Year()
# 	month := time.Now().Month()
# 	day := time.Now().Day()
# 	hour := time.Now().Hour()
# 	minute := time.Now().Minute()
# 	second := time.Now().Second()
# 	db.Statement.SetColumn("modified_on", fmt.Sprintf("%d-%d-%d %d:%d:%d",
# 		year, month, day, hour, minute, second))
# 	return nil
# }


# // BeforeCreate 建议抽象为接口
# func (blogArticle *Article) BeforeCreate(db *gorm.DB) error {
# 	year := time.Now().Year()
# 	month := time.Now().Month()
# 	day := time.Now().Day()
# 	hour := time.Now().Hour()
# 	minute := time.Now().Minute()
# 	second := time.Now().Second()
# 	db.Statement.SetColumn("created_on", fmt.Sprintf("%d-%d-%d %d:%d:%d", year, month, day, hour, minute, second))
# 	return nil
# }


# func (blogArticle *Article) BeforeUpdate(db *gorm.DB) error {
# 	year := time.Now().Year()
# 	month := time.Now().Month()
# 	day := time.Now().Day()
# 	hour := time.Now().Hour()
# 	minute := time.Now().Minute()
# 	second := time.Now().Second()
# 	db.Statement.SetColumn("modified_on", fmt.Sprintf("%d-%d-%d %d:%d:%d", year, month, day, hour, minute, second))
# 	return nil
# }

# //type Auth struct {
# //	Id       int    `json:"id" gorm:"column:id;type:int(10) unsigned;not null;primary_key;AUTO_INCREMENT"`
# //	Username string `json:"username" gorm:"column:username;type:varchar(50);not null;default:''"`
# //	Password string `json:"password" gorm:"column:password;type:varchar(50);not null;default:''"`
# //}
