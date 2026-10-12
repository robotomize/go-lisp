package main

import (
	"bufio"
	"fmt"
	"os"
	"robotomize/go-lisp/internal/lexer"
)

func init() {

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("go-lisp>>")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		l := lexer.New(line)
		for {
			token, ok := l.Next()
			if !ok {
				break
			}

			fmt.Println(token)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
