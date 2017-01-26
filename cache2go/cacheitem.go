package cache2go

import (
	"sync"
	"time"
)

type CacheItem struct {
	sync.RWMutex
	key           interface{}
	data          interface{}
	lifeSpan      time.Duration //在缓存中的存在时间
	createdOn     time.Time     //创建时间
	accessedOn    time.Time     //最近访问时间
	accessCount   int64         //访问次数
	aboutToExpire func(key interface{})
}

func NewCacheItem(key interface{}, lifeSpan time.Duration, data interface{}) *CacheItem {
	t := time.Now()
	return &CacheItem{
		key:           key,
		lifeSpan:      lifeSpan,
		createdOn:     t,
		accessedOn:    t,
		accessCount:   0,
		aboutToExpire: nil,
		data:          data,
	}
}

func (item *CacheItem) KeepAlive() {
	item.Lock()
	defer item.Unlock()
	item.accessedOn = time.Now()
	item.accessCount++
}

func (item *CacheItem) LifeSpan() time.Duration {
	//存活时间不需要同步,并发访问也可
	return item.lifeSpan
}

func (item *CacheItem) AccessedOn() time.Time {
	item.RLock()
	defer item.RUnlock()
	return item.accessedOn
}

func (item *CacheItem) AccessCount() int64 {
	item.RLock()
	defer item.RUnlock()
	return item.accessCount
}

func (item *CacheItem) CreatedOn() time.Time {
	return item.createdOn
}

func (item *CacheItem) Key() interface{} {
	return item.key
}

func (item *CacheItem) Data() interface{} {
	return item.data
}

func (item *CacheItem) SetAboutToExpireCallback(f func(interface{})) {
	item.Lock()
	defer item.Unlock()
	item.aboutToExpire = f
}
