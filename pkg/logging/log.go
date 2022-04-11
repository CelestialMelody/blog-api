package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// 找找如何用zap实现

type Level int

var (
	F                  *os.File // 日志文件
	DefaultPrefix      = ""     // 如果不设置前缀，默认为
	DefaultCallerDepth = 2      // 调用层级
	loggo              *log.Logger
	logPrefix          = "" // 日志前缀
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func init() {
	F = openLogFile() // 写入日志文件
	loggo = log.New(F, DefaultPrefix, log.LstdFlags)
}

// 设置日志前缀; 有点难理解
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth) // 获取调用层级
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], file, line) // 格式化前缀
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	loggo.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	loggo.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	loggo.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	loggo.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	loggo.Fatalln(v)
}
