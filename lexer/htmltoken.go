package lexer

import (
	"fmt"
	"strings"

	"github.com/phillip-england/gtml/stur"
)

type HtmlTokenType string

const (
	EmptySpace HtmlTokenType = "EmptySpace"
	HtmlOpen   HtmlTokenType = "HtmlOpen"
	HtmlClose  HtmlTokenType = "HtmlClose"
	HtmlVoid   HtmlTokenType = "HtmlVoid"
	Text       HtmlTokenType = "Text"
)

type HtmlToken struct {
	Lexeme string
	Type   HtmlTokenType
}

// GetLexeme returns the string content of the token.
func (tok HtmlToken) GetLexeme() string {
	return tok.Lexeme
}

// GetType returns the type of the token as a string.
func (tok HtmlToken) GetType() string {
	return string(tok.Type)
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
		if tok.GetType() != string(HtmlOpen) {
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
	if tok.GetType() == string(HtmlVoid) {
		return nil, -1, nil
	}
	if tok.GetType() != string(HtmlOpen) {
		return tok, -1, fmt.Errorf(`attempted to extract the closing tag from an invalid token: %s`, tok.GetLexeme())
	}
	name, err := GetTagName(tok)
	if err != nil {
		return tok, -1, err
	}
	found := 1
	for i1, tok1 := range toks {
		if i1 <= i {
			continue
		}
		if tok1.GetType() != string(HtmlOpen) && tok1.GetType() != string(HtmlClose) {
			continue
		}
		name1, err := GetTagName(tok1)
		if err != nil {
			return tok, -1, err
		}
		if name != name1 {
			continue
		}
		if tok1.GetType() == string(HtmlOpen) {
			found += 1
		}
		if tok1.GetType() == string(HtmlClose) {
			found -= 1
		}
		if found == 0 && tok1.GetType() == string(HtmlClose) {
			return tok1, i1, nil
		}
	}
	// failed to find a matching closing tag
	return nil, -1, nil
}

func IsSelfContained(toks []Token) (bool, error) {
	toks = RemoveEmptySpace(toks)
	innerToks, err := ShedOuterHtml(toks)
	if err != nil {
		return false, err
	}
	if len(innerToks) == len(toks) {
		return false, nil
	}
	return true, nil
}

// RemoveEmptySpace filters out all tokens of type EmptySpace from the provided token slice.
func RemoveEmptySpace(toks []Token) []Token {
	filtered := make([]Token, 0, len(toks))
	for _, tok := range toks {
		if tok.GetType() != string(EmptySpace) {
			filtered = append(filtered, tok)
		}
	}
	return filtered
}


func ShedOuterHtml(toks []Token) ([]Token, error) {
	out := []Token{}
	if len(toks) == 0 {
		return toks, nil
	}
	firstTok := toks[0]
	if firstTok.GetType() == string(HtmlClose) {
		return out, fmt.Errorf("you cannot shed the outerhtml of an html closing tag: %s", firstTok.GetType())
	}
	filteredToks := []Token{}
	for _, tok := range toks {
		if tok.GetType() == string(EmptySpace) {
			continue
		}
		filteredToks = append(filteredToks, tok)
	}
	_, closeTagIndex, err := GetClosingTag(filteredToks[0], 0, filteredToks)
	if err != nil {
		return out, err
	}
	if closeTagIndex == len(filteredToks)-1 {
		filteredToks = filteredToks[1:]
		filteredToks = filteredToks[:len(filteredToks)-1]
		return filteredToks, nil
	} else {
		return filteredToks, nil
	}
}




// GetTagName extracts the tag name from an HtmlOpen or HtmlClose token's lexeme.
// Strips angle brackets, slashes, and attributes, returning just the tag name.
func GetTagName(tok Token) (string, error) {
	if HtmlTokenType(tok.GetType()) != HtmlOpen && HtmlTokenType(tok.GetType()) != HtmlClose {
		return "", fmt.Errorf(`tag names can only be extracted from HtmlOpen and HtmlClose but you attempted to extract on: %s`, tok)
	}
	s := tok.GetLexeme()
	s = strings.Replace(s, "<", "", 1)
	s = stur.ReplaceLast(s, '>', "")
	s = stur.ReplaceLast(s, '/', "")
	parts := strings.Split(s, " ")
	for _, part := range parts {
		sq := stur.Squeeze(part)
		if sq != "" {
			return part, nil
		}
	}
	return "", fmt.Errorf(`failed to extract tag name from: %s`, tok.GetLexeme())
}

// firstPass performs an initial walk over the input runes and splits the input
// into basic tokens: HtmlOpen, HtmlClose, Text, and EmptySpace.
// It uses the lexer to handle quote-skipping and whitespace correctly.
func firstPass(input []rune) ([]Token, error) {
	toks := []Token{}
	l := NewLexer(input)
	for {
		if l.Terminated {
			break
		}
		// avoid capturing the leading '<' as text
		if l.Pos != 0 {
			l.Mark()
			l.WalkUntilSkipQuotes('<')
			l.StepBack()
			if len(stur.Squeeze(string(l.FlushFromMark()))) == 0 {
				toks = append(toks, HtmlToken{
					Lexeme: string(l.FlushFromMark()),
					Type:   EmptySpace,
				})
			} else {
				toks = append(toks, HtmlToken{
					Lexeme: string(l.FlushFromMark()),
					Type:   Text,
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
			if strings.Contains(buf, "/") {
				toks = append(toks, HtmlToken{
					Lexeme: string(l.FlushFromMark()),
					Type:   HtmlClose,
				})
			} else {
				toks = append(toks, HtmlToken{
					Lexeme: string(l.FlushFromMark()),
					Type:   HtmlOpen,
				})
			}
			l.Step()
			continue
		}
		l.Step()
	}
	return toks, nil
}
