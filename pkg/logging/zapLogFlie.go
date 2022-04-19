package logging

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type ZapLevel int

var LoggoZap *zap.Logger

//var DefaultCallerDepthZap = 3

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 获取日志编码器
func getLogWriter(fileName string) zapcore.WriteSyncer {
	dir, _ := os.Getwd() // 获取当前目录
	dir = dir + "/runtime/logs"
	if !pathExists(dir) {
		_ = os.Mkdir(dir, os.ModePerm)
		logs.Warn("create dir %s failed", dir)
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   dir + "/" + fileName, // 日志文件路径
		MaxSize:    1,                    // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	//return zapcore.AddSync(ZapF)
	// 这里说明输出是在文件中
	return zapcore.AddSync(lumberJackLogger)
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := getLogWriter(fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(), writer, level)
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	//if format == "json" {
	//	return zapcore.NewJSONEncoder(getEncoderConfig())
	//}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
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

// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// v0.4.2 这里不需要配置 zap.AddCaller() 就可以获取到文件名和行号; 这里配置会输出zap包的文件名和行号 影响使用
	//_, file, line, _ := runtime.Caller(DefaultCallerDepthZap) // 获取调用层级
	//logPrefix = fmt.Sprintf("[%s:%d]", file, line)            // 格式化前缀
	enc.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
}

func init() {
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
		getEncoderCore(fmt.Sprintf("./%s/server_debug.log", "./zap_log"), debugPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_info.log", "./zap_log"), infoPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_warn.log", "./zap_log"), warnPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_error.log", "./zap_log"), errorPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_panic.log", "./zap_log"), panicPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_fatal.log", "./zap_log"), fatalPriority),
	}

	// zap.AddCaller() 可以获取到文件名和行号
	LoggoZap = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
}
