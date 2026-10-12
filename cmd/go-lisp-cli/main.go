package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"robotomize/go-lisp/internal/lexer"
	"robotomize/go-lisp/internal/parser"
	"robotomize/go-lisp/internal/types"
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

	var collectedTokens []lexer.Token

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

			collectedTokens = append(collectedTokens, token)
		}
	}

	p := parser.New(collectedTokens)
	first := true
	for !p.IsEOF() {
		val, err := p.Parse()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		if val == nil {
			break
		}
		if !first {
			fmt.Println()
		}
		first = false
		fmt.Println(types.Pretty(val))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
