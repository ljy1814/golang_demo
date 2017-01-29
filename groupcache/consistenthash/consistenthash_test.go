package consistenthash

import (
	"demo/display"
	"fmt"
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, err := strconv.Atoi(string(key))
		if err != nil {
			panic(err)
		}
		return uint32(i)
	})

	hash.Add("6", "4", "2")
	testCases := map[string]string{
		"2":  "2", // 2 >= 2
		"11": "2", //12 > 11
		"23": "4", //24 > 23 故为4
		"27": "2", //27 > 26 故还原为0
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
	/*  这是使用Dpsplay得出的结果
	display.Display("hash", hash)
		Display hash (*consistenthash.Map):
		(*hash).hash = consistenthash.Hash0x46e720
		(*hash).replicas = 3  //迭代3层;2,4,6增序排列
		(*hash).keys[0] = 2
		(*hash).keys[1] = 4
		(*hash).keys[2] = 6
		(*hash).keys[3] = 12
		(*hash).keys[4] = 14
		(*hash).keys[5] = 16
		(*hash).keys[6] = 22
		(*hash).keys[7] = 24
		(*hash).keys[8] = 26
		(*hash).hashMap[6] = "6"
		(*hash).hashMap[2] = "2"
		(*hash).hashMap[12] = "2"
		(*hash).hashMap[22] = "2"
		(*hash).hashMap[16] = "6"
		(*hash).hashMap[26] = "6"
		(*hash).hashMap[4] = "4"
		(*hash).hashMap[14] = "4"
		(*hash).hashMap[24] = "4"
	*/

	hash.Add("8")
	testCases["27"] = "8"
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
	//	display.Display("hash", hash)
}

func TestConsistency(t *testing.T) {
	hash1 := New(1, nil)
	hash2 := New(1, nil)

	hash1.Add("Bill", "Bob", "Bonny")
	hash2.Add("Bob", "Bonny", "Bill")

	if hash1.Get("Ben") != hash2.Get("Ben") {
		t.Errorf("Fetching 'ben' from both hashes shoule be the same")
	}
	/*
		display.Display("hash1", hash1)
		display.Display("hash2", hash2)
		Display hash1 (*consistenthash.Map):
		(*hash1).hash = consistenthash.Hash0x4bc1c0
		(*hash1).replicas = 1
		(*hash1).keys[0] = 1679827945
		(*hash1).keys[1] = 2622760538
		(*hash1).keys[2] = 3819440399
		(*hash1).hashMap[2622760538] = "Bill"
		(*hash1).hashMap[3819440399] = "Bob"
		(*hash1).hashMap[1679827945] = "Bonny"
		Display hash2 (*consistenthash.Map):
		(*hash2).hash = consistenthash.Hash0x4bc1c0
		(*hash2).replicas = 1
		(*hash2).keys[0] = 1679827945
		(*hash2).keys[1] = 2622760538
		(*hash2).keys[2] = 3819440399
		(*hash2).hashMap[1679827945] = "Bonny"
		(*hash2).hashMap[2622760538] = "Bill"
		(*hash2).hashMap[3819440399] = "Bob"
	*/

	hash2.Add("Becky", "Ben", "Bobby")
	display.Display("hash1", hash1)
	display.Display("hash2", hash2)
	if hash1.Get("Ben") != hash2.Get("Ben") ||
		hash1.Get("Bob") != hash2.Get("Bob") ||
		hash1.Get("Bonny") != hash2.Get("Bonny") {
		t.Errorf("Direct matches should always return the same entry")
	}
	//由于Ben的哈希值最小,而最小的hash值为键对于的值为Bonny
	fmt.Println(hash1.Get("Ben")) //Bonny
	fmt.Println(hash2.Get("Ben"))
	fmt.Println(hash1.Get("Bob")) //Bob 对应自己
	fmt.Println(hash2.Get("Bob"))
	fmt.Println(hash1.Get("Bonny")) //Bonny  Bonny刚好对应自己
	fmt.Println(hash2.Get("Bonny"))
	/*
		Display hash1 (*consistenthash.Map):
		(*hash1).hash = consistenthash.Hash0x4bc880
		(*hash1).replicas = 1
		(*hash1).keys[0] = 1679827945
		(*hash1).keys[1] = 2622760538
		(*hash1).keys[2] = 3819440399
		(*hash1).hashMap[2622760538] = "Bill"
		(*hash1).hashMap[3819440399] = "Bob"
		(*hash1).hashMap[1679827945] = "Bonny"
		Display hash2 (*consistenthash.Map):
		(*hash2).hash = consistenthash.Hash0x4bc880
		(*hash2).replicas = 1
		(*hash2).keys[0] = 284274094
		(*hash2).keys[1] = 1679827945
		(*hash2).keys[2] = 2117248155
		(*hash2).keys[3] = 2622760538
		(*hash2).keys[4] = 3247412609
		(*hash2).keys[5] = 3819440399
		(*hash2).hashMap[3819440399] = "Bob"
		(*hash2).hashMap[1679827945] = "Bonny"
		(*hash2).hashMap[2622760538] = "Bill"
		(*hash2).hashMap[2117248155] = "Becky"
		(*hash2).hashMap[284274094] = "Ben"
		(*hash2).hashMap[3247412609] = "Bobby"
		Bonny
		Bonny
		Bob
		Bob
		Bonny
		Bonny
	*/
}

func BenchmarkGet8(b *testing.B) {
	benchmarkGet(b, 8)
}
func BenchmarkGet32(b *testing.B) {
	benchmarkGet(b, 32)
}
func BenchmarkGet128(b *testing.B) {
	benchmarkGet(b, 128)
}
func BenchmarkGet512(b *testing.B) {
	benchmarkGet(b, 512)
}

func benchmarkGet(b *testing.B, shards int) {
	hash := New(50, nil)
	var buckets []string
	for i := 0; i < shards; i++ {
		buckets = append(buckets, fmt.Sprintf("shard-%d", i))
	}

	hash.Add(buckets...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash.Get(buckets[i&(shards-1)])
	}
}
