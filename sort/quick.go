package main

import (
	"fmt"
	"math/rand"

	"github.com/fatih/color"
)

type Item interface {
	Less(other Item) bool
}

type IItem int8

func (i IItem) Less(other Item) bool {
	return i < other.(IItem)
}

func partition(items []IItem, l, r int) int {
	item := items[l]
	color.Blue("l : %d , r : %d %v\n", l, r, l < r)
	for l < r {
		for item.Less(items[r]) {
			r--
		}
		if l < r {
			items[l], items[r] = items[r], items[l]
			l++
		}
		for items[l].Less(item) {
			l++
		}
		if l < r {
			items[l], items[r] = items[r], items[l]
			r--
		}
	}
	return l
}

func quickSort(items []IItem, l, r int) {
	if len(items) > 0 {
		color.Green("len[%d], cap[%d]\n", len(items), cap(items))
		mid := partition(items, l, r)
		color.Red("partition : [%v]\n", mid)
		if mid > l {
			quickSort(items, l, mid-1)
			quickSort(items, mid+1, r)
		}
	}
}

func main() {
	fmt.Println("vim-go")
	items := make([]IItem, 0)
	for i := 0; i < 20; i++ {
		tt := rand.Int()
		tt = tt % 128
		items = append(items, IItem(int8(tt)))
	}
	fmt.Println(items)
	fmt.Println(items[0:2])
	quickSort(items, 0, len(items)-1)
	fmt.Println(items)
}
