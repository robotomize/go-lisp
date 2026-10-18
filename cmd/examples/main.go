package main

import (
	"fmt"
	"robotomize/go-lisp/internal/types"
)

func main() {
	list := types.List(types.NewSymbol("ok"), types.NewInt(10), types.NewInt(20), types.NewSymbol("ok"), types.NewCons(types.NewSymbol("ok"), types.NewInt(30)))
	fmt.Println(list.(*types.Cons).Car)
	fmt.Println(list)
}
