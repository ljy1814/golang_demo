package itempipeline

import (
	"demo/webcrawler/base"
)

//处理条目的函数
type ProcessItem func(item base.Item) (result base.Item, err error)
