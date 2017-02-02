package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(getAbsFilePath(dir))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil
	}
	return entries
}

func walkDir(dir string, n *sync.WaitGroup, filelength chan<- int64) {
	defer n.Done()
	if cancelled() {
		fmt.Println("cancelled called---")
		return
	}

	for _, entry := range dirents(dir) {
		if strings.Contains(entry.Name(), ".git") {
			continue
		}
		if strings.Contains(entry.Name(), ".idea") {
			continue
		}
		if entry.IsDir() && entry.Name() != ".git" {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, filelength)
		} else if compareExt(entry.Name(), *ext) {
			absPath := getAbsFilePath(dir + "/" + entry.Name())
			lines := getLinesOfFile(absPath)
			fmt.Printf("%s %d\n", absPath, lines)
			filelength <- lines
		}
	}
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

var verbose = flag.Bool("v", false, "show verbose progress message")
var ext = flag.String("ext", "go", "filename extension")
var done = make(chan struct{})
var sema = make(chan struct{}, 20)

func main() {
	flag.Parse()
	roots := flag.Args()

	if 0 == len(roots) {
		roots = []string{"."}
	}

	filelength := make(chan int64)
	var tick <-chan time.Time

	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}

	var n sync.WaitGroup

	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, filelength)
	}

	go func() {
		n.Wait()
		close(filelength)
	}()

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	var nfiles, nbytes int64

loop:
	for {
		select {
		case <-done:
			for range filelength {

			}
			return
		case size, ok := <-filelength:
			if !ok {
				break loop
			}
			nfiles++
			nbytes += size
		case <-tick:
			fmt.Printf("%d %d\n", nfiles, nbytes)
		}
	}
	fmt.Printf("%d %d\n", nfiles, nbytes)
}

func getLinesOfFile(filename string) int64 {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	buf := bufio.NewReader(f)
	var lines int64 = 0
	for {
		_, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		lines++
	}
	return lines
}

func getAbsFilePath(path string) string {
	if path == "" {
		return ""
	}
	//	fmt.Println(path)
	path = filepath.FromSlash(path)
	if !filepath.IsAbs(path) {
		baseDir, err := os.Getwd()
		if err != nil {
			panic(errors.New(fmt.Sprintf("Can not get current work dir : %s\n", err)))
		}
		path = filepath.Join(baseDir, path)
	}
	return path
}

func compareExt(filename, ext string) bool {
	if path.Ext(filename) == "."+ext {
		return true
	}
	return false
}
