// 功能：使用通道实现 future/promise
// 实现：读通道作为函数返回值
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func longTimeRequest() <- chan int32 {
	r := make(chan int32)
	go func() {
		time.Sleep(time.Second * 3)
		r <- rand.Int31n(100)
	}()
	return r
}

func sumSquares(a, b int32) int32 {
	return a*a + b*b
}

func main() {
	a, b := longTimeRequest(), longTimeRequest()
	fmt.Println(sumSquares(<-a, <-b))
}
