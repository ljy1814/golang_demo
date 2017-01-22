package base

import (
	"bytes"
	"fmt"
)

//error type
type ErrorType string

//错误类型
const (
	DOWNLOADER_ERROR	ErrorType = "Downloader Error"
	ANALYZER_ERROR	ErrorType = "Analyzer Error"
	ITEM_PROCESSER_ERROR	ErrorType = "Item Processor Error"
)

//爬虫错误处理
type CrawlerError interface {
	Type() ErrorType  //错误除磷类型
	Error() string	  //错误提示信息
}

//爬虫错误结构
type myCrawlerError struct {
	errType ErrorType //错误类型
	errMsg  string		//错误提示信息
	fullErrMsg	string	//完整的错误提示信息
}

//创建新的爬虫错误提示
func NewCrawlerError(errType ErrorType, errMsg string) CrawlerError {
	return &myCrawlerError{errType: errType, errMsg: errMsg}
}


//获取错误类型
func (ce *myCrawlerError) Type() ErrorType {
	return ce.errType
}

//获取错误提示信息
func (ce *myCrawlerError) Error() string {
	if ce.fullErrMsg == "" {
		ce.genFullErrMsg()
	}
	return ce.fullErrMsg
}

//生成错误提示信息,并给相应的字段赋值
func (ce *myCrawlerError) genFullErrMsg() {
	var buffer  bytes.Buffer
	buffer.WriteString("Crawler Error: ")
	if ce.errType != "" {
		buffer.WriteString(string(ce.errType))
		buffer.WriteString(": ")
	}
	buffer.WriteString(ce.errMsg)
	ce.fullErrMsg = fmt.Sprintf("%s\n", buffer.String())
	return
}
