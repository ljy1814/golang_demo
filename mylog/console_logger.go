package mylog

import (
	"log"
)

type ConsoleLogger struct {
	position Position
	//	invokingNumber uint
}

func (logger *ConsoleLogger) GetPosition() Position {
	return logger.position
}

func (logger *ConsoleLogger) SetPosition(pos Position) {
	logger.position = pos
}

func (logger *ConsoleLogger) Error(v ...interface{}) string {
	content := generateLogContent(getErrorLogTag(), logger.GetPosition(), "", v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Errorf(format string, v ...interface{}) string {
	content := generateLogContent(getErrorLogTag(), logger.GetPosition(), format, v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Errorln(v ...interface{}) string {
	content := generateLogContent(getErrorLogTag(), logger.GetPosition(), "", v...)
	log.Println(content)
	return content
}

func (logger *ConsoleLogger) Fatal(v ...interface{}) string {
	content := generateLogContent(getFatalLogTag(), logger.GetPosition(), "", v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Fatalf(format string, v ...interface{}) string {
	content := generateLogContent(getFatalLogTag(), logger.GetPosition(), format, v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Fatalln(v ...interface{}) string {
	content := generateLogContent(getFatalLogTag(), logger.GetPosition(), "", v...)
	log.Println(content)
	return content
}

func (logger *ConsoleLogger) Info(v ...interface{}) string {
	content := generateLogContent(getInfoLogTag(), logger.GetPosition(), "", v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Infof(format string, v ...interface{}) string {
	content := generateLogContent(getInfoLogTag(), logger.GetPosition(), format, v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Infoln(v ...interface{}) string {
	content := generateLogContent(getInfoLogTag(), logger.GetPosition(), "", v...)
	log.Println(content)
	return content
}

func (logger *ConsoleLogger) Warn(v ...interface{}) string {
	content := generateLogContent(getWarnLogTag(), logger.GetPosition(), "", v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Warnf(format string, v ...interface{}) string {
	content := generateLogContent(getWarnLogTag(), logger.GetPosition(), format, v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Warnln(v ...interface{}) string {
	content := generateLogContent(getPanicLogTag(), logger.GetPosition(), "", v...)
	log.Println(content)
	return content
}

func (logger *ConsoleLogger) Panic(v ...interface{}) string {
	content := generateLogContent(getPanicLogTag(), logger.GetPosition(), "", v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Panicf(format string, v ...interface{}) string {
	content := generateLogContent(getPanicLogTag(), logger.GetPosition(), format, v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Panicln(v ...interface{}) string {
	content := generateLogContent(getPanicLogTag(), logger.GetPosition(), "", v...)
	log.Println(content)
	return content
}

func (logger *ConsoleLogger) SetDefaultInvokingNumber() {
	//	logger.invokingNumber = 1
	logger.position = 1
}

func (logger *ConsoleLogger) SetInvokingNumber(invoking Position) {
	logger.position = invoking
}

func (logger *ConsoleLogger) GetInvokingNumber() Position {
	return logger.position
}

//TODO 写一个带color的控制台日志工具
//red \x1b[%dm%s\x1b[0m]]
//green \x1b[%dm%s\x1b[0m]]
//yellow \x1b[%dm%s\x1b[0m]]
//blue \x1b[%dm%s\x1b[0m]]
//mageenta \x1b[%dm%s\x1b[0m]] 洋红
