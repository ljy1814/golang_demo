package analyzer

import (
	"errors"
	"fmt"
	"reflect"
	"demo/webcrawler/middleware"
)

//生成分析器的函数
type  GenAnalyzer func() Analyzer

//分析池的接口
type AnalyzerPool interface {
	Take() (Analyzer, error) //取出一个分析器
	Return( analyzer Analyzer) error //归还分析器
	Total() uint32	//获取总量
	Used() uint32 //正在被使用的分析器数量
}

func NewAnalyzerPool(total uint32, gen GenAnalyzer) (AnalyzerPool, error) {
	etype := reflect.TypeOf(gen())
	genEntity := func() middleware.Entity {
		return gen()
	}
	pool, err := middleware.NewPool(total, etype, genEntity)
	if err != nil {
		return nil,  err
	}
	dlpool := &myAnalyzerPool{pool: pool, etype: etype}
	return dlpool, nil
}

type myAnalyzerPool struct {
	pool middleware.Pool //实体池
	etype reflect.Type //实体池类型
}

func (spdpool *myAnalyzerPool) Take() (Analyzer, error) {
	entity, err := spdpool.pool.Take()
	if err != nil {
		return nil, err
	}
	analyzer, ok := entity.(Analyzer)
	if !ok {
		errMsg := fmt.Sprintf("The type of entity is NOT %s!\n", spdpool.etype)
		panic(errors.New(errMsg))
	}
	return analyzer, nil
}

func (spdpool *myAnalyzerPool) Return(analyzer Analyzer) error {
	return spdpool.pool.Return(analyzer)
}

func (spdpool *myAnalyzerPool) Total() uint32 {
	return spdpool.pool.Total()
}

func (spdpool *myAnalyzerPool) Used() uint32 {
	return spdpool.pool.Used()
}
