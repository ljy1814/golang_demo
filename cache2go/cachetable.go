package cache2go

import (
	"log"
	"sort"
	"sync"
	"time"
)

type CacheTable struct {
	sync.RWMutex

	name  string                     //表名
	items map[interface{}]*CacheItem //所有的items

	cleanupTimer    *time.Timer   //定时清理
	cleanupInterval time.Duration //清理时间间隔

	logger *log.Logger //日志

	loadData  func(key interface{}, args ...interface{}) *CacheItem //回调函数,当处理不存在的键
	addedItem func(item *CacheItem)                                 //添加一个新键值对

	aboutToDeleteItem func(item *CacheItem) //删除一个键值对
}

func (table *CacheTable) Count() int {
	table.RLock()
	defer table.RUnlock()
	return len(table.items)
}

//循环遍历所有键值对
func (table *CacheTable) Foreach(trans func(key interface{}, item *CacheItem)) {
	table.RLock()
	defer table.RUnlock()

	for k, v := range table.items {
		trans(k, v)
	}
}

//设置数据加载器
func (table *CacheTable) SetDataLoader(f func(interface{}, ...interface{}) *CacheItem) {
	table.Lock()
	defer table.Unlock()
	table.loadData = f
}

//设置添加键值对的函数
func (table *CacheTable) SetAddedItemCallback(f func(*CacheItem)) {
	table.Lock()
	defer table.Unlock()
	table.addedItem = f
}

//设置删除键值对函数
func (table *CacheTable) SetAboutToDeleteItemCallback(f func(*CacheItem)) {
	table.Lock()
	defer table.Unlock()
	table.aboutToDeleteItem = f
}

//设置日志记录器
func (table *CacheTable) SetLogger(logger *log.Logger) {
	table.Lock()
	defer table.Unlock()
	table.logger = logger
}

//删除过期的键值对
func (table *CacheTable) expirationCheck() {
	table.Lock()
	if table.cleanupTimer != nil {
		table.cleanupTimer.Stop()
	}
	if table.cleanupInterval > 0 {
		table.log("Expiration check triggered after", table.cleanupInterval, "for table ", table.name)
	} else {
		table.log("Expiration check installed for table", table.name)
	}
	//缓存value,不用打断数据同步
	items := table.items
	table.Unlock()

	//更新计时器时间
	now := time.Now()
	smallestDuration := 0 * time.Second
	for key, item := range items {
		//此处采用读写锁
		item.RLock()
		lifeSpan := item.lifeSpan
		accessedOn := item.accessedOn
		item.RUnlock()

		if lifeSpan == 0 {
			continue
		}
		if now.Sub(accessedOn) >= lifeSpan {
			//过期
			table.Delete(key)
		} else {
			if smallestDuration == 0 || lifeSpan-now.Sub(accessedOn) < smallestDuration {
				smallestDuration = lifeSpan - now.Sub(accessedOn)
			}
		}
	}

	//设置下一个时间间隔
	table.RLock()
	table.cleanupInterval = smallestDuration
	if smallestDuration > 0 {
		table.cleanupTimer = time.AfterFunc(smallestDuration, func() {
			go table.expirationCheck()
		})
	}
	table.RUnlock()
}

func (table *CacheTable) addInternal(item *CacheItem) {
	table.log("Add item with key", item.key, "and lifespan of", item.lifeSpan, "to table", table.name)
	table.items[item.key] = item

	//缓存value,因此打开锁,当未加锁时不要运行此方法
	expDur := table.cleanupInterval
	addedItem := table.addedItem
	table.Unlock()

	//添加键值对到缓存
	if addedItem != nil {
		addedItem(item)
	}

	if item.lifeSpan > 0 && (expDur == 0 || item.lifeSpan < expDur) {
		table.expirationCheck()
	}
}

//添加一个键值对,key是键,lifeSpan是超多这个时间间隔未访问就从cache移除的时间间隔,data是数据
func (table *CacheTable) Add(key interface{}, lifeSpan time.Duration, data interface{}) *CacheItem {
	item := NewCacheItem(key, lifeSpan, data)

	table.Lock()
	table.addInternal(item)
	return item
}

func (table *CacheTable) Delete(key interface{}) (*CacheItem, error) {
	table.RLock()
	r, ok := table.items[key]
	if !ok {
		table.RUnlock()
		return nil, ErrKeyNotFound
	}

	//缓存结构的值,因此不用保持加锁状态
	aboutToDeleteItem := table.aboutToDeleteItem
	table.RUnlock()

	if aboutToDeleteItem != nil {
		aboutToDeleteItem(r)
	}

	r.RLock()
	defer r.RUnlock()
	if r.aboutToExpire != nil {
		r.aboutToExpire(key)
	}

	table.Lock()
	defer table.Unlock()
	table.log("Delete item with key", key, "created on", r.createdOn, "and hit", r.accessCount, "time from table", table.name)
	delete(table.items, key)
	return r, nil
}

//检查key是否存在
func (table *CacheTable) Exists(key interface{}) bool {
	table.RLock()
	defer table.RUnlock()
	_, ok := table.items[key]
	return ok
}

//不存在时添加
func (table *CacheTable) NotFoundAdd(key interface{}, lifeSpan time.Duration, data interface{}) bool {
	table.Lock()
	if _, ok := table.items[key]; ok {
		table.Unlock()
		return false
	}

	item := NewCacheItem(key, lifeSpan, data)
	table.addInternal(item)
	return true
}

func (table *CacheTable) Value(key interface{}, args ...interface{}) (*CacheItem, error) {
	table.RLock()
	r, ok := table.items[key]
	loadData := table.loadData

	table.RUnlock()

	//存在这个键值对,更新其生命周期,
	if ok {
		r.KeepAlive()
		return r, nil
	}
	//如果不存在,尝试从数据加载器中寻找,找不到返回nil
	if loadData != nil {
		item := loadData(key, args...)
		if item != nil {
			table.Add(key, item.lifeSpan, item.data)
			return item, nil
		}
		//没有加载成功
		return nil, ErrKeyNotFoundOrLoaded
	}
	return nil, ErrKeyNotFound
}

//Flush删除所有的键值对
func (table *CacheTable) Flush() {
	table.Lock()
	defer table.Unlock()
	table.log("Flushing table", table.name)
	//重新初始化items,旧的Items被回收
	table.items = make(map[interface{}]*CacheItem)
	table.cleanupInterval = 0
	if table.cleanupTimer != nil {
		table.cleanupTimer.Stop()
	}
}

//key - 访问次数
type CacheItemPair struct {
	Key         interface{}
	AccessCount int64
}

type CacheItemPairList []CacheItemPair //实现了sort接口,通过AccessCount排序

func (p CacheItemPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p CacheItemPairList) Len() int           { return len(p) }
func (p CacheItemPairList) Less(i, j int) bool { return p[i].AccessCount > p[j].AccessCount }

//访问最频繁的键值对
func (table *CacheTable) MostAccessed(count int64) []*CacheItem {
	table.RLock()
	defer table.RUnlock()

	p := make(CacheItemPairList, len(table.items))
	i := 0
	for k, v := range table.items {
		p[i] = CacheItemPair{k, v.accessCount}
		i++
	}
	//根据访问次数进行排序
	sort.Sort(p)

	var r []*CacheItem
	c := int64(0)
	for _, v := range p {
		if c >= count {
			break
		}

		item, ok := table.items[v.Key]
		if ok {
			r = append(r, item)
		}
		c++
	}
	return r
}

//日志记录
func (table *CacheTable) log(v ...interface{}) {
	if table.logger == nil {
		return
	}
	table.logger.Println(v)
}
