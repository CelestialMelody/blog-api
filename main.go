package main

import (
	"context"
	"fmt"
	"gin-gorm-practice/conf"
	_ "gin-gorm-practice/docs" // dont forget
	"gin-gorm-practice/models/blogArticle"
	"gin-gorm-practice/models/blogAuth"
	"gin-gorm-practice/models/blogTag"
	"gin-gorm-practice/pkg/log"
	"gin-gorm-practice/pkg/mysql"
	"gin-gorm-practice/pkg/redis"
	"gin-gorm-practice/routers"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/CelestialMelody/gin-gorm-practice
// @license.name MIT
// @license.url https://github.com/CelestialMelody/gin-gorm-practice/blob/main/LICENSE
func main() {
	var err error

	conf.Init()
	err = mysql.Init()
	if err != nil {
		return
	}
	err = redis.Init()
	if err != nil {
		return
	}
	router := routers.InitRouter()

	log.Init()
	blogAuth.Init()
	blogArticle.Init()
	blogTag.Init()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", conf.AppSetting.Port),
		Handler:        router,
		ReadTimeout:    conf.AppSetting.ReadTimeout,
		WriteTimeout:   conf.AppSetting.WriteTimeout,
		MaxHeaderBytes: conf.AppSetting.MaxHeaderBytes,
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
	//blogTag.ClearAllTag()
	//blogArticle.CleanAllArticle()
	fmt.Println("End Exit...")
	os.Exit(0)
}
