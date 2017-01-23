package scheduler

import (
	"fmt"
	"sync"
	"demo/webcrawler/base"
)

//状态字典
var statusMap = map[requestCacheStatus]string{
	0: "running",
	1: "closed",
}

type requestCacheStatus byte

const (
	REQUEST_CACHE_STATUS_RUNNING requestCacheStatus = 0
	REQUEST_CACHE_STATUS_CLOSED requestCacheStatus = 1
)

//请求缓存的接口类型
type requestCache interface {
	//将请求放入缓存
	put(req *base.Request) bool
	//从请求缓存获取最早放入且仍在其中的请求
	get() *base.Request
	//获取请求缓存的容量
	capacity() int
	//获取请求缓存实际长度,即请求的及时数量
	length() int
	//关闭请求
	close()
	//获取请求缓存的摘要信息
	summary() string
}

//创建请求缓存
func newRequestCache() requestCache {
	rc := &reqCacheBySlice{
		cache: make([]*base.Request, 0),
	}
	return rc
}

//请求缓存的实现类型
type reqCacheBySlice struct {
	cache []*base.Request //请求的存储介质
	mutex sync.Mutex  //互斥锁
	status requestCacheStatus //缓存状态,0-正在运行,1-已关闭
}

func (rcache *reqCacheBySlice) put(req *base.Request) bool {
	if req == nil {
		return false
	}
	if rcache.status == REQUEST_CACHE_STATUS_CLOSED {
		return false
	}
	rcache.mutex.Lock()
	defer rcache.mutex.Unlock()
	rcache.cache = append(rcache.cache, req)
	return true
}

func (rcache *reqCacheBySlice) get() *base.Request {
	if rcache.length() == 0 {
		return nil
	}
	if rcache.status == REQUEST_CACHE_STATUS_CLOSED {
		return nil
	}
	rcache.mutex.Lock()
	defer rcache.mutex.Unlock()
	req := rcache.cache[0]
	rcache.cache = rcache.cache[1:]
	return req
}

func (rcache *reqCacheBySlice) capacity() int {
	return cap(rcache.cache)
}

func (rcache *reqCacheBySlice) length() int {
	return len(rcache.cache)
}

func (rcache *reqCacheBySlice) close() {
	if rcache.status == REQUEST_CACHE_STATUS_CLOSED {
		return
	}
	rcache.status = REQUEST_CACHE_STATUS_CLOSED
}

//摘要信息模板
var summaryTemplate = "status: %s, length: %d, capacity: %d"

func (rcache *reqCacheBySlice) summary() string {
	summary := fmt.Sprintf(summaryTemplate, statusMap[rcache.status],
		rcache.length(), rcache.capacity())
	fmt.Printf("xxx-----%s\n", summary)
	return summary
}
