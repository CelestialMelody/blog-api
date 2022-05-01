package logging

import (
	"fmt"
	"gin-gorm-practice/conf/setting"
	"gin-gorm-practice/pkg/file"
	"go.uber.org/zap"
	"os"
	"time"
)

var logger = zap.NewExample()

func getLogFilePath() string {
	dir, _ := os.Getwd()
	return fmt.Sprintf("%s/%s%s",
		dir,
		setting.AppSetting.RuntimeRootPath,
		setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

// getLogFileName get the save name of the log file
func getLogFileFullPath() string {
	preFixPath := getLogFilePath()
	suffixPath := getLogFileName()
	return fmt.Sprintf("%s%s", preFixPath, suffixPath)
}

func makeDir() error {
	if err := os.MkdirAll(getLogFilePath(), os.ModePerm); err != nil {
		logger.Debug("创建日志目录失败", zap.Error(err))
		return err
	}
	logger.Debug("创建日志目录成功", zap.String("dir", getLogFilePath()))
	return nil
}

func openLogFile(fileName, filePath string) (*os.File, error) {
	path := getLogFilePath()
	if perm := file.CheckPermission(path); perm == true {
		logger.Debug("日志目录权限不够", zap.String("path", path))
		return nil, fmt.Errorf("file.CheckPermission Permission denied , path: %s", path)
	}
	if err := file.IsNotExistMkDir(path); err != nil {
		logger.Debug("创建日志目录失败", zap.Error(err))
		return nil, fmt.Errorf("file.IsNotExistMkDir failed , path : %s", path)
	}
	f, err := file.Open(path+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Debug("打开日志文件失败", zap.Error(err))
		return nil, fmt.Errorf("fail to OpenFile :%v", err)
	}
	return f, nil
}

//func openLogFile() *os.File {
//	// 日志文件保存路径
//	logSavePath := getLogFileFullPath()
//	// 判断日志文件保存路径是否存在
//	_, err := os.Stat(logSavePath)
//	switch { // switch 不接表达式，switch - case 就相当于 if - elseif - else
//	case os.IsNotExist(err):
//		// 如果不存在则创建
//		err := makeDir()
//		if err != nil {
//			logger.Error("创建日志目录失败", zap.String("err", err.Error()))
//			return nil
//		}
//	case os.IsPermission(err):
//		// 如果没有权限则报错
//		logger.Error("没有权限创建日志目录", zap.String("err", err.Error()))
//		return nil
//	}
//	// 打开日志文件
//	file, err := os.OpenFile(getLogFileFullPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // 644 权限 rw-r--r--
//	if err != nil {
//		logger.Error("打开日志文件失败", zap.String("err", err.Error()))
//		return nil
//	}
//	return file
//}
