package main

import (
	"blog-api/api/v1/controller"
	_ "blog-api/docs"
	"blog-api/pkg/jwt"
	"blog-api/pkg/upload"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func InitRouter() *gin.Engine {
	router := gin.New()      // 初始化路由
	router.Use(gin.Logger()) // 日志
	/*
		Recovery 中间件会恢复(recovers) 任何恐慌(panics)
		如果存在恐慌，中间件将会写入500。
		当你程序里有些异常情况你没考虑到的时候，程序就退出了，服务就停止了，所以是必要的。
	*/
	router.Use(gin.Recovery()) // 异常处理
	router.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiV1 := router.Group("/api/v1")
	// 接入中间件
	apiV1.Use(jwt.JWT())
	{
		apiV1.GET("/authorInfo", controller.GetAuthorInfo)   // JWT 验证
		apiV1.POST("/register", controller.Register)         // 注册
		apiV1.POST("/login", controller.Login)               // 登录
		apiV1.POST("/upload", controller.UploadImage)        // 上传图片
		apiV1.POST("/uploadImages", controller.UploadImages) // 上传多张图片

		apiV1.GET("/tags", controller.GetTagLists)       // 获取标签列表
		apiV1.POST("/tags", controller.AddTags)          // 新建标签
		apiV1.PUT("/tags/:id", controller.EditTags)      // 更新指定标签
		apiV1.DELETE("/tags/:id", controller.DeleteTags) // 删除指定标签

		apiV1.GET("/articles", controller.GetArticleLists)      // 获取文章列表
		apiV1.GET("/articles/:id", controller.GetArticle)       // 获取指定文章
		apiV1.POST("/articles", controller.AddArticle)          // 新建文章
		apiV1.PUT("/articles/:id", controller.EditArticle)      // 更新指定文章
		apiV1.DELETE("/articles/:id", controller.DeleteArticle) // 删除指定文章
	}

	return router
}
