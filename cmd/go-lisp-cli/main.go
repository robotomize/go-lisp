package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"robotomize/go-lisp/internal/lexer"
)

func init() {

}

func main() {
	argv := os.Args[1:]

	var input io.Reader
	isInteractive := true

	if len(argv) > 0 {
		filename := argv[0]
		if _, err := os.Stat(filename); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		defer f.Close()
		input = f
		isInteractive = false
	}

	scanner := bufio.NewScanner(input)

	for {
		if isInteractive {
			fmt.Print("go-lisp>>")
		}

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
