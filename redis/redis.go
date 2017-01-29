package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func main() {
	fmt.Println("vim-go")
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	fmt.Println(c)

	v, err := c.Do("SET", "name", "red")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)

	//	v, err = c.Do("lpush", "redlist", "qqq")
	//	v, err = c.Do("lpush", "redlist", "www")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)
	//	v, err = c.Do("lpush", "redlist", "xxx")
	if err != nil {
		fmt.Println(err)
		return
	}
	//v为redis库存入的列表数据条数
	fmt.Println(v)

	values, _ := redis.Values(c.Do("lrange", "redlist", "0", "100"))
	for _, v := range values {
		fmt.Println(v)
		fmt.Println(string(v.([]byte)))
	}

	fmt.Println("-------------------")
	var v1 string
	redis.Scan(values, &v1)
	fmt.Println(v1)
	fmt.Println("---------pipeline----------")
	c.Send("SET", "value", "green")
	c.Send("GET", "value")
	c.Flush()
	c.Receive()
	c.Receive()
	fmt.Println("-------------------")
	go subscribe()
	go subscribe()
	go subscribe()
	go subscribe()
	go subscribe()

	for {
		var s string
		fmt.Scanln(&s)
		_, err := c.Do("PUBLISH", "redChatRoom", s)
		if err != nil {
			fmt.Println("pub err: ", err)
			return
		}
	}
}

func subscribe() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	psc := redis.PubSubConn{c}
	psc.Subscribe("redChatRoom")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s : message :  %s\n", v.Channel, v.Data)
		case redis.Subscription:
			fmt.Printf("%s : %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			fmt.Println(v)
			return
		}
	}
}
