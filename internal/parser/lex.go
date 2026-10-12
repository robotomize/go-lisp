package parser

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func stripLeft(text string) string {
	for len(text) > 0 {
		r, size := utf8.DecodeRuneInString(text)
		switch {
		case r == ',' || unicode.IsSpace(r):
			text = text[size:]
		case r == ';':
			end := strings.IndexByte(text, '\n')
			if end == -1 {
				return ""
			}

			text = text[end:]
		default:
			return text
		}
	}

	return text
}
