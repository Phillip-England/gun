package token

import (
	"fmt"
	"strings"

	"github.com/phillip-england/gtml/lexer"
	"github.com/phillip-england/gtml/stur"
)

type HtmlTokenType string

const (
	// EmptySpace HtmlTokenType = "EmptySpace"
	HtmlOpen   HtmlTokenType = "HtmlOpen"
	HtmlClose  HtmlTokenType = "HtmlClose"
	HtmlVoid   HtmlTokenType = "HtmlVoid"
	Text       HtmlTokenType = "Text"
)

type HtmlToken struct {
	Lexeme string
	Type   HtmlTokenType
	Line int
	Column int
}

// GetLexeme returns the string content of the token.
func (tok HtmlToken) GetLexeme() string {
	return tok.Lexeme
}

// GetType returns the type of the token as a string.
func (tok HtmlToken) GetType() HtmlTokenType {
	return tok.Type
}

func (tok HtmlToken) GetLine() int {
	return tok.Line
}

func (tok HtmlToken) GetColumn() int {
	return tok.Column
}

// TokenizeHtml tokenizes a slice of runes representing HTML input
// into a list of tokens through two passes: raw token extraction
// and structural classification (e.g., identifying void elements).
func TokenizeHtml(input []rune) ([]Token, error) {
	toks, err := firstPass(input)
	if err != nil {
		return toks, err
	}
	toks, err = secondPass(toks)
	if err != nil {
		return toks, err
	}

	return toks, nil
}

// secondPass processes tokens from the first pass and determines if
// HtmlOpen tokens are actually HtmlVoid (self-closing) by checking
// for corresponding HtmlClose tags later in the sequence.
func secondPass(toks []Token) ([]Token, error) {
	out := []Token{}
	for i, tok := range toks {
		if tok.GetType() != HtmlOpen {
			out = append(out, tok)
			continue
		}
		closingTag, _, err := GetClosingTag(tok, i, toks)
		if err != nil {
			return out, err
		}
		if closingTag == nil {
			out = append(out, HtmlToken{
				Lexeme: tok.GetLexeme(),
				Type:   HtmlVoid,
				Line: tok.GetLine(),
				Column: tok.GetColumn(),
			})
		} else {
			out = append(out, tok)
		}
	}

	return out, nil
}

// GetClosingTag searches for the matching HtmlClose token corresponding
// to the given HtmlOpen token, respecting nesting depth.
// Returns nil if no matching closing tag is found.
func GetClosingTag(tok Token, i int, toks []Token) (Token, int, error) {
	if tok.GetType() == HtmlVoid {
		return nil, -1, nil
	}
	if tok.GetType() != HtmlOpen {
		return tok, -1, fmt.Errorf(`attempted to extract the closing tag from an invalid token: %s`, tok.GetLexeme())
	}
	name := GetTagName(tok)
	found := 1
	for i1, tok1 := range toks {
		if i1 <= i {
			continue
		}
		if tok1.GetType() != HtmlOpen && tok1.GetType() != HtmlClose {
			continue
		}
		name1 := GetTagName(tok1)
		if name != name1 {
			continue
		}
		if tok1.GetType() == HtmlOpen {
			found += 1
		}
		if tok1.GetType() == HtmlClose {
			found -= 1
		}
		if found == 0 && tok1.GetType() == HtmlClose {
			return tok1, i1, nil
		}
	}
	// failed to find a matching closing tag
	return nil, -1, nil
}

func IsSelfContained(toks []Token) (bool, error) {
	innerToks, err := ShedOuterHtml(toks)
	if err != nil {
		return false, err
	}
	if len(innerToks) == len(toks) {
		return false, nil
	}
	return true, nil
}

func ExtractFullElement(tok Token, i int, toks []Token) {
	
}


func ShedOuterHtml(toks []Token) ([]Token, error) {
	out := []Token{}
	if len(toks) == 0 {
		return toks, nil
	}
	firstTok := toks[0]
	if firstTok.GetType() == HtmlClose {
		return out, fmt.Errorf("you cannot shed the outerhtml of an html closing tag: %s", firstTok.GetType())
	}
	_, closeTagIndex, err := GetClosingTag(toks[0], 0, toks)
	if err != nil {
		return out, err
	}
	if closeTagIndex == len(toks)-1 {
		toks = toks[1:]
		toks = toks[:len(toks)-1]
		return toks, nil
	} else {
		return toks, nil
	}
}




// GetTagName extracts the tag name from an HtmlOpen or HtmlClose token's lexeme.
// Strips angle brackets, slashes, and attributes, returning just the tag name.
func GetTagName(tok Token) string {
	if HtmlTokenType(tok.GetType()) != HtmlOpen && HtmlTokenType(tok.GetType()) != HtmlClose && HtmlTokenType(tok.GetType()) != HtmlVoid {
		return ""
	}
	s := tok.GetLexeme()
	s = strings.Replace(s, "<", "", 1)
	s = stur.ReplaceLast(s, '>', "")
	s = stur.ReplaceLast(s, '/', "")
	parts := strings.Split(s, " ")
	for _, part := range parts {
		sq := stur.Squeeze(part)
		if sq != "" {
			return part
		}
	}
	return ""
}

// firstPass performs an initial walk over the input runes and splits the input
// into basic tokens: HtmlOpen, HtmlClose, Text, and EmptySpace.
// It uses the lexer to handle quote-skipping and whitespace correctly.
func firstPass(input []rune) ([]Token, error) {
	toks := []Token{}
	l := lexer.NewLexer(input)
	for {
		if l.Terminated {
			break
		}
		// avoid capturing the leading '<' as text
		if l.Pos != 0 {
			l.Mark()
			l.WalkUntilSkipQuotes('<')
			l.StepBack()
			buf := string(l.FlushFromMark())
			if len(stur.Squeeze(buf)) == 0 {
				// OPTING OUT OF COLLECTING EMPTY SPACE
				// toks = append(toks, HtmlToken{
				// 	Lexeme: buf,
				// 	Type:   EmptySpace,
				// 	Line: l.Line,
				// 	Column: l.Column,
				// })
			} else {
				toks = append(toks, HtmlToken{
					Lexeme: buf,
					Type:   Text,
					Line: l.Line,
					Column: l.Column-len(buf),
				})
			}
			l.Step()
		}
		if l.CharIs("<") {
			l.Mark()
			found := l.WalkUntilSkipQuotes('>')
			if !found {
				return toks, fmt.Errorf(`SYNTAX ERROR: failed to close html element: %s`, string(input))
			}
			buf := string(l.FlushFromMark())
			sq := stur.Squeeze(buf)
			if len(sq) > 2 && sq[1] == '/' {
				toks = append(toks, HtmlToken{
					Lexeme: buf,
					Type:   HtmlClose,
					Line: l.Line,
					Column: l.Column-len(buf),
				})
			} else {
				toks = append(toks, HtmlToken{
					Lexeme: buf,
					Type:   HtmlOpen,
					Line: l.Line,
					Column: l.Column-len(buf),
				})
			}
			l.Step()
			continue
		}
		l.Step()
	}
	return toks, nil
}
