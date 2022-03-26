package blogAuth

type BlogAuth struct {
	Id       int    `json:"id" gorm:"column:id;type:int(10) unsigned;not null;primary_key;AUTO_INCREMENT"`
	Username string `json:"username" gorm:"column:username;type:varchar(50);not null;default:''"`
	Password string `json:"password" gorm:"column:password;type:varchar(50);not null;default:''"`
}

//CREATE TABLE `blog_auth` (
//`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
//`username` varchar(50) DEFAULT '' COMMENT '账号',
//`password` varchar(50) DEFAULT '' COMMENT '密码',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8;
//
//INSERT INTO `blog`.`blog_auth` (`id`, `username`, `password`) VALUES (null, 'test', 'test123456');
