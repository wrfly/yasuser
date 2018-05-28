package utils

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

func AddLineNum(v ...interface{}) string {
	var msg string
	// notice that we're using 1, so it will actually log where
	// the error happened, 0 = this function, we don't want that.
	_, fn, line, _ := runtime.Caller(1)
	if len(v) == 0 {
		msg = fmt.Sprintf("%s", v)
	} else {
		defer func() {
			x := recover()
			if x != nil {
				logrus.Errorf("PANIC: %s:%d %s", fn, line, x)
			}
		}()
		switch v[0].(type) {
		case string:
			msg = fmt.Sprintf(v[0].(string), v[1:])
		default:
			msg = fmt.Sprintf("%v", v)
		}
	}
	return fmt.Sprintf("%s:%d %s", fn, line, msg)
}

func LineNum() string {
	_, fn, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", fn, line)
}
