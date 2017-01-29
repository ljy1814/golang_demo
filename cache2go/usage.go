package main

import (
	"cache2go"
	"fmt"
	"time"
)

type myStruct struct {
	text     string
	moreData []byte
}

func main() {
	cache := cache2go.Cache("myCache")

	val := myStruct{"This is a test!", []byte{}}

	cache.Add("simeKey", 5*time.Second, &val)

	res, err := cache.Value("someKey")
	if err == nil {
		fmt.Println("Found value in cache:", res.Data().(*myStruct).text)
	} else {
		fmt.Println("Error retrieving value from cache:", err)
	}

	time.Sleep(6 * time.Second)
	res, err = cache.Value("someKey")
	if err != nil {
		fmt.Println("Item is not cached (anymore).")
	}

	cache.Add("someKey", 0, &val)

	cache.SetAboutToDeleteItemCallback(func(e *cache2go.CacheItem) {
		fmt.Println("Deleting:", e.Key(), e.Data().(*myStruct).text, e.CreatedOn())
	})

	cache.Delete("someKey")
	cache.Flush()
}
