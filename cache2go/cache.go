package cache2go

import "sync"

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

//根据table名返回table,如果不存在则创建一个table
func Cache(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		t = &CacheTable{
			name:  table,
			items: make(map[interface{}]*CacheItem),
		}

		mutex.Lock()
		cache[table] = t
		mutex.Unlock()
	}
	return t
}
