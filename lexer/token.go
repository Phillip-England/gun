package lexer

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	TokenElementOpenTag TokenType = iota
	TokenElementCloseTag
	TokenElementText
)

type Token struct {
	Lexeme string
	Type   TokenType
}

func TokenizeHtml(s string) ([]Token, error) {
	toks := []Token{}
	l := NewLexer(s)
	for {
		if l.Done {
			break
		}

		// jumping to the first "<"
		if len(toks) == 0 {
			l.SkipWhitespace()
			if l.Char() != "<" {
				return toks, fmt.Errorf(`SYNTAX ERR: expected first character to be "<" but found "%s"`, l.Char())
			}
			l.MarkPos()
			l.WalkToWithQuoteSkip(">")
			l.CollectFromMark()
			s := l.FlushBuffer()
			toks = append(toks, Token{
				Lexeme: s,
				Type: TokenElementOpenTag,
			})
			l.Step()
			continue
		}

		prevTok := toks[len(toks)-1] // used to determine where we are

		if prevTok.Type == TokenElementOpenTag || prevTok.Type == TokenElementCloseTag {
			l.MarkPos()
			l.WalkToWithQuoteSkip("<")
			l.StepBack()
			l.CollectFromMark()
			s := l.FlushBuffer()
			toks = append(toks, Token{
				Lexeme: s,
				Type: TokenElementText,
			})
			l.Step()
			continue
		}

		if prevTok.Type == TokenElementText {
			l.MarkPos()
			l.WalkToWithQuoteSkip(">")
			l.CollectFromMark()
			s := l.FlushBuffer()
			sq := strings.ReplaceAll(s, " ", "")
			if len(sq) < 2 {
				return toks, fmt.Errorf(`SYNTAX ERR: found html tag with less then 2 characters (not counting spaces): %s`, s)
			}
			if string(sq[1]) == "/" {
				toks = append(toks, Token{
					Lexeme: s,
					Type: TokenElementCloseTag,
				})
			} else {
				toks = append(toks, Token{
					Lexeme: s,
					Type: TokenElementOpenTag,
				})
			}
			l.Step()
			continue
		}

		l.Step() // increment
	}

	return toks, nil
}
