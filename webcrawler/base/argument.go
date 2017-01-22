package base

import (
	"errors"
	"fmt"
)

//参数容器接口
type Args interface {
	//自检参数的有效性,并在必要时返回可以说明问题的错误值
	//若结果值为nil,则说明未发现问题,否则就意味着自检通过
	Check() error
	//获得参数容器的字符串表现形式
	String() string
}

//通道参数的容器描述版本
var channelArgsTemplate string = "{ reqChanLen: %d, respChanLen: %d.,  itemChanLen: %d, errChanLen: %d }"

//通道参数容器
type ChannelArgs struct {
	reqChanLen	uint //请求通道的长度
	respChanLen	uint //响应通道的长度
	itemChanLen	uint //条目通道的长度
	errorChanLen	uint //错误通道的长度
	description	string	//描述
} 

//创建通道参数的容器
func NewChannelArgs(
	reqChanLen uint,
	respChanLen uint,
	itemChanLen uint,
	errorChanLen uint) ChannelArgs {
	return ChannelArgs{
		reqChanLen: reqChanLen,
		respChanLen: respChanLen,
		errorChanLen: errorChanLen,
		itemChanLen: itemChanLen,
	}
}

//
func (args *ChannelArgs) Check() error {
	if args.reqChanLen == 0 {
		return errors.New("The request channel max length (capacity) can not be 0!.\n")
	}
	if args.respChanLen == 0 {
		return errors.New("The response channel max length (capacity) can not be 0!.\n")
	}
	if args.itemChanLen == 0 {
		return errors.New("The item channel max length (capacity) can not be 0!.\n")
	}
	if args.errorChanLen == 0 {
		return errors.New("The error channel max length (capacity) can not be 0!.\n")
	}
	return nil
}

//字符串函数
func (args *ChannelArgs) String() string {
	if args.description == "" {
		args.description = fmt.Sprintf(
			channelArgsTemplate,
			args.reqChanLen,
			args.respChanLen,
			args.itemChanLen,
			args.errorChanLen)
	}
	return args.description
}

//获取请求通道的长度
func (args *ChannelArgs) ReqChanLen() uint {
	return args.reqChanLen
}
//获取响应通道的长度
func (args *ChannelArgs) RespChanLen() uint {
	return args.respChanLen
}
//获取错误通道的长度
func (args *ChannelArgs) ErrorChanLen() uint {
	return args.errorChanLen
}
//获取通道的长度
func (args *ChannelArgs) ItemChanLen() uint {
	return args.itemChanLen
}

var poolBaseArgsTemplate string = "{ pageDownloaderPoolSize: %d, analyzerPoolSize: %d }"

//池基本参数
type PoolBaseArgs struct {
	pageDownloaderPoolSize	uint32//网页下载池尺寸
	analyzerPoolSize	uint32	//分析池的尺寸
	description		string //描述
}

//创建池基本参数的容器
func NewPoolBaseArgs(pageDownloaderPoolSize uint32,
	analyzerPoolSize uint32) PoolBaseArgs {
	return PoolBaseArgs{
		pageDownloaderPoolSize: pageDownloaderPoolSize,
		analyzerPoolSize: analyzerPoolSize,
	}
}

func (args *PoolBaseArgs) Check() error {
	if args.pageDownloaderPoolSize == 0 {
		return errors.New("The page downloader pool size can not be 0.\n")
	}
	if args.analyzerPoolSize == 0 {
		return errors.New("The analyzer pool size can not be 0.\n")
	}
	return nil
}

func (args *PoolBaseArgs) String() string {
	if args.description == "" {
		args.description = fmt.Sprintf(poolBaseArgsTemplate, args.pageDownloaderPoolSize,
		args.analyzerPoolSize)
	}
	return args.description
}

//获得下载器池的尺寸
func (args *PoolBaseArgs) PageDownloaderPoolSize() uint32 {
	return args.pageDownloaderPoolSize
}
//获得分析器池的尺寸
func (args *PoolBaseArgs) AnalyzerPoolSize() uint32 {
	return args.analyzerPoolSize
}
