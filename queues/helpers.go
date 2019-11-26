package queues

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

// TODO: REMOVE THIS FILE COMPLETELY
// true if error exists
func ErrorHandler(err error) (b bool) {
	b = false
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, fn, line, _ := runtime.Caller(1)

		logrus.Error(fmt.Sprintf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err))
		logrus.Error(fmt.Sprintf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err))
		b = true
	}
	return b
}

func ErrorLogger(err error) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, fn, line, _ := runtime.Caller(1)

		fmt.Printf(fmt.Sprintf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err))
		logrus.Error(fmt.Sprintf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err))
	}
}
