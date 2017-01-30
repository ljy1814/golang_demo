package mylog

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

type Position uint

const (
	POSITION_SINGLE     Position = 1
	POSITION_IN_MANAGER Position = 2
)

func init() {
	log.SetFlags(log.LstdFlags)
}

type Logger interface {
	GetPosition() Position
	SetPosition(pos Position)
	Error(v ...interface{}) string
	Errorf(format string, v ...interface{}) string
	Errorln(v ...interface{}) string
	Fatal(v ...interface{}) string
	Fatalf(format string, v ...interface{}) string
	Fatalln(v ...interface{}) string
	Info(v ...interface{}) string
	Infof(format string, v ...interface{}) string
	Infoln(v ...interface{}) string
	Panic(v ...interface{}) string
	Panicf(format string, v ...interface{}) string
	Panicln(v ...interface{}) string
	Warn(v ...interface{}) string
	Warnf(format string, v ...interface{}) string
	Warnln(v ...interface{}) string
}

func getInvokerLocation(skipNumber int) string {
	//pc程序计数器,file调用的文件名,line调用语句所在行
	//skipNumber跳几个栈帧,0表示当前调用的函数,1表示调用者函数(即上一级函数),
	/*
		test() {	//2
			test1() { //1
				test2() {  //0
					Caller() //此处执行Caller
				}
			}
		}
	*/
	pc, file, line, ok := runtime.Caller(skipNumber)
	if !ok {
		return ""
	}
	simpleFileName := ""
	if index := strings.LastIndex(file, "/"); index > 0 {
		simpleFileName = file[index+1 : len(file)]
	}
	funcPath := ""
	funcPtr := runtime.FuncForPC(pc)
	if funcPtr != nil {
		funcPath = funcPtr.Name()
	}
	return fmt.Sprintf("%s : (%s:%d)", funcPath, simpleFileName, line)
}

func generateLogContent(logTag LogTag, pos Position, format string, v ...interface{}) string {
	skipNumber := int(pos) + 2
	baseInfo := fmt.Sprintf("%s %s %s - ", logTag.Prefix(), time.Now().Format("2006-01-02 15:04:05"), getInvokerLocation(skipNumber))
	var result string
	if len(format) > 0 {
		result = fmt.Sprintf((baseInfo + format), v...)
	} else {
		vLen := len(v)
		params := make([]interface{}, (vLen + 1))
		params[0] = baseInfo
		for i := 1; i <= vLen; i++ {
			params[i] = v[i-1]
		}
		result = fmt.Sprint(params...)
	}
	return result
}

func NewSimpleLogger() Logger {
	logger := &ConsoleLogger{}
	logger.SetPosition(POSITION_SINGLE)
	return logger
}

func NewFileLogger() Logger {
	logger := &FileLogger{}
	logger.SetPosition(POSITION_SINGLE)
	return logger
}
func NewLogger(loggers []Logger) Logger {
	for _, logger := range loggers {
		logger.SetPosition(POSITION_IN_MANAGER)
	}
	return &LogManager{loggers: loggers}
}
