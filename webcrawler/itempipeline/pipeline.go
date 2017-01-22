package itempipeline

import (
	"errors"
	"fmt"
	"sync/atomic"
	"demo/webcrawler/base"
)


//条目处理管道接口
type ItemPipeline interface {
	//发送条目
	Send(item base.Item) []error
	//返回一个布尔值,用于判断当前条目处理管道是否快速失败
	//快速失败:只有在条目处理流程中某一个步骤出错
	FailFast() bool
	//设置是否快速失败
	SetFailFast(failFast bool)
	//获取已发送,已接收和已处理的条目的计数值.作为结果值的切片总会有3个元素,这3个值分别代表前述的3个计数
	Count() []uint64
	//获取正在处理的条目数量
	ProcessingNumber() uint64
	//获取摘要信息
	Summary() string
}

//创建条目处理管道
func NewItemPipeline(itemProcessors []ProcessItem) ItemPipeline {
	if itemProcessors == nil {
		panic(errors.New(fmt.Sprintf("Invalid item processor list!")))
	}
	innerItemProcessors := make([]ProcessItem, 0)
	for i, ip := range itemProcessors {
		if ip == nil {
			panic(errors.New(fmt.Sprintf("Invalid item processor[%d]!\n", i)))
		}
		innerItemProcessors = append(innerItemProcessors, ip)
	}
	return &myItemPipeline{itemProcessors: innerItemProcessors}
}

//条目处理管道
type myItemPipeline struct {
	itemProcessors []ProcessItem //条目处理器列表
	failFast bool //是否需要快速失败标志位
	sent uint64 //已发送的条目数量
	accepted uint64 //已接受的条目数量
	processed uint64 //已处理的条目数量
	processingNumber uint64 //正在被处理的条目数量
}

func (ip *myItemPipeline) Send(item base.Item) []error {
	atomic.AddUint64(&ip.processingNumber, 1)
	defer atomic.AddUint64(&ip.processingNumber, uint64(0))
	atomic.AddUint64(&ip.sent, 1)
	errs := make([]error, 0)
	if item == nil {
		errs = append(errs, errors.New("The item is invalid!"))
		return errs
	}
	atomic.AddUint64(&ip.accepted, 1)
	var currentItem base.Item = item
	for _, itemProcessor := range ip.itemProcessors {
		processItem, err := itemProcessor(currentItem) 
		if err != nil {
			errs = append(errs, err)
			if ip.failFast {
				break
			}
		}
		if  processItem != nil {
			currentItem = processItem
		}
	}
	atomic.AddUint64(&ip.processed, 1)
	return errs
}

func (ip *myItemPipeline) FailFast() bool {
	return ip.failFast
}

func (ip *myItemPipeline) SetFailFast(failFast bool) {
	ip.failFast = failFast
}

func (ip *myItemPipeline) Count() []uint64 {
	//记录3个数量
	counts := make([]uint64, 3)
	counts[0] = atomic.LoadUint64(&ip.sent)
	counts[1] = atomic.LoadUint64(&ip.accepted)
	counts[2] = atomic.LoadUint64(&ip.processed)
	return counts
}

func (ip *myItemPipeline) ProcessingNumber() uint64 {
	return atomic.LoadUint64(&ip.processingNumber)
}

var summaryTemplate = "failFast: %v, processorNumber: %d, sent:%d, accepted:%d, processed: %d"

func (ip *myItemPipeline) Summary() string {
	counts := ip.Count()
	summary := fmt.Sprintf(summaryTemplate, ip.failFast, len(ip.itemProcessors), counts[0], counts[1], counts[2], ip.ProcessingNumber())
	return summary
}
