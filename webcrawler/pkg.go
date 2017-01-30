package main

import (
    "fmt"
//    "os"
	"strings"
//    do "gopkg.in/godo.v2"
)

/*
func tasks(p *do.Project) {
    if pwd, err := os.Getwd(); err == nil {
        do.Env = fmt.Sprintf("GOPATH=%s/vendor::$GOPATH", pwd)
    }

    p.Task("server", nil, func(c *do.Context) {
        c.Start("main.go ./config/page.yaml", do.M{"$in": "./"})
    }).Src("**///*.go")
//}


func main() {
//    do.Godo(tasks)
	strs()
}

func strs() {
	http1 := "http"
	https1 := "https"
	if strings.ToLower(http1) == "http" {
		fmt.Println(http1)
	}
	if strings.ToLower(https1) == "https" {
		fmt.Println(https1)
	}
	if !(strings.ToLower(http1) == "http" || strings.ToLower(http1) == "https") {
		fmt.Printf("%s in https or http\n", http1)
	}
}
