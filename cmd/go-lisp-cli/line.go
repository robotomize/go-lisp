package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

var errAbort = errors.New("aborted")

type keyKind int

const (
	keyIgnore keyKind = iota
	keyRune
	keyLeft
	keyRight
	keyUp
	keyDown
	keyBackspace
	keyDelete
	keyEnter
	keyCtrlC
	keyCtrlD
)

type key struct {
	kind keyKind
	r    rune
}

type editState struct {
	buf     []rune
	cursor  int
	history []string
	histPos int
	draft   []rune
}

func (s *editState) insert(r rune) {
	s.buf = append(s.buf, 0)
	copy(s.buf[s.cursor+1:], s.buf[s.cursor:])
	s.buf[s.cursor] = r
	s.cursor++
}

func (s *editState) left() {
	if s.cursor > 0 {
		s.cursor--
	}
}

func (s *editState) right() {
	if s.cursor < len(s.buf) {
		s.cursor++
	}
}

func (s *editState) backspace() {
	if s.cursor == 0 {
		return
	}
	s.buf = append(s.buf[:s.cursor-1], s.buf[s.cursor:]...)
	s.cursor--
}

func (s *editState) delete() {
	if s.cursor >= len(s.buf) {
		return
	}
	s.buf = append(s.buf[:s.cursor], s.buf[s.cursor+1:]...)
}

func (s *editState) up() {
	if s.histPos == 0 {
		return
	}
	if s.histPos == len(s.history) {
		s.draft = append([]rune(nil), s.buf...)
	}
	s.histPos--
	s.buf = []rune(s.history[s.histPos])
	s.cursor = len(s.buf)
}

func (s *editState) down() {
	if s.histPos >= len(s.history) {
		return
	}
	s.histPos++
	if s.histPos == len(s.history) {
		s.buf = append([]rune(nil), s.draft...)
	} else {
		s.buf = []rune(s.history[s.histPos])
	}
	s.cursor = len(s.buf)
}

func (s *editState) submit() string {
	text := string(s.buf)
	if text != "" {
		s.history = append(s.history, text)
	}
	s.buf = nil
	s.cursor = 0
	s.draft = nil
	s.histPos = len(s.history)
	return text
}

func (s *editState) abort() {
	s.buf = nil
	s.cursor = 0
	s.draft = nil
	s.histPos = len(s.history)
}

type lineReader interface {
	next(prompt string) (string, error)
}

type scannerReader struct {
	scanner *bufio.Scanner
	out     io.Writer
}

func (r scannerReader) next(prompt string) (string, error) {
	if prompt != "" {
		fmt.Fprint(r.out, prompt)
	}
	if !r.scanner.Scan() {
		if err := r.scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return r.scanner.Text(), nil
}

type termReader struct {
	in    *os.File
	out   *os.File
	state editState
}

func (r *termReader) next(prompt string) (string, error) {
	fd := int(r.in.Fd())
	raw, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}

	defer term.Restore(fd, raw)

	r.redraw(prompt)
	input := bufio.NewReader(r.in)
	for {
		k, err := readKey(input)
		if err != nil {
			return "", err
		}
		switch k.kind {
		case keyRune:
			r.state.insert(k.r)
		case keyLeft:
			r.state.left()
		case keyRight:
			r.state.right()
		case keyUp:
			r.state.up()
		case keyDown:
			r.state.down()
		case keyBackspace:
			r.state.backspace()
		case keyDelete:
			r.state.delete()
		case keyEnter:
			fmt.Fprint(r.out, "\r\n")
			return r.state.submit(), nil
		case keyCtrlC:
			fmt.Fprint(r.out, "^C\r\n")
			r.state.abort()
			return "", errAbort
		case keyCtrlD:
			if len(r.state.buf) == 0 {
				fmt.Fprint(r.out, "\r\n")
				return "", io.EOF
			}
		}
		r.redraw(prompt)
	}
}

func (r *termReader) redraw(prompt string) {
	fmt.Fprintf(r.out, "\r%s%s\033[K", prompt, string(r.state.buf))
	tail := len(r.state.buf) - r.state.cursor
	if tail > 0 {
		fmt.Fprintf(r.out, "\033[%dD", tail)
	}
}

func readKey(r *bufio.Reader) (key, error) {
	b, err := r.ReadByte()
	if err != nil {
		return key{}, err
	}

	switch b {
	case 3:
		return key{kind: keyCtrlC}, nil
	case 4:
		return key{kind: keyCtrlD}, nil
	case 8, 127:
		return key{kind: keyBackspace}, nil
	case '\r', '\n':
		return key{kind: keyEnter}, nil
	case 0x1b:
		return readEscape(r)
	}

	if b < 0x80 {
		if b < 0x20 {
			return key{kind: keyIgnore}, nil
		}
		return key{kind: keyRune, r: rune(b)}, nil
	}

	if err := r.UnreadByte(); err != nil {
		return key{}, err
	}
	rn, _, err := r.ReadRune()
	if err != nil {
		return key{}, err
	}
	return key{kind: keyRune, r: rn}, nil
}

func readEscape(r *bufio.Reader) (key, error) {
	b, err := r.ReadByte()
	if err != nil {
		return key{}, err
	}
	if b != '[' && b != 'O' {
		return key{kind: keyIgnore}, nil
	}

	var seq []byte
	for {
		c, err := r.ReadByte()
		if err != nil {
			return key{}, err
		}
		seq = append(seq, c)
		if c >= 0x40 && c <= 0x7e {
			break
		}
	}

	switch seq[len(seq)-1] {
	case 'A':
		return key{kind: keyUp}, nil
	case 'B':
		return key{kind: keyDown}, nil
	case 'C':
		return key{kind: keyRight}, nil
	case 'D':
		return key{kind: keyLeft}, nil
	case '~':
		if len(seq) >= 2 && seq[0] == '3' {
			return key{kind: keyDelete}, nil
		}
	}
	return key{kind: keyIgnore}, nil
}
