package groupcache

import (
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	peerAddrs = flag.String("test_peer_addrs", "", "Comma-separated list of peer addresses; used by TestHTTPPool")
	peerIndex = flag.Int("test_peer_index", -1, "Index of which peer this child is; used by TestHTTPPool")
	peerChild = flag.Bool("test_peer_child", false, "True if running as a child process; used by TestHTTPPool")
)

func TestHTTPPool(t *testing.T) {
	if *peerChild {
		beChildForTestHTTPPool()
		os.Exit(0)
	}

	const (
		nChild = 4
		nGets  = 100
	)

	var childAddr []string
	for i := 0; i < nChild; i++ {
		childAddr = append(childAddr, pickFreeAddr(t))
		//		fmt.Println(childAddr)
	}

	var cmds []*exec.Cmd
	var wg sync.WaitGroup

	for i := 0; i < nChild; i++ {
		cmd := exec.Command(os.Args[0],
			"--test.run=TestHTTPPool",
			"--test_peer_child",
			"--test_peer_addrs="+strings.Join(childAddr, ","),
			"--test_peer_index="+strconv.Itoa(i),
		)
		cmds = append(cmds, cmd)
		//		fmt.Println(cmds)
		wg.Add(1)
		if err := cmd.Start(); err != nil {
			t.Fatal("failed to start child process: ", err)
		}
		go awaitAddrReady(t, childAddr[i], &wg)
	}
	defer func() {
		for i := 0; i < nChild; i++ {
			if cmds[i].Process != nil {
				cmds[i].Process.Kill()
			}
		}
	}()
	wg.Wait()

	//使用无效URL,保证不处理get请求
	p := NewHTTPPool("should-be-ignored")
	p.Set(addrToURL(childAddr)...)

	getter := GetterFunc(func(ctx Context, key string, dest Sink) error {
		return errors.New("parent getter called; someething's wrong")
	})
	//本机测试环境
	g := NewGroup("httpPoolTest", 1<<20, getter)

	for _, key := range testKeys(nGets) {
		var value string
		if err := g.Get(nil, key, StringSink(&value)); err != nil {
			t.Fatal(err)
		}
		// value	3:127.0.0.1:44549:99
		if suffix := ":" + key; !strings.HasSuffix(value, suffix) {
			t.Errorf("Get(%q), want value ending in %q", key, value, suffix)
		}
		t.Logf("Get key = %q, value=%q (peer:key)", key, value)
	}
}

func testKeys(n int) (keys []string) {
	keys = make([]string, n)
	for i := range keys {
		keys[i] = strconv.Itoa(i)
	}
	return
}

//网络上的
func beChildForTestHTTPPool() {
	addrs := strings.Split(*peerAddrs, ",")
	//连接的端口是80
	//	*peerIndex = 0
	p := NewHTTPPool("http://" + addrs[*peerIndex])
	p.Set(addrToURL(addrs)...)
	//	fmt.Println(p)

	getter := GetterFunc(func(ctx Context, key string, dest Sink) error {
		//value从此处返回
		dest.SetString(strconv.Itoa(*peerIndex) + ":" + addrs[*peerIndex] + ":" + key)
		return nil
	})
	NewGroup("httpPoolTest", 1<<20, getter)
	log.Fatal(http.ListenAndServe(addrs[*peerIndex], p))
}

func pickFreeAddr(t *testing.T) string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().String()
}

func addrToURL(addr []string) []string {
	url := make([]string, len(addr))
	for i := range addr {
		url[i] = "http://" + addr[i]
	}
	return url
}

func awaitAddrReady(t *testing.T, addr string, wg *sync.WaitGroup) {
	defer wg.Done()
	const max = 1 * time.Second
	tries := 0
	for {
		tries++
		c, err := net.Dial("tcp", addr)
		//		fmt.Println(c)
		//		fmt.Println(err)
		if err == nil {
			c.Close()
			return
		}
		delay := time.Duration(tries) * 25 * time.Millisecond
		if delay > max {
			delay = max
		}
		time.Sleep(delay)
	}
}
