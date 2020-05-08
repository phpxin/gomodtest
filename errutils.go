/*
 @Title  错误以及崩溃处理
 @Description
 @Author  Leo  2020/4/20 5:05 下午
 @Update  Leo  2020/4/20 5:05 下午
*/

package main

import (
	"fmt"
	"runtime"
)

func PrintPanicStackError() {
	if x := recover(); x != nil {
		fmt.Println("panic ",x)
		PrintPanicStack()
	}
}

func PrintPanicStack() {
	for i := 0; i < 10; i++ {
		funcName, file, line, ok := runtime.Caller(i)
		if ok {
			funcName := runtime.FuncForPC(funcName).Name()
			errInfo := fmt.Sprintf("frame %d:[func:%s, file: %s, line:%d]", i, funcName, file, line)
			fmt.Println(errInfo)
		}
	}
}
