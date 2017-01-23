package downloader
import (
	"errors"
	"fmt"
	"reflect"
	"demo/webcrawler/middleware"
)

//生成网页下载器
type GenPageDownloader func() PageDownloader

//网页下载器池
type PageDownloaderPool interface {
	Take() (PageDownloader, error) //取出一个网页下载器
	Return(dl PageDownloader) error 
	Total() uint32 //获取池子总量
	Used() uint32	//正在使用的网页下载器总量 
}

func NewPageDownloaderPool(total uint32, gen GenPageDownloader) (PageDownloaderPool, error) {
	etype := reflect.TypeOf(gen())
	genEntity := func() middleware.Entity {
		return gen()
	}
	pool, err := middleware.NewPool(total, etype, genEntity)
	if err != nil {
		return nil, err
	}
	dlpool := &myDownloaderPool { pool: pool, etype: etype }
	return dlpool, nil
}

type myDownloaderPool struct {
	pool middleware.Pool //实体池
	etype reflect.Type //实体池的类型
}


func (dlpool *myDownloaderPool) Take() (PageDownloader, error) {
	entity, err := dlpool.pool.Take()
	if err != nil {
		return nil, err
	}
	dl, ok := entity.(PageDownloader)
	if !ok {
		errMsg := fmt.Sprintf("The type of entity is NOT %s\n", dlpool.etype)
		panic(errors.New(errMsg))
	}
	return dl, nil
}

func (dlpool *myDownloaderPool) Return(dl PageDownloader) error {
	return dlpool.pool.Return(dl)
}

func (dlpool *myDownloaderPool) Total() uint32 {
	return dlpool.pool.Total()
}

func (dlpool *myDownloaderPool) Used() uint32 {
	return dlpool.pool.Used()
}
