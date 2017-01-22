package loadgen

import (
	"demo/mylog"
)

var logger mylog.Logger

func init() {
	logger = mylog.NewSimpleLogger()
}
