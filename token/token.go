package token

import (
	"fmt"
	"strings"

	"github.com/phillip-england/gun/lexer"
	"github.com/phillip-england/gun/stur"
)

type Token struct {
	Type   TokenType
	Lexeme string
}

// converts standard html input into tokens
func TokenizeHtml(s string) ([]Token, error) {
	toks := []Token{}
	l := lexer.NewLexer(s)
	for {
		if l.Done {
			toks = append(toks, Token{
				Lexeme: "",
				Type: TokenEndOfFile,
			})
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
				Type: TokenHtmlOpenTag,
			})
			l.Step()
			continue
		}

		prevTok := toks[len(toks)-1] // used to determine where we are

		if prevTok.Type == TokenHtmlOpenTag || prevTok.Type == TokenHtmlCloseTag || prevTok.Type == TokenHtmlSelfClosingTag {
			l.MarkPos()
			l.WalkToWithQuoteSkip("<")
			l.StepBack()
			l.CollectFromMark()
			s := l.FlushBuffer()
			sq := stur.Squeeze(s)
			if sq == "" {
				toks = append(toks, Token{
					Lexeme: s,
					Type: TokenHtmlWhiteSpace,
				})
			} else {
				toks = append(toks, Token{
					Lexeme: s,
					Type: TokenHtmlText,
				})
			}
			l.Step()
			continue
		}

		if prevTok.Type == TokenHtmlText || prevTok.Type == TokenHtmlWhiteSpace {
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
					Type: TokenHtmlCloseTag,
				})
			} else if string(sq[len(sq)-2]) == "/" {
				toks = append(toks, Token{
					Lexeme: s,
					Type: TokenHtmlSelfClosingTag,
				})
			} else {
				toks = append(toks, Token{
					Lexeme: s,
					Type: TokenHtmlOpenTag,
				})
			}
			l.Step()
			continue
		}

		l.Step() // increment
	}

	return toks, nil
}


// takes in a series of tokens, checks for deconstructable tokens,
// and breaks them down into their smaller bits
func Deconstruct(toks []Token) ([]Token, error) {
	out := []Token{}
	for _, tok := range toks {
		
		if tok.Type == TokenHtmlOpenTag {
			out = append(out, Token{
				Lexeme: "<",
				Type: TokenHtmlOpenTagOpeningBracket,
			})
			s := tok.Lexeme
			s = strings.Replace(s, "<", "", 1)
			s = s[:len(s)-1]
			parts := stur.SplitWithStringPreserve(s, " ")
			// just a tagname no attributes
			if len(parts) == 1 {
				out = append(out, Token{
					Lexeme: parts[0],
					Type: TokenHtmlOpenTagName,
				})
				out = append(out, Token{
					Lexeme: ">",
					Type: TokenHtmlOpenTagClosingBracket,
				})
				continue
			}
			if len(parts) == 0 {
				return out, fmt.Errorf(`SYNTAX ERROR: found TokenElementOpenTag that was split by " " but now has a length of 0: %s`, s)
			}
			for index, part := range parts {
				if index == 0 {
					out = append(out, Token{
						Lexeme: part,
						Type: TokenHtmlOpenTagName,
					})
					continue
				}
				if strings.Contains(part, "=") {
					out = append(out, Token{
						Lexeme: part,
						Type: TokenHtmlAttribute,
					})
				} else {
					out = append(out, Token{
						Lexeme: part,
						Type: TokenHtmlBooleanAttribute,
					})
				}
			}
			out = append(out, Token{
				Lexeme: ">",
				Type: TokenHtmlOpenTagClosingBracket,
			})
			continue
		}

		out = append(out, tok)

	}
	return out, nil
}