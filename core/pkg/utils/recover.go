package utils

import (
	"fmt"
	"runtime/debug"

	"github.com/kunlun-sec/lunying/pkg/logger"
)

// RecoverPanic 捕获并处理panic
func RecoverPanic() {
	if r := recover(); r != nil {
		logger.Error("发生panic: %v\n堆栈信息:\n%s", r, string(debug.Stack()))
	}
}

// SafeGo 安全地启动goroutine，自动捕获panic
func SafeGo(fn func()) {
	go func() {
		defer RecoverPanic()
		fn()
	}()
}

// TryCatch 类似try-catch的错误处理
func TryCatch(try func() error, catch func(error)) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v\n%s", r, string(debug.Stack()))
			logger.Error("TryCatch捕获panic: %v", err)
			if catch != nil {
				catch(err)
			}
		}
	}()

	if err := try(); err != nil {
		if catch != nil {
			catch(err)
		}
	}
}

// Result 包装可能panic的函数返回值
type Result struct {
	Value interface{}
	Error error
}

// SafeCall 安全调用可能panic的函数
func SafeCall(fn func() interface{}) (result Result) {
	defer func() {
		if r := recover(); r != nil {
			result.Error = fmt.Errorf("panic: %v\n%s", r, string(debug.Stack()))
			logger.Error("SafeCall捕获panic: %v", result.Error)
		}
	}()

	result.Value = fn()
	return
}
