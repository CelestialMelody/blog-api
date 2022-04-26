package logging

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	//logSavePath = "runtime/logs/"
	logSavePath = "\\runtime\\logs\\"
	logSaveName = "log"
	logFileExt  = "log"      // 日志文件后缀
	timeFormat  = "20060102" // go 诞生日期
	//timeFormat = "20060102150405" // go 诞生日期
	logger = zap.NewExample()
)

func getLogFilePath() string {
	return logSavePath
}

func getLogFileFullPath() string {
	dir, _ := os.Getwd()                 // 获取当前文件所在目录
	preFixPath := dir + getLogFilePath() // 日志文件目录

	zap.L().Debug("logger", zap.String("preFixPath", preFixPath))

	suffixPath := fmt.Sprintf("%s%s.%s", logSaveName,
		time.Now().Format(timeFormat), logFileExt)
	return fmt.Sprintf("%s%s", preFixPath, suffixPath)
}

func makeDir() error {
	dir, _ := os.Getwd() // 返回当前工作目录
	err := os.MkdirAll(dir+getLogFilePath(), os.ModePerm)
	if err != nil {
		logger.Debug("创建日志目录失败", zap.Error(err))
		return err
	}
	logger.Debug("创建日志目录成功", zap.String("dir", dir+getLogFilePath()))
	return nil
}

func openLogFile() *os.File {
	// 日志文件保存路径
	logSavePath := getLogFileFullPath()
	// 判断日志文件保存路径是否存在
	_, err := os.Stat(logSavePath)
	switch { // switch 不接表达式，switch - case 就相当于 if - elseif - else
	case os.IsNotExist(err):
		// 如果不存在则创建
		err := makeDir()
		if err != nil {
			logger.Error("创建日志目录失败", zap.String("err", err.Error()))
			return nil
		}
	case os.IsPermission(err):
		// 如果没有权限则报错
		logger.Error("没有权限创建日志目录", zap.String("err", err.Error()))
		return nil
	}
	// 打开日志文件
	file, err := os.OpenFile(getLogFileFullPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // 644 权限 rw-r--r--
	if err != nil {
		logger.Error("打开日志文件失败", zap.String("err", err.Error()))
		return nil
	}
	return file
}
