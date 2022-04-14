package logging

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"time"
)

// 找找如何用zap实现

type ZapLevel int

var LoggoZap *zap.Logger
var DefaultCallerDepthZap = 3

func init() {
	dir, _ := os.Getwd()
	path := dir + "/runtime/loggoZaps"
	zap.L().Debug("logging.init zapLogFile", zap.String("path", path))
	if ok := pathExists(path); !ok { // 判断是否有Director文件夹
		zap.L().Debug("创建文件夹zapLogFile", zap.String("path", path)) // 创建文件夹
		_ = os.Mkdir(path, os.ModePerm)
	}
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

	// v0.4 未生效
	cores := [...]zapcore.Core{
		getEncoderCore(fmt.Sprintf("./%s/server_debug.log", "./zap_log"), debugPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_info.log", "./zap_log"), infoPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_warn.log", "./zap_log"), warnPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_error.log", "./zap_log"), errorPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_panic.log", "./zap_log"), panicPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_fatal.log", "./zap_log"), fatalPriority),
	}

	LoggoZap = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	dir, _ := os.Getwd() // 获取当前目录
	dir = dir + "/runtime/logs"
	lumberJackLogger := &lumberjack.Logger{
		Filename:   dir + "/" + fileName, // 日志文件路径
		MaxSize:    1,                    // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	//return zapcore.AddSync(ZapF)
	return zapcore.AddSync(lumberJackLogger)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// getEncoderConfig 获取zapcore.EncoderConfig
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

// getEncoder 获取zapcore.Encoder
func getEncoder(format string) zapcore.Encoder {
	if format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := getLogWriter(fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(""), writer, level)
}

// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	_, file, line, _ := runtime.Caller(DefaultCallerDepthZap) // 获取调用层级
	logPrefix = fmt.Sprintf("[%s:%d]", file, line)            // 格式化前缀
	enc.AppendString(t.Format(logPrefix + "2006/01/02 - 15:04:05.000"))
}
