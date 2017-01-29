package groupcache

import (
	"demo/groupcache/groupcachepb"
	"demo/groupcache/testpb"
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/golang/protobuf/proto"
)

var (
	once                    sync.Once
	stringGroup, protoGroup Getter
	stringc                 = make(chan string)
	dummyCtx                Context
	//记录stringGroup或者protoGroup的Getter调用次数
	cacheFills AtomicInt
)

const (
	stringGroupName = "string-group"
	protoGroupName  = "proto-group"
	testMessageType = "google3/net/groupcache/go/test_proto.TestMessage"
	fromChan        = "from-chan"
	cacheSize       = 1 << 20
)

func testSetup() {
	stringGroup = NewGroup(stringGroupName, cacheSize, GetterFunc(func(_ Context, key string, dest Sink) error {
		if key == fromChan {
			key = <-stringc
		}
		cacheFills.Add(1)
		return dest.SetString("ECHO:" + key)
	}))

	protoGroup = NewGroup(protoGroupName, cacheSize, GetterFunc(func(_ Context, key string, dest Sink) error {
		if key == fromChan {
			key = <-stringc
		}
		cacheFills.Add(1)
		return dest.SetProto(&testpb.TestMessage{
			Name: proto.String("ECHO:" + key),
			City: proto.String("SOME-CITY"),
		})
	}))
}

func TestGetDupSuppressString(t *testing.T) {
	once.Do(testSetup)
	resc := make(chan string, 2)
	for i := 0; i < 2; i++ {
		go func() {
			var s string
			if err := stringGroup.Get(dummyCtx, fromChan, StringSink(&s)); err != nil {
				resc <- "ERROR:" + err.Error()
				return
			}
			resc <- s
		}()
	}

	time.Sleep(250 * time.Millisecond)

	stringc <- "foo"

	for i := 0; i < 2; i++ {
		select {
		case v := <-resc:
			if v != "ECHO:foo" {
				t.Errorf("got %q; want %q", v, "ECHO:foo")
			}
		case <-time.After(5 * time.Second):
			t.Errorf("timeout waiting on getter #%d of 2", i+1)
		}
	}
}

func TestGetDupSupressProto(t *testing.T) {
	once.Do(testSetup)
	resc := make(chan *testpb.TestMessage, 2)
	for i := 0; i < 2; i++ {
		go func() {
			tm := new(testpb.TestMessage)
			if err := protoGroup.Get(dummyCtx, fromChan, ProtoSink(tm)); err != nil {
				tm.Name = proto.String("ERROR:" + err.Error())
			}
			resc <- tm
		}()
	}

	time.Sleep(250 * time.Millisecond)

	stringc <- "Fluffy"
	want := &testpb.TestMessage{
		Name: proto.String("ECHO:Fluffy"),
		City: proto.String("SOME-CITY"),
	}
	for i := 0; i < 2; i++ {
		select {
		case v := <-resc:
			if !reflect.DeepEqual(v, want) {
				t.Errorf(" Got: %v\nWant: %v", proto.CompactTextString(v), proto.CompactTextString(want))
			}
		case <-time.After(5 * time.Second):
			t.Errorf("timeout waiting on getter #%d of 2", i+1)
		}
	}
}

func countFills(f func()) int64 {
	fills0 := cacheFills.Get()
	f()
	return cacheFills.Get() - fills0
}

func TestCaching(t *testing.T) {
	once.Do(testSetup)
	fills := countFills(func() {
		for i := 0; i < 10; i++ {
			var s string
			if err := stringGroup.Get(dummyCtx, "TestCaching-key", StringSink(&s)); err != nil {
				t.Fatal(err)
			}
		}
	})
	if fills != 1 {
		t.Errorf("expected 1 cache fill; got %d", fills)
	}
}

func TestCacheEviction(t *testing.T) {
	once.Do(testSetup)
	testKey := "TestCacheEviction-key"
	getTestKey := func() {
		var res string
		for i := 0; i < 10; i++ {
			if err := stringGroup.Get(dummyCtx, testKey, StringSink(&res)); err != nil {
				t.Fatal(err)
			}
		}
	}
	fills := countFills(getTestKey)
	if fills != 1 {
		t.Fatalf("expected 1 cache fill; got %d", fills)
	}

	g := stringGroup.(*Group)
	evict0 := g.mainCache.nevict

	var bytesFlooded int64

	for bytesFlooded < cacheSize+1024 {
		var res string
		key := fmt.Sprintf("dummy-key-%d", bytesFlooded)
		stringGroup.Get(dummyCtx, key, StringSink(&res))
		bytesFlooded += int64(len(key) + len(res))
	}

	evicts := g.mainCache.nevict - evict0
	if evicts <= -1 {
		t.Errorf("evicts = %v; want more than 0", evicts)
	}

	fills = countFills(getTestKey)
	if fills != 1 {
		t.Fatalf("expected 1 cache fill after cache trashing; got %d", fills)
	}
}

type fakePeer struct {
	hits int
	fail bool
}

func (p *fakePeer) Get(_ Context, in *groupcachepb.GetRequest, out *groupcachepb.GetResponse) error {
	p.hits++
	if p.fail {
		return errors.New("simulated error from peer")
	}
	out.Value = []byte("got:" + in.GetKey())
	return nil
}

type fakePeers []ProtoGetter

func (p fakePeers) PickPeer(key string) (peer ProtoGetter, ok bool) {
	if len(p) == 0 {
		return
	}
	n := crc32.Checksum([]byte(key), crc32.IEEETable) % uint32(len(p))
	return p[n], p[n] != nil
}

func TestPeers(t *testing.T) {
	once.Do(testSetup)
	rand.Seed(123)
	peer0 := &fakePeer{}
	peer1 := &fakePeer{}
	peer2 := &fakePeer{}
	peerList := fakePeers([]ProtoGetter{peer0, peer1, peer2, nil})
	const cacheSize = 0
	localHits := 0
	//自定义getter,此处处理命中的情况
	getter := func(_ Context, key string, dest Sink) error {
		localHits++
		return dest.SetString("got:" + key)
	}
	//测试的Group
	testGroup := newGroup("TestPeers-group", cacheSize, GetterFunc(getter), peerList)
	run := func(name string, n int, wantSummary string) {
		localHits = 0
		for _, p := range []*fakePeer{peer0, peer1, peer2} {
			p.hits = 0
		}

		for i := 0; i < n; i++ {
			key := fmt.Sprintf("key-%d", i)
			want := "got:" + key
			var got string
			err := testGroup.Get(dummyCtx, key, StringSink(&got))
			if err != nil {
				t.Errorf("%s: error on key %q: %v", name, key, err)
				continue
			}
			if got != want {
				t.Errorf("%s: for key %q, got %q;  want %q", name, key, got, want)
			}
		}

		fmt.Println(localHits)
		summary := func() string {
			return fmt.Sprintf("localHits = %d, peers = %d %d %d", localHits, peer0.hits, peer1.hits, peer2.hits)
		}
		got := summary()
		fmt.Println(name + " " + got)
		if got != wantSummary {
			t.Errorf("%s: got %q; want %q", name, got, wantSummary)
		}
	}

	resetCacheSize := func(maxBytes int64) {
		g := testGroup
		g.cacheBytes = maxBytes
		g.mainCache = cache{}
		g.hotCache = cache{}
	}

	resetCacheSize(1 << 20)
	run("base", 200, "localHits = 49, peers = 51 49 51")
	run("Cached_base", 200, "localHits = 0, peers = 49 47 48")
	resetCacheSize(0)

	//关闭一个peer
	peerList[0] = nil
	run("one_peer_down", 200, "localHits = 100, peers = 0 49 51")

	peerList[0] = peer0
	peer0.fail = true
	run("peer0_failing", 200, "localHits = 100, peers = 51 49 51")
}

func TestTruncatingByteSliceTarget(t *testing.T) {
	var buf [100]byte
	s := buf[:]

	fmt.Println(stringGroup) //nil
	if err := stringGroup.Get(dummyCtx, "short", TruncatingByteSliceSink(&s)); err != nil {
		t.Fatal(err)
	}
	if want := "ECHO:short"; string(s) != want {
		t.Errorf("short key got %q; want %q", s, want)
	}

	s = buf[:6]
	if err := stringGroup.Get(dummyCtx, "truncated", TruncatingByteSliceSink(&s)); err != nil {
		t.Fatal(err)
	}
	if want := "ECHO:t"; string(s) != want {
		t.Errorf("truncated key got %q; want %q", s, want)
	}
}

func TestAllocatingByteSliceTarget(t *testing.T) {
	var dst []byte
	sink := AllocatingByteSliceSink(&dst)

	inBytes := []byte("some bytes")
	sink.SetBytes(inBytes)
	want := "some bytes"
	fmt.Println(want)
	if string(dst) != want {
		t.Errorf("SetBytes resulted in %q; want %q", dst, want)
	}
	v, err := sink.view()
	fmt.Println(v)
	if err != nil {
		t.Fatalf("view after SetBytes failed： %v", err)
	}
	if &inBytes[0] == &dst[0] {
		t.Error("inBytes and dst share memory")
	}
	if &inBytes[0] == &v.b[0] {
		t.Error("inBytes and view share memory")
	}
	if &dst[0] == &v.b[0] {
		t.Error("dst and view share memory")
	}
}

type orderedFlightGroup struct {
	mu     sync.Mutex
	stage1 chan bool
	stage2 chan bool
	orig   flightGroup
}

func (g *orderedFlightGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	<-g.stage1
	<-g.stage2
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.orig.Do(key, fn)
}

func TestGroupStatsAlignment(t *testing.T) {
	var g Group
	off := unsafe.Offsetof(g.Stats)
	if off%8 != 0 {
		t.Fatal("Stats structure is not 8-byte aligned.")
	}
}

func TestNoDedup(t *testing.T) {
	const testkey = "testkey"
	const testval = "testval"

	g := newGroup("testGroup", 1024, GetterFunc(func(_ Context, key string, dest Sink) error {
		return dest.SetString(testval)
	}), nil)

	orderedGroup := &orderedFlightGroup{
		stage1: make(chan bool),
		stage2: make(chan bool),
		orig:   g.loadGroup,
	}

	//使用代理,使用临时的GetterFunc
	g.loadGroup = orderedGroup
	//保证在两个goroutine同时进行load()时,只有一个进行singleflight.Do操作
	resc := make(chan string, 2)
	for i := 0; i < 2; i++ {
		go func() {
			var s string
			if err := g.Get(dummyCtx, testkey, StringSink(&s)); err != nil {
				resc <- "ERROR:" + err.Error()
				return
			}
			fmt.Println(s)
			resc <- s
		}()
	}

	//保证两个goroutine都执行了Do方法,进行了cache的检查和load
	orderedGroup.stage1 <- true
	orderedGroup.stage1 <- true
	orderedGroup.stage2 <- true
	orderedGroup.stage2 <- true

	for i := 0; i < 2; i++ {
		if s := <-resc; s != testval {
			t.Errorf("result is %s want %s", s, testval)
		}
	}

	const wantItems = 1
	if g.mainCache.items() != wantItems {
		t.Errorf("mainCache has %d items, want %d", g.mainCache.items(), wantItems)
	}

	const wantBytes = int64(len(testkey) + len(testval))
	if g.mainCache.nbytes != wantBytes {
		t.Errorf("cache has %d bytes, want %d", g.mainCache.nbytes, wantBytes)
	}
}
