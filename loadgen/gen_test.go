package loadgen

import  (
	"demo/loadgen/lib"
	"demo/loadgen/testhelper"
	"runtime"
	"testing"
	"time"
//	"demo/display"
//	"fmt"
)

var printDetail = false
func TestStart(t *testing.T) {
	//设置P最大数量
	runtime.GOMAXPROCS(runtime.NumCPU())
	//初始化服务器
	server := testhelper.NewTcpServer()
	defer server.Close()
//	display.Display("server", server)
	serverAddr := "127.0.0.1:10000"
	t.Logf("Startup TCP server(%s)...\n", serverAddr)
	if err := server.Listen(serverAddr); err != nil {
		t.Fatalf("TCP Server startup failing! (addr=%s)!\n", serverAddr)
		t.FailNow()
	}
//	fmt.Println("OK....")
//	display.Display("XXXXXXXXXX", server)
//	time.Sleep(2 * time.Second)
	//初始化调度器
	comm := testhelper.NewTcpComm(serverAddr)
	//初始化载荷发生器
	resultCh := make(chan *lib.CallResult, 50)
	timeoutNs := 3 * time.Millisecond
	lps := uint32(200)
	durationNs := 12 * time.Second
	t.Logf("Initialize load generator (timeoutNs=%v, lps=%d, durationNs=%v)...", timeoutNs, lps, durationNs)
	gen, err := NewGenerator( comm, timeoutNs, lps, durationNs, resultCh)
	if err != nil {
		t.Fatalf("Load generator initialization failing: %s.\n", err)
		t.FailNow()
	}
	//开始
	t.Log("Start load generator...")
	gen.Start()
	//显示结果
	countMap := make(map[lib.ResultCode]int)
	for  r := range resultCh {
//		display.Display("r", r)
		countMap[r.Code] = countMap[r.Code] + 1
		if printDetail {
			t.Logf("Result: Id=%d, Code=%d, Msg=%d, Elapse=%v.\n", r.Id, r.Code, r.Msg, r.Elapse)
		}
	}
//	fmt.Println(countMap)
//	display.Display("countMap", countMap)
	var total int
	t.Log("Code Count:")
	for k, v := range countMap {
		codePlain := lib.GetResultCodePlain(k)
		t.Logf(" Code Plain: %s (%d), Count: %d.\n", codePlain, k, v)
		total += v
	}
	t.Logf("Total load: %d.\n", total)
	successCount := countMap[lib.RESULT_CODE_SUCCESS]
	tps := float64(successCount) / float64(durationNs / 1e9)
	t.Logf("Loads per second: %d; Treatments per second: %f.\n", lps, tps)
}

func TestStop(t *testing.T) {
	//设置P最大数量
	runtime.GOMAXPROCS(runtime.NumCPU())
	//初始化服务器
	server := testhelper.NewTcpServer()
	defer server.Close()
	serverAddr := "127.0.0.1:11000"
	t.Logf("Startup TCP server(%s)...\n", serverAddr)
	if err := server.Listen(serverAddr); err != nil {
		t.Fatalf("TCP Server startup failing! (addr=%s)!\n", serverAddr)
		t.FailNow()
	}
	//初始化调度器
	comm := testhelper.NewTcpComm(serverAddr)
	//初始化载荷发生器
	resultCh := make(chan *lib.CallResult, 50)
	timeoutNs := 3 * time.Millisecond
	lps := uint32(200)
	durationNs := 10 * time.Second
	t.Logf("Initialize load generator (timeoutNs=%v, lps=%d, durationNs=%v)...", timeoutNs, lps, durationNs)
	gen, err := NewGenerator( comm, timeoutNs, lps, durationNs, resultCh)
	if err != nil {
		t.Fatalf("Load generator initialization failing: %s.\n", err)
		t.FailNow()
	}
	//开始
	t.Log("Start load generator...")
	gen.Start()
	//显示结果
	countMap := make(map[lib.ResultCode]int)
	count := 0
	for  r := range resultCh {
		countMap[r.Code] = countMap[r.Code] + 1
		if printDetail {
			t.Logf("Result: Id=%d, Code=%d, Msg=%d, Elapse=%v.\n", r.Id, r.Code, r.Msg, r.Elapse)
		}
		count++
		if count > 3 {
			gen.Stop()
		}
	}
	var total int
	t.Log("Code Count:")
	for k, v := range countMap {
		codePlain := lib.GetResultCodePlain(k)
		t.Logf(" Code Plain: %s (%d), Count: %d.\n", codePlain, k, v)
		total += v
	}
	t.Logf("Total load: %d.\n", total)
	successCount := countMap[lib.RESULT_CODE_SUCCESS]
	tps := float64(successCount) / float64(durationNs / 1e9)
	t.Logf("Loads per second: %d; Treatments per second: %f.\n", lps, tps)
}
