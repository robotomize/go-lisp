package main

import (
	"fmt"
	"robotomize/go-lisp/internal/types"
)

func main() {
	list := types.List(types.NewSymbol("ok"), types.NewNumber(10), types.NewNumber(20), types.NewSymbol("ok"), types.NewCons(types.NewSymbol("ok"), types.NewNumber(30)))
	fmt.Println(list.(*types.Cons).Car)
	fmt.Println(list)
}
