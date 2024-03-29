package log

import (
	"blog-api/conf"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type ZapLevel int

var Logger *zap.Logger

func Init() {
	if conf.AppConfig.RunMode == "debug" {
		// 开发模式 日志输出到终端
		core := zapcore.NewTee(
			zapcore.NewCore(getEncoder(),
				zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
		Logger = zap.New(core, zap.AddCaller())
	} else {
		fileLog()
	}
}

func fileLog() {
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

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	dir, _ := os.Getwd() // 获取当前目录
	dir = dir + "/logs"
	if !pathExists(dir) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			logs.Warn("create dir %s failed\n\treson is %s", dir, err.Error())
		}
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
	writer := getLogWriter(fileName)
	return zapcore.NewCore(getEncoder(), writer, level)
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(getEncoderConfig())
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
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 将日志级别字符串转化为小写
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行消耗时间转化成浮点型的秒
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 以包/文件:行号 格式化调用堆栈
	}
	return config
}

// CustomTimeEncoder 自定义日志输出时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(conf.LogConfig.TimeFormat))
}

//func Error(msg string, err error) {
//	Logger.Error(msg, zap.Error(err))
//}
//func Info(msg string) {
//	Logger.Info(msg)
//}
