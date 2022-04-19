package main

import (
	"context"
	"fmt"
	"gin-gorm-practice/conf/setting"
	"gin-gorm-practice/routers"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title           gin-gorm-practice
// @version         1.0
// @description   	gin-gorm-practice
// @contact.name    API Support
// @license.name    MIT
// @host            localhost:8088
// @BasePath        /api/v1
func main() {
	router := routers.InitRouter()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    time.Duration(setting.ReadTimeout),
		WriteTimeout:   time.Duration(setting.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			zap.L().Info("Listen: %s\n", zap.Any("err", err))
		}
	}()

	quit := make(chan os.Signal)      // 创建一无缓冲通道
	signal.Notify(quit, os.Interrupt) // 当接收到中断信号时，会发送到quit通道中
	<-quit                            // 阻塞直到quit通道被关闭

	zap.L().Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 创建一个5秒的上下文
	defer cancel()                                                          // 当上下文被关闭时，会调用cancel函数

	// 关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown: %s\n", zap.Any("err", err))
	}

	zap.L().Info("Server exiting")
	exit()
}

func exit() {
	fmt.Println("Start Exit...")
	fmt.Println("Exit Clean...")
	fmt.Println("End Exit...")
	os.Exit(0)
}

//func main() { // 用不了 应改linux环境可以
//	// 初始化配置
//	endless.DefaultReadTimeOut = time.Duration(setting.ReadTimeout)
//	endless.DefaultWriteTimeOut = time.Duration(setting.WriteTimeout)
//	endless.DefaultMaxHeaderBytes = 1 << 20 // 1MB
//	endPoint := fmt.Sprintf(":%d", setting.HTTPPort)
//
//	server := endless.NewServer(endPoint, routers.InitRouter())
//	// 启动服务
//	server.BeforeBegin = func(add string) {
//		// 输出当前进程pid
//		zap.L().Info("Actual pid is :", zap.Int("pid", syscall.Getpid()))
//	}
//	err := server.ListenAndServe()
//	if err != nil {
//		zap.L().Panic("Server err", zap.Error(err))
//		return
//	}
//}

//func main() {
//	router := routers.InitRouter()
//	server := &http.Server{
//		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
//		Handler:        router,
//		ReadTimeout:    time.Duration(setting.ReadTimeout),
//		WriteTimeout:   time.Duration(setting.WriteTimeout),
//		MaxHeaderBytes: 1 << 20,
//	}
//	err := server.ListenAndServe()
//	if err != nil {
//		zap.L().Panic("Server err", zap.Error(err))
//		return
//	}
//}

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
