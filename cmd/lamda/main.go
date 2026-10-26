package main

import "fmt"

func main() {
	type selfFunc func(selfFunc, int) int

	result := func(f selfFunc, n int) int {
		return f(f, n)
	}(
		func(self selfFunc, currentN int) int {
			if currentN == 0 {
				return 1
			}
			return currentN * self(self, currentN-1)
		},
		5,
	)
	fmt.Println(result)

}
