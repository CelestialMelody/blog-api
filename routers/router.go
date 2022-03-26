package routers

import (
	v1 "gin-gorm-practice/routers/api/v1"
	"github.com/gin-gonic/gin"
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

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/tags", v1.GetTags)           // 获取标签列表
		apiV1.POST("/tags", v1.AddTags)          // 新建标签
		apiV1.PUT("/tags/:id", v1.EditTags)      // 更新指定标签
		apiV1.DELETE("/tags/:id", v1.DeleteTags) // 删除指定标签
	}

	return router
}

//func InitRouter() *gin.Engine {
//	router := gin.New()
//	router.Use(gin.Logger())
//	router.Use(gin.Recovery())
//	gin.SetMode(gin.DebugMode)
//	router.GET("/test", func(c *gin.Context) {
//		c.JSON(http.StatusOK, gin.H{
//			"message": "test",
//		})
//	})
//	return router
//}
