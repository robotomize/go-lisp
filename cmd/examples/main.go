package main

import (
	"fmt"
	"time"
)

func fib(n int) int {
	if n <= 1 {
		return n
	}

	return fib(n-1) + fib(n-2)
}

func main() {
	t := time.Now()
	fib(40)
	fmt.Println(time.Since(t))
}
