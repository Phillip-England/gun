package lexer

import (
	"fmt"
	"testing"
)


func TestLexer(t *testing.T) {

	input := []rune("\n\n\n\n\n    hello\n\n\nhello\n")
	l := NewLexer(input)
	for {
		if l.Terminated {
			break
		}
		if l.Char() != "h" {
			l.Step()
			continue
		}
		fmt.Println(l.Line, l.Column)
		fmt.Println(l.Char())
		l.Step()
	}

}