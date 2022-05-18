package log

import (
	"fmt"
	"gin-gorm-practice/conf"
	"github.com/astaxie/beego/logs"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type ZapLevel int

var Logger *zap.Logger

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	dir, _ := os.Getwd() // 获取当前目录
	dir = dir + "/runtime/logs"
	if !pathExists(dir) {
		_ = os.Mkdir(dir, os.ModePerm)
		logs.Warn("create dir %s failed", dir)
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   dir + "/" + fileName, // 日志文件路径
		MaxSize:    1,                    // 设置日志文件最大尺寸
		MaxBackups: 3,                    // 设置日志文件最多保存多少个备份
		MaxAge:     28,                   // 设置日志文件最多保存多少天
		Compress:   true,                 // 是否压缩 disabled by default
	}
	// 返回同步方式写入日志文件的zapcore.WriteSyncer
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := getLogWriter(fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(), writer, level)
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(getEncoderConfig())
	//	return zapcore.NewJSONEncoder(getEncoderConfig())
}

func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return config
}

// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(conf.LogSetting.TimeFormat))
}

func Init() {
	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.ErrorLevel
	})
	// panic级别
	panicPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.PanicLevel
	})
	// fatal级别
	fatalPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.FatalLevel
	})

	cores := [...]zapcore.Core{
		getEncoderCore(fmt.Sprintf("./debug.log"), debugPriority),
		getEncoderCore(fmt.Sprintf("./info.log"), infoPriority),
		getEncoderCore(fmt.Sprintf("./warn.log"), warnPriority),
		getEncoderCore(fmt.Sprintf("./error.log"), errorPriority),
		getEncoderCore(fmt.Sprintf("./panic.log"), panicPriority),
		getEncoderCore(fmt.Sprintf("./fatal.log"), fatalPriority),
	}

	// zap.AddCaller() 可以获取到文件名和行号
	Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
}
