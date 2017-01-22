package middleware

import (
	"math"
	"sync"
)

//ID生成器
type IdGenerator interface {
	GetUint32() uint32 //获得一个uint32类型的ID
}

//创建ID生成器
func NewIdGenerator() IdGenerator {
	return &cycliIdGenerator{}
}

//ID生成器
type cycliIdGenerator struct {
	sn uint32 //当前ID
	ended bool //是否达到uint32最大值
	mutex sync.Mutex //互斥锁
}

func (gen *cycliIdGenerator) GetUint32() uint32 {
	gen.mutex.Lock()
	defer gen.mutex.Unlock()
	if gen.ended {
		defer func() {
			gen.ended = false
		}()
		gen.sn = 0
		return gen.sn
	}
	id := gen.sn
	if id < math.MaxUint32 {
		gen.sn++
	} else {
		gen.ended = true
	}
	return id
}

//64位ID生成器接口
type IdGenerator64 interface {
	GetUint64() uint64 //获得一个uint64类型的ID
}

//创建64位ID生成器
func NewIdGenerator64() IdGenerator64 {
	return &cycliIdGenerator64{}
}

type cycliIdGenerator64 struct {
	base cycliIdGenerator //基本的ID生成器
	cycleCount uint64 //基于uint32类型的取值范围的周期计数
}

func (gen *cycliIdGenerator64) GetUint64() uint64 {
	var id64 uint64
	if gen.cycleCount % 2 == 1 {
		id64 += math.MaxUint32
	}
	id32 := gen.base.GetUint32()
	if id32 == math.MaxUint32 {
		gen.cycleCount++
	}
	id64 += uint64(id32)
	return id64
}
