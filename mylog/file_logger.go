package mylog

import (
	"io"
	"os"
	"strings"
	"time"
)

type FileLogger struct {
	position Position
	//	invokingNumber uint
}

const (
	Ldate = 1 << iota // the date in the local time zone: 2009/01/23
	Ltime             // the time in the local time zone: 01:23:23
)

var (
	LstdFlags           = Ldate | Ltime // initial values for the standard logger
	prefix              = ""
	skipNumber Position = 1
)

func New(out io.Writer, prefix string, flag int) *FileLogger {
	//	dateString := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0]
	//	errorFile, err := os.OpenFile("error-"+dateString+".log", os.O_WRONLY|os.O_CREATE, 0666)
	//	if err != nil {
	//		panic("open error log file failed")
	//	}
	//	errorLogger := log.New(errorFile, prefix, LstdFlags)
	return nil
}

func (logger *FileLogger) GetPosition() Position {
	return logger.position
}

func (logger *FileLogger) SetPosition(pos Position) {
	logger.position = pos
}

func writeContent(logType string, content *string) {
	var fileName string
	dateString := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0]
	switch logType {
	case "error":
		fileName = "error-" + dateString + ".log"
	case "fatal":
		fileName = "fatal-" + dateString + ".log"
	case "panic":
		fileName = "panic-" + dateString + ".log"
	case "info":
		fileName = "info-" + dateString + ".log"
	case "warn":
		fileName = "warn-" + dateString + ".log"
	}
	errorFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("open error log file failed")
	}
	defer errorFile.Close()
	errorFile.WriteString(*content + "\r\n")
}

func (logger *FileLogger) Error(v ...interface{}) string {
	content := generateLogContent(getErrorLogTag(), skipNumber, "", v...)
	writeContent("error", &content)
	return content
}

func (logger *FileLogger) Errorf(format string, v ...interface{}) string {
	content := generateLogContent(getErrorLogTag(), skipNumber, format, v...)
	writeContent("error", &content)
	return content
}

func (logger *FileLogger) Errorln(v ...interface{}) string {
	content := generateLogContent(getErrorLogTag(), skipNumber, "", v...)
	writeContent("error", &content)
	return content
}

func (logger *FileLogger) Fatal(v ...interface{}) string {
	content := generateLogContent(getFatalLogTag(), skipNumber, "", v...)
	writeContent("fatal", &content)
	return content
}

func (logger *FileLogger) Fatalf(format string, v ...interface{}) string {
	content := generateLogContent(getFatalLogTag(), skipNumber, format, v...)
	writeContent("fatal", &content)
	return content
}

func (logger *FileLogger) Fatalln(v ...interface{}) string {
	content := generateLogContent(getFatalLogTag(), skipNumber, "", v...)
	writeContent("fatal", &content)
	return content
}

func (logger *FileLogger) Info(v ...interface{}) string {
	content := generateLogContent(getInfoLogTag(), skipNumber, "", v...)
	writeContent("info", &content)
	return content
}

func (logger *FileLogger) Infof(format string, v ...interface{}) string {
	content := generateLogContent(getInfoLogTag(), skipNumber, format, v...)
	writeContent("info", &content)
	return content
}

func (logger *FileLogger) Infoln(v ...interface{}) string {
	content := generateLogContent(getInfoLogTag(), skipNumber, "", v...)
	writeContent("info", &content)
	return content
}

func (logger *FileLogger) Warn(v ...interface{}) string {
	content := generateLogContent(getWarnLogTag(), skipNumber, "", v...)
	writeContent("warn", &content)
	return content
}

func (logger *FileLogger) Warnf(format string, v ...interface{}) string {
	content := generateLogContent(getWarnLogTag(), skipNumber, format, v...)
	writeContent("warn", &content)
	return content
}

func (logger *FileLogger) Warnln(v ...interface{}) string {
	content := generateLogContent(getPanicLogTag(), skipNumber, "", v...)
	writeContent("warn", &content)
	return content
}

func (logger *FileLogger) Panic(v ...interface{}) string {
	content := generateLogContent(getPanicLogTag(), skipNumber, "", v...)
	writeContent("panic", &content)
	return content
}

func (logger *FileLogger) Panicf(format string, v ...interface{}) string {
	content := generateLogContent(getPanicLogTag(), skipNumber, format, v...)
	writeContent("panic", &content)
	return content
}

func (logger *FileLogger) Panicln(v ...interface{}) string {
	content := generateLogContent(getPanicLogTag(), skipNumber, "", v...)
	writeContent("panic", &content)
	return content
}

func (logger *FileLogger) SetDefaultInvokingNumber() {
	//	logger.invokingNumber = 1
	logger.position = 1
}

func (logger *FileLogger) SetInvokingNumber(invoking Position) {
	logger.position = invoking
}

func (logger *FileLogger) GetInvokingNumber() Position {
	return logger.position
}
