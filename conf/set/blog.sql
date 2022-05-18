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

INSERT INTO `blog`.`blog_auth` (`id`, `username`, `password`)
VALUES (1, 'melody', '20011717');
