# groupcache

## 官方介绍

### 总结
groupcache是一个缓存和缓存填充库，在许多情况下用于替代memcached。

### 与memcache相同支出
shard by key选择哪个对等体负责该密钥

### 对比memcache
不需要运行单独的一组服务器，从而大量减少部署/配置的痛苦。 groupcache是​​一个客户端库以及一个服务器。它连接到自己的同伴。

自带缓存填充机制。而memcached只是说“对不起，缓存未命中”，经常导致从无限数量的客户端（这导致几个有趣的停机）数据库（或任何）负载的群体，groupcache协调缓存填充，使只有一个负载整个复制的一组进程的一个进程填充高速缓存，然后将加载的值复用到所有呼叫者。

不支持版本化值。如果键“foo”是值“bar”，则键“foo”必须始终是“bar”。既没有缓存过期时间，也没有显式缓存逐出。因此也没有CAS，也没有增量/减少。这也意味着groupcache ....

...支持将超热项目自动镜像到多个进程。这防止了memcached热点，其中机器的CPU和/或NIC由非常流行的键/值过载。

目前只能用于Go。这是不太可能的，我（bradfitz @）将代码移植到任何其他语言。

### 进程加载
简而言之，Get（“foo”）的groupcache查找如下所示：

（在运行相同代码的一组N个机器的机器＃5上）

是本地内存中的“foo”的值，因为它是超级热？ 如果是，使用它。

是本地内存中的“foo”的值，因为对等体＃5（当前对等体）是它的所有者吗？ 如果是，使用它。

在我的一套N的所有同行中，我是钥匙“foo”的所有者？ （例如，它是否一致哈希到5？）如果是，加载它。 如果其他调用者通过相同的过程或通过对等体的RPC请求进入，它们阻止等待加载完成并获得相同的答案。 如果不是，RPC给作为所有者的对等体并得到答案。 如果RPC失败，只需在本地加载它（仍然使用本地dup抑制）。

### 代码总结

#### 哈希
```c
这是一批测试数据
hash.Add("6", "4", "2")
Display hash (*consistenthash.Map):
    (*hash).hash = consistenthash.Hash0x46e720
    (*hash).replicas = 3  //迭代3层;2,4,6增序排列
    (*hash).keys[0] = 2
    (*hash).keys[1] = 4
    (*hash).keys[2] = 6
    (*hash).keys[3] = 12
    (*hash).keys[4] = 14
    (*hash).keys[5] = 16
    (*hash).keys[6] = 22
    (*hash).keys[7] = 24
    (*hash).keys[8] = 26
    (*hash).hashMap[6] = "6"
    (*hash).hashMap[2] = "2"
    (*hash).hashMap[12] = "2"
    (*hash).hashMap[22] = "2"
    (*hash).hashMap[16] = "6"
    (*hash).hashMap[26] = "6"
    (*hash).hashMap[4] = "4"
    (*hash).hashMap[14] = "4"
    (*hash).hashMap[24] = "4"

testCases := map[string]string{
    "2":  "2", // 2 >= 2
    "11": "2", //12 > 11
    "23": "4", //24 > 23 故为4
    "27": "2", //27 > 26 故还原为0
}

hash1.Add("Bill", "Bob", "Bonny")
hash2.Add("Bob", "Bonny", "Bill")

Display hash1 (*consistenthash.Map):
    (*hash1).hash = consistenthash.Hash0x4bc1c0
    (*hash1).replicas = 1
    (*hash1).keys[0] = 1679827945
    (*hash1).keys[1] = 2622760538
    (*hash1).keys[2] = 3819440399
    (*hash1).hashMap[2622760538] = "Bill"
    (*hash1).hashMap[3819440399] = "Bob"
    (*hash1).hashMap[1679827945] = "Bonny"
Display hash2 (*consistenthash.Map):
    (*hash2).hash = consistenthash.Hash0x4bc1c0
    (*hash2).replicas = 1
    (*hash2).keys[0] = 1679827945
    (*hash2).keys[1] = 2622760538
    (*hash2).keys[2] = 3819440399
    (*hash2).hashMap[1679827945] = "Bonny"
    (*hash2).hashMap[2622760538] = "Bill"
    (*hash2).hashMap[3819440399] = "Bob"

//由于Ben的哈希值最小,而最小的hash值为键对于的值为Bonny
fmt.Println(hash1.Get("Ben")) //Bonny
fmt.Println(hash2.Get("Ben"))
fmt.Println(hash1.Get("Bob")) //Bob 对应自己
fmt.Println(hash2.Get("Bob"))
fmt.Println(hash1.Get("Bonny")) //Bonny  Bonny刚好对应自己
fmt.Println(hash2.Get("Bonny"))

哈希get返回结果
hash := int(m.hash([]byte(key)))
//返回哈希值刚刚大于等于当前哈希值的元素
idx := sort.Search(len(m.keys), func(i int) bool { return m.keys[i] >= hash })
if idx == len(m.keys) {
    idx = 0
}
return m.hashMap[m.keys[idx]]

```

#### 缓存设计
```c
type Cache struct {
    MaxEntries int                              //条目最多数,0表示无限制
    OnEvicted  func(key Key, value interface{}) //当清除元素时调用
    ll         *list.List
    cache      map[interface{}]*list.Element
}

//添加操作
func (c *Cache) Add(key Key, value interface{}) {
    //cache未初始化,
    if c.cache == nil {
        c.cache = make(map[interface{}]*list.Element)
        c.ll = list.New()
    }
    //存在该元素
    if ee, ok := c.cache[key]; ok {
        c.ll.MoveToFront(ee)
        ee.Value.(*entry).value = value
        return
    }
    //在列表最前边插入该元素
    ele := c.ll.PushFront(&entry{key, value})
    //在cache里面存储该元素的指针
    c.cache[key] = ele
    //缓存已满,缓存移除最早的元素,但是在列表里面并不移除,伪删除
    if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
        c.RemoveOldest()
    }
}


func (c *Cache) Get(key Key) (value interface{}, ok bool) {
    if c.cache == nil {
        return
    }
    if ele, hit := c.cache[key]; hit {
    //访问之后立即将它移到列表最前边
        c.ll.MoveToFront(ele)
        return ele.Value.(*entry).value, true
    }
    return
}
func (c *Cache) Remove(key Key) {
    if c.cache == nil {
        return
    }
    if ele, hit := c.cache[key]; hit {
        //缓存中删除映射关系,从链表删除
        c.removeElement(ele)
    }
}

func (c *Cache) removeElement(e *list.Element) {
    c.ll.Remove(e)
    kv := e.Value.(*entry)
    delete(c.cache, kv.key)
    if c.OnEvicted != nil {
        c.OnEvicted(kv.key, kv.value)
    }
}

```

#### 快速操作
```c
type Group struct {
    mu sync.Mutex
    m  map[string]*call //延迟初始化
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
    g.mu.Lock()
    if g.m == nil {
        g.m = make(map[string]*call)
    }
    //如果map中存在,直接返回现有数据
    if c, ok := g.m[key]; ok {
        g.mu.Unlock()
        c.wg.Wait()
        return c.val, c.err
    }
    //生成新的数据
    c := new(call)
    c.wg.Add(1)
    g.m[key] = c
    g.mu.Unlock()
    c.val, c.err = fn()
    c.wg.Done()
    g.mu.Lock()
    //删除map中的key
    delete(g.m, key)
    g.mu.Unlock()
    return c.val, c.err
}

```
#### 值存储
byte
string
proto
alloc
trunc

#### groupcache核心
```c
func (g *Group) Get(ctx Context, key string, dest Sink) error {
    g.peersOnce.Do(g.initPeers)
    //get计数加1
    g.Stats.Gets.Add(1)
    if dest == nil {
        return errors.New("groupcache: nil dest Sink")
    }
    value, cacheHit := g.lookupCache(key)
    //命中缓存
    if cacheHit {
        g.Stats.CacheHits.Add(1)
        return setSinkView(dest, value)
    }
    //两个cache都没有查找到
    //避免两次反序列化或者复制,多goroutine同时访问
    destPopulated := false
    value, destPopulated, err := g.load(ctx, key, dest)
    if err != nil {
        return err
    }
    if destPopulated {
        return nil
    }
    return setSinkView(dest, value)
}

//加载数据,本地加载或者从网络上获取
func (g *Group) load(ctx Context, key string, dest Sink) (value ByteView, destPopulated bool, err error) {
    g.Stats.Loads.Add(1)
    viewi, err := g.loadGroup.Do(key, func() (interface{}, error) {
        //查询cache,防止两个goroutine同时访问出现差错,singleflight只会调用那个函数
        /*
        1 Get("key")
        2 Get("key")
        1 lookupCache("key")
        2 lookupCache("key")
        1 load("key")
        2 load("key")
        1 loadGroup("key", fn)
        1 fn
        2 loadGroup("key", fn)
        2 fn
        */
        //防止有一个goroutine已经将其加载到缓存里面了
        if value, cacheHit := g.lookupCache(key); cacheHit {
            g.Stats.CacheHits.Add(1)
            return value, nil
        }
        g.Stats.LoadsDeduped.Add(1)
        var value ByteView
        var err error
        //获取网络上的对等体,同类客户端或者服务器
        if peer, ok := g.peers.PickPeer(key); ok {
            value, err = g.getFromPeer(ctx, peer, key)
            if err == nil {
                g.Stats.PeerLoads.Add(1)
                return value, nil
            }
            //其它对等体没有查找到,
            g.Stats.PeerErrors.Add(1)
        }
        //没有从网络上查找到,从本地查找
        value, err = g.getLocally(ctx, key, dest)
        if err != nil {
            g.Stats.LocalLoadErrors.Add(1)
            return nil, err
        }
        g.Stats.LocalLoads.Add(1)
        destPopulated = true
        g.populateCache(key, value, &g.mainCache)
        return value, nil
    })
    if err == nil {
        value = viewi.(ByteView)
    }
    return
}

func (g *Group) getFromPeer(ctx Context, peer ProtoGetter, key string) (ByteView, error) {
    req := &groupcachepb.GetRequest{
        Group: &g.name,
        Key:   &key,
    }
    res := &groupcachepb.GetResponse{}
    err := peer.Get(ctx, req, res)
    if err != nil {
        return ByteView{}, err
    }
    value := ByteView{b: res.Value}
    //随机清理cache中的元素
    if rand.Intn(10) == 0 {
        g.populateCache(key, value, &g.hotCache)
    }
    return value, nil
}


func (g *Group) lookupCache(key string) (value ByteView, ok bool) {
    if g.cacheBytes <= 0 {
        return
    }
    //先从mainCache里面查找,没有则去hotCache里面查找
    value, ok = g.mainCache.get(key)
    if ok {
        return
    }
    value, ok = g.hotCache.get(key)
    return
}


func (g *Group) populateCache(key string, value ByteView, cache *cache) {
    if g.cacheBytes <= 0 {
        return
    }
    cache.add(key, value)
    //Evict items from cache(s) if necessary
    for {
        mainBytes := g.mainCache.bytes()
        hotBytes := g.hotCache.bytes()
        //缓存未满
        if mainBytes+hotBytes <= g.cacheBytes {
            return
        }
        victim := &g.mainCache
        //到了一定程度就清理最旧的元素
        if hotBytes > mainBytes/8 {
            victim = &g.hotCache
        }
        victim.removeOldest()
    }
}
```

#### 网络对等体peer池
采用HTTPPool池组成一个pool,形成一个group

```c
//取出一个peer
func (p *HTTPPool) PickPeer(key string) (ProtoGetter, bool) {
    //  display.Display("HTTPPool", p)
    p.mu.Lock()
    defer p.mu.Unlock()
    if p.peers.IsEmpty() {
        return nil, false
    }
    //不是自身
    if peer := p.peers.Get(key); peer != p.self {
        return p.httpGetters[peer], true
    }
    return nil, false
}
```
