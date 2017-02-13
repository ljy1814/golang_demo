//hash操作
package consistenthash

import (
	"demo/display"
	"hash/crc32"
	"os"
	"sort"
	"strconv"
	"time"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int            //重复层数
	keys     []int          //Sorted,排序
	hashMap  map[int]string //哈希存储
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//是否是空哈希表
func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

//添加一些key
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
	GDisplay("m", m)
}

func GDisplay(name string, value interface{}) {
	fdisplay, err := os.OpenFile("display-"+time.Now().Format("2006-01-02")+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic("open display filoe failed")
	}
	defer fdisplay.Close()
	display.Fdisplay(fdisplay, name, value)
}

//获取hash值与之最接近的key
func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool { return m.keys[i] >= hash })
	if idx == len(m.keys) {
		idx = 0
	}
	return m.hashMap[m.keys[idx]]
}
