package routers

import (
	//_ "gin-gorm-practice/docs" // 不要忘了导入把你上一步生成的docs
	"gin-gorm-practice/middleware/jwt"
	"gin-gorm-practice/pkg/upload"
	"gin-gorm-practice/routers/api"
	v1 "gin-gorm-practice/routers/api/v1"
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
	gin.SetMode(gin.DebugMode) //设置gin的模式，debug模式
	router.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/auth", api.GetAuth) // JWT 验证
	router.POST("/register", api.Register)
	router.POST("/upload", api.UploadImage)
	router.POST("/uploadImages", api.UploadImages)

	apiV1 := router.Group("/api/v1")
	// 接入中间件
	apiV1.Use(jwt.JWT())
	{
		apiV1.GET("/tags", v1.GetTagLists)       // 获取标签列表
		apiV1.POST("/tags", v1.AddTags)          // 新建标签
		apiV1.PUT("/tags/:id", v1.EditTags)      // 更新指定标签
		apiV1.DELETE("/tags/:id", v1.DeleteTags) // 删除指定标签

		apiV1.GET("/articles", v1.GetArticleLists)      // 获取文章列表
		apiV1.GET("/articles/:id", v1.GetArticle)       // 获取指定文章
		apiV1.POST("/articles", v1.AddArticle)          // 新建文章
		apiV1.PUT("/articles/:id", v1.EditArticle)      // 更新指定文章
		apiV1.DELETE("/articles/:id", v1.DeleteArticle) // 删除指定文章
	}

	return router
}
