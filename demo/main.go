package main

import (
	"fmt"
	"gin-gorm-practice/conf/setting"
	"gin-gorm-practice/routers"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	router := routers.InitRouter()
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    time.Duration(setting.ReadTimeout),
		WriteTimeout:   time.Duration(setting.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		logrus.Panic("ListenAndServe: ", err)
		return
	}
}

//func main() {
//	router := gin.Default()
//	router.GET("/test", func(c *gin.Context) {
//		c.JSON(200, gin.H{
//			"message": "test",
//		})
//	})
//	// 启动服务 这里居然知道处于debug模式
//	server := &http.Server{
//		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
//		Handler:        router,
//		ReadTimeout:    time.Duration(setting.ReadTimeout),
//		WriteTimeout:   time.Duration(setting.WriteTimeout),
//		MaxHeaderBytes: 1 << 20,
//	}
//	err := server.ListenAndServe()
//	if err != nil {
//		logrus.Printf("start server error: %v", err)
//		return
//	}
//	// gin 直接写不知道处于debug模式
//	//err := router.Run(fmt.Sprintf(":%d", setting.HTTPPort))
//	//if err != nil {
//	//	return
//	//}
//}
