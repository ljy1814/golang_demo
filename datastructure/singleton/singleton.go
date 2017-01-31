package singleton

import "sync"

type Single struct {
	data int
}

var (
	s    *Single     //指针才可以与nil比较
	lock *sync.Mutex = &sync.Mutex{}
	once sync.Once
)

func GetInstance() interface{} {
	if s == nil {
		s = &Single{data: 0}
	}
	return s
}
func GetInstance2() interface{} {
	if s == nil {
		lock.Lock()
		defer lock.Unlock()
		s = &Single{data: 2}
	}
	return s
}
func GetInstance3() interface{} {
	once.Do(func() {
		s = &Single{data: 3}
	})
	return s
}
