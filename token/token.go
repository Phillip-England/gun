package token

import (
	"fmt"

	"github.com/phillip-england/gtml/lexer"
)

type TokenType string

type Token struct {
	Lexeme string
	Type TokenType
}

func Tokenize(runes []rune) (string, error) {

	l := lexer.NewLexer(runes)

	for {
		if l.Terminated {
			break
		}
		l.SkipWhiteSpace()
		fmt.Println(l.Char)
		l.Step()
	}
	

	


	return "", nil
}
