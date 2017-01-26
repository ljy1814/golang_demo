package cache2go

import "errors"

var (
	ErrKeyNotFound         = errors.New("Key not found in cache")
	ErrKeyNotFoundOrLoaded = errors.New("Key not found Or could not be loaded into cache")
)
