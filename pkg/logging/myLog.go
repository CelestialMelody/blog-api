package logging

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
)

// 实际未使用 体验一下接口编程

type Log interface {
	getLogFilePath() string
	getLogFileFullPath() string
	openLogFile() (*os.File, error)
	makeDir() error
}

type myLog struct {
	logSavePath string
	logSaveName string
	logFileExt  string
	timeFormat  string
	logger      *zap.Logger
}

//var myLogger myLog

//func init() {
//	myLogger = myLog{
//		logSavePath: "runtime/logs/",
//		logSaveName: "log",
//		logFileExt:  "log", // 日志文件后缀
//		timeFormat:  "20220411",
//		logger:      zap.NewExample(),
//	}
//}

func (myLog *myLog) getLogFilePath() string {
	return myLog.logSavePath
}

func (myLog *myLog) getLogFileFullPath() string {
	preFixPath := myLog.getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", myLog.logSaveName,
		time.Now().Format(myLog.timeFormat), myLog.logFileExt)
	return fmt.Sprintf("%s%s", preFixPath, suffixPath)
}

func (myLog *myLog) makeDir() error {
	dir, _ := os.Getwd() // 返回当前工作目录
	err := os.MkdirAll(dir+myLog.getLogFilePath(), os.ModePerm)
	if err != nil {
		myLog.logger.Error("create dir fail", zap.String("dir", dir+myLog.getLogFilePath()))
		return err
	}
	return nil
}

func (myLog *myLog) openLogFile() (*os.File, error) {
	// 日志文件保存路径
	logSavePath := myLog.getLogFilePath()
	// 判断日志文件保存路径是否存在
	_, err := os.Stat(logSavePath)
	switch { // switch 不接表达式，switch - case 就相当于 if - elseif - else
	case os.IsNotExist(err):
		// 如果不存在则创建
		err := myLog.makeDir()
		if err != nil {
			myLog.logger.Error("创建日志目录失败", zap.String("err", err.Error()))
			return nil, err
		}
	case os.IsPermission(err):
		// 如果没有权限则报错
		myLog.logger.Error("没有权限创建日志目录", zap.String("err", err.Error()))
		return nil, err
	}
	// 打开日志文件
	file, err := os.OpenFile(myLog.getLogFileFullPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // 644 权限 rw-r--r--
	if err != nil {
		myLog.logger.Error("打开日志文件失败", zap.String("err", err.Error()))
		return nil, err
	}
	return file, nil
} //
