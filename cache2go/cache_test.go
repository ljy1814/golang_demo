package cache2go

import (
	"bytes"
	"demo/display"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	k = "testkey"
	v = "testvalue"
)

func TestCache(t *testing.T) {
	table := Cache("testCache")
	table.Add(k+"_1", 0*time.Second, v+"_1")
	table.Add(k+"_2", 1*time.Second, v+"_2")

	p, err := table.Value(k + "_1")
	if err != nil || p == nil || p.Data().(string) != v+"_1" {
		t.Error("Error retrieving non expiring data from cache", err)
	}

	p, err = table.Value(k + "_2")
	if err != nil || p == nil || p.Data().(string) != v+"_2" {
		t.Error("Error retrieving non expiring data from cache", err)
	}

	if p.AccessCount() != 1 {
		t.Error("Error getting correct access count")
	}
	if p.LifeSpan() != 1*time.Second {
		t.Error("Error getting access time")
	}
	if p.AccessedOn().Unix() == 0 {
		t.Error("Error getting access time")
	}
	if p.CreatedOn().Unix() == 0 {
		t.Error("Error getting creation time")
	}
	display.Display("cache", table)
}

func TestCacheExpire(t *testing.T) {
	table := Cache("testCache")
	table.Add(k+"_1", 0*time.Second, v+"_1")
	table.Add(k+"_2", 1*time.Second, v+"_2")

	time.Sleep(75 * time.Millisecond)
	_, err := table.Value(k + "_1")
	if err != nil {
		t.Error("Error retrieving value from cache:", err)
	}
	_, err = table.Value(k + "_2")
	if err != nil {
		t.Error("Found key which should have been expired by now")
	}
}

func TestExists(t *testing.T) {
	table := Cache("testCache")
	table.Add(k, 0, v)
	if !table.Exists(k) {
		t.Error("Error verifying existing data in cache")
	}
}

func TestNotFoundAdd(t *testing.T) {
	table := Cache("testCache")
	if !table.NotFoundAdd(k, 0, v) {
		t.Error("Error verifying NotFoundAdd, data not in cache")
	}
	if table.NotFoundAdd(k, 0, v) {
		t.Error("Error verifying NotFoundAdd")
	}
}

func TestNotFoundAddConcurrency(t *testing.T) {
	table := Cache("testNotFoundAdd")
	var finish sync.WaitGroup
	var added int32
	var idle int32
	fn := func(id int) {
		for i := 0; i < 100; i++ {
			if table.NotFoundAdd(i, 0, i+id) {
				atomic.AddInt32(&added, 1)
			} else {
				atomic.AddInt32(&idle, 1)
			}
			time.Sleep(0)
		}
		finish.Done()
	}
	finish.Add(10)
	go fn(0x0000)
	go fn(0x1100)
	go fn(0x2200)
	go fn(0x3300)
	go fn(0x4400)
	go fn(0x5500)
	go fn(0x6600)
	go fn(0x7700)
	go fn(0x8800)
	go fn(0x9900)
	finish.Wait()

	t.Log(added, idle)
	table.Foreach(func(key interface{}, item *CacheItem) {
		v, _ := item.Data().(int)
		k, _ := key.(int)
		t.Logf("%02x %04x\n", k, v)
	})
}

func TestCacheKeepAlive(t *testing.T) {
	table := Cache("testKeepAlive")
	p := table.Add(k, 100*time.Millisecond, v)

	time.Sleep(50 * time.Millisecond)
	p.KeepAlive()

	time.Sleep(75 * time.Millisecond)
	if !table.Exists(k) {
		t.Error("Error expiring item after keeping it alive")
	}

	time.Sleep(75 * time.Millisecond)
	if table.Exists(k) {
		t.Error("Error expiring item after keeping it alive")
	}
}

func TestDelete(t *testing.T) {
	table := Cache("testDelete")
	table.Add(k, 0, v)

	p, err := table.Value(k)
	if err != nil || p == nil || p.Data().(string) != v {
		t.Error("Error retrieving data from cache", err)
	}

	table.Delete(k)
	p, err = table.Value(k)
	if err == nil || p != nil {
		t.Error("Error deleting data")
	}

	_, err = table.Delete(k)
	if err == nil {
		t.Error("Excepted error deleting item")
	}
}

func TestFlush(t *testing.T) {
	table := Cache("testFlush")
	table.Add(k, 10*time.Second, v)
	table.Flush()
	p, err := table.Value(k)
	if err == nil || p != nil {
		t.Error("Error flushing table")
	}
	if table.Count() != 0 {
		t.Error("Error verifying count of flushed table")
	}
}

func TestCount(t *testing.T) {
	table := Cache("testCount")
	count := 100000
	for i := 0; i < count; i++ {
		key := k + strconv.Itoa(i)
		table.Add(key, 10*time.Second, v)
	}

	for i := 0; i < count; i++ {
		key := k + strconv.Itoa(i)
		p, err := table.Value(key)
		if err != nil || p == nil || p.Data().(string) != v {
			t.Error("Error retrieving data")
		}
	}

	if table.Count() != count {
		t.Error("Data count mismatc")
	}
}

func TestDataLoader(t *testing.T) {
	table := Cache("testDataLoader")
	table.SetDataLoader(func(key interface{}, args ...interface{}) *CacheItem {
		var item *CacheItem
		if key.(string) != "nil" {
			val := k + key.(string)
			i := NewCacheItem(key, 500*time.Millisecond, val)
			item = i
		}
		return item
	})

	p, err := table.Value("nil")
	if err == nil || table.Exists("nil") {
		t.Error("Error validating data loader for nil values")
	}

	for i := 0; i < 10; i++ {
		key := k + strconv.Itoa(i)
		vp := k + key
		p, err = table.Value(key)
		if err != nil || p == nil || p.Data().(string) != vp {
			t.Error("Error validating data loader")
		}
	}
}

func TestAccessCount(t *testing.T) {
	count := 100
	table := Cache("testAccessCount")
	for i := 0; i < count; i++ {
		table.Add(i, 10*time.Second, v)
	}

	for i := 0; i < count; i++ {
		for j := 0; j < i; j++ {
			table.Value(i)
		}
	}

	ma := table.MostAccessed(int64(count))
	for i, item := range ma {
		if item.Key() != count-1-i {
			t.Error("Most accessed items seem to be sorted incorrectly")
		}
	}
	ma = table.MostAccessed(int64(count - 1))
	if len(ma) != count-1 {
		t.Error("MostAccessed returns incorrect amount of items")
	}
}

func TestCallbacks(t *testing.T) {
	addedKey := ""
	removedKey := ""
	expired := false

	table := Cache("testCallbacks")
	table.SetAddedItemCallback(func(item *CacheItem) {
		addedKey = item.Key().(string)
	})
	table.SetAboutToDeleteItemCallback(func(item *CacheItem) {
		removedKey = item.Key().(string)
	})
	i := table.Add(k, 500*time.Millisecond, v)
	i.SetAboutToExpireCallback(func(key interface{}) {
		expired = true
	})

	time.Sleep(250 * time.Millisecond)
	if addedKey != k {
		t.Error("AddedItem callback not working")
	}

	time.Sleep(500 * time.Millisecond)
	if removedKey != k {
		t.Error("AboutToDeleteItem callback not working:" + k + "_" + removedKey)
	}

	if !expired {
		t.Error("AboutToExpire callback not working")
	}
}

func TestLogger(t *testing.T) {
	out := new(bytes.Buffer)
	l := log.New(out, "cache2go", log.Ldate|log.Ltime)

	table := Cache("testLogger")
	table.SetLogger(l)
	table.Add(k, 1*time.Second, v)

	if out.Len() == 0 {
		t.Error("Logger is empty")
	}
}
