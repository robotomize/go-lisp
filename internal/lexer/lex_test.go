package lexer

import "testing"

func TestNextList(t *testing.T) {
	l := New("(+ 1 2)")
	want := []Token{
		{Typ: TokenTypOpen, Value: "("},
		{Typ: TokenTypSymbol, Value: "+"},
		{Typ: TokenTypNumber, Value: "1"},
		{Typ: TokenTypNumber, Value: "2"},
		{Typ: TokenTypClose, Value: ")"},
	}

	for i, expected := range want {
		got, ok := l.Next()
		if !ok {
			t.Fatalf("token %d: lexer stopped early", i)
		}
		if got != expected {
			t.Fatalf("token %d: got %#v, want %#v", i, got, expected)
		}
	}

	if _, ok := l.Next(); ok {
		t.Fatal("lexer produced a token after the input ended")
	}
}

func TestQuoteToken(t *testing.T) {
	l := New("'x")
	want := []Token{
		{Typ: TokenTypQuote, Value: "'"},
		{Typ: TokenTypSymbol, Value: "x"},
	}

	for i, expected := range want {
		got, ok := l.Next()
		if !ok {
			t.Fatalf("token %d: lexer stopped early", i)
		}
		if got != expected {
			t.Fatalf("token %d: got %#v, want %#v", i, got, expected)
		}
	}
}
