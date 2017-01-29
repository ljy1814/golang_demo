package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	//	fmt.Println("vim-go")
	//	getAllFileLines(os.Args[1])
	getFileList(os.Args[1])
	fmt.Printf("all files line : %d\n", totalLines)
}

func getAllFileLines(path string) uint {
	fd, err := os.Open(path)
	check(err)
	finfo, err := fd.Stat()
	if finfo.IsDir() && finfo.Name() != ".git" {
		fmt.Println(finfo.Name())
	}
	//	fmt.Println(finfo)
	//	display.Display("finfo", finfo)
	return 0
}

func getLinesOfFile(filename string) uint {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	buf := bufio.NewReader(f)
	var lines uint = 0
	for {
		_, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		lines++
	}
	return lines
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var totalLines uint = 0

func getFileList(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() || f.Name() == ".git" {
			return nil
		}
		totalLines += getLinesOfFile(path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
