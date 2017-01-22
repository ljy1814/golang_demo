package base

import (
	"demo/mylog"
)

//创建日志记录器
func NewLogger() mylog.Logger {
	return mylog.NewSimpleLogger()
}
