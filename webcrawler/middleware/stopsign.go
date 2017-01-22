package middleware

import (
	"fmt"
	"sync"
)

//停止信号
type StopSign interface {
	//发送停止信号,如果先前已经发送,返回false
	Sign() bool
	//判断停止信号是否被发出
	Signed() bool 
	//处理停止信号,参数code对应停止信号处理方的代号,该代号会出现在停止信号的处理记录中
	Deal(code string)
	//获取一个停止信号处理方的处理计数,该处理计数会从相应的停止信号处理记录中获取
	DealCount(code string) uint32
	//获取摘要信息,其中包括所以的停止信号处理记录
	Summary() string
	Reset()
}

//创建停止信号
func NewStopSign() StopSign {
	ss := &myStopSign {
		dealCountMap: make(map[string]uint32),
	}
	return ss
}
//停止信号实现
type myStopSign struct  {
	rwmutex sync.RWMutex //读写锁
	signed bool //表示信号是否已经发出的标志
	dealCountMap map[string]uint32 //处理计数的字典
}

func (ss *myStopSign) Sign() bool {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	if ss.signed {
		return false
	}
	ss.signed = true
	return true
}

func (ss *myStopSign) Signed() bool {
	return ss.signed
}

func (ss *myStopSign) Reset() {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	ss.signed = false
	ss.dealCountMap = make(map[string]uint32)
}

func (ss *myStopSign) Deal(code string) {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	if ss.signed {
		return
	}
	if _, ok := ss.dealCountMap[code]; !ok {
		ss.dealCountMap[code] = 1
	} else {
		ss.dealCountMap[code] += 1
	}
}

func (ss *myStopSign) DealCount(code string) uint32 {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	return ss.dealCountMap[code]
}

func (ss *myStopSign) Summary() string {
	if ss.signed {
		return fmt.Sprintf("signed: true, dealCount: %v", ss.dealCountMap)
	} else {
		return "signed: false"
	}
}
