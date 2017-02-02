package api

//一个 "Load" rpc 指令
type Load struct {
	Key string
}

//a "store" rpc instruction
type Store struct {
	Key   string
	Value string
}

type NullResult int

type ValueResult struct {
	Value string
}
