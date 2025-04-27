package token

import (
	"fmt"
	"strings"

	"github.com/phillip-england/gun/lexer"
	"github.com/phillip-england/gun/logi"
	"github.com/phillip-england/gun/stur"
)

type Token struct {
	Type   TokenType
	Lexeme string
}

func (t Token) String() string {
	return fmt.Sprintf("TYPE: %s\n%s", t.Type, t.Lexeme)
}

func Tokenize(s string) ([]Token, error) {

	toks, err := TokenizeStageOne(s)
	if err != nil {
		return toks, err
	}

	toks, err = TokenizeStageTwo(toks)
	if err != nil {
		return toks, err
	}

	toks, err = TokenizeStageThree(toks)
	if err != nil {
		return toks, err
	}

	toks, err = TokenizeStageFour(toks)
	if err != nil {
		return toks, err
	}

	return toks, nil

}

func TokenizeStageOne(s string) ([]Token, error) {
	toks := []Token{}
	l := lexer.NewLexer(s)
	for {
		if l.Done {
			toks = append(toks, Token{
				Lexeme: "",
				Type: TokenEndOfInput,
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
					Type: TokenHtmlTextWhiteSpace,
				})
			} else {
				toks = append(toks, Token{
					Lexeme: s,
					Type: TokenHtmlTextNode,
				})
			}
			l.Step()
			continue
		}

		if prevTok.Type == TokenHtmlTextNode || prevTok.Type == TokenHtmlTextWhiteSpace {
			l.MarkPos()
			l.WalkToWithQuoteSkip(">")
			l.CollectFromMark()
			s := l.FlushBuffer()
			sq := strings.ReplaceAll(s, " ", "")
			voidElementNames := []string{"<area", "<base", "<br", "<col", "<embed", "<hr", "<img", "<input", "<link", "<meta", "<param", "<source", "<track", "<wbr"}
			isVoidElement := false
			for _, voidName := range voidElementNames {
				if len(voidName)+1 > len(sq)-1 {
					continue
				}
				if stur.StartsWith(sq, voidName) {
					isVoidElement = true
				}
			}
			if isVoidElement {
				toks = append(toks, Token{
					Lexeme: s,
					Type: TokenHtmlSelfClosingTag,
				})
				l.Step()
				continue
			}
			if len(sq) < 2 {
				return toks, fmt.Errorf(`SYNTAX ERR: found html tag with less then 2 characters (not counting spaces): %s`, s)
			}
			if string(sq[1]) == "/" {
				toks = append(toks, Token{
					Lexeme: s,
					Type: TokenHtmlCloseTag,
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

func TokenizeStageTwo(toks []Token) ([]Token, error) {
	out := []Token{}
	for _, tok := range toks {

		if tok.Type == TokenHtmlCloseTag {
			s := stur.Squeeze(tok.Lexeme)
			s = strings.ReplaceAll(s, "</", "")
			s = strings.ReplaceAll(s, ">", "")
			out = append(out, Token{
				Lexeme: "</",
				Type: TokenHtmlCloseTagOpeningBracket,
			})
			out = append(out, Token{
				Lexeme: s,
				Type: TokenHtmlCloseTagName,
			})
			out = append(out, Token{
				Lexeme: ">",
				Type: TokenHtmlCloseTagClosingBracket,
			})
			continue
		}

		if tok.Type == TokenHtmlSelfClosingTag {
			out = append(out, Token{
				Lexeme: "<",
				Type: TokenHtmlSelfClosingTagOpeningBracket,
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
				return out, fmt.Errorf(`SYNTAX ERROR: found TokenHtmlSelfClosingTag that was split by " " but now has a length of 0: %s`, s)
			}
			for index, part := range parts {
				if index == 0 {
					out = append(out, Token{
						Lexeme: part,
						Type: TokenHtmlSelfClosingTagName,
					})
					continue
				}
				if index != len(part)-1 {
					out = append(out, Token{
						Lexeme: " ",
						Type: TokenHtmlTagWhiteSpace,
					})
				}
				if stur.LastChar(part) == "/" {
					part = stur.RemoveLastChar(part)
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
				Lexeme: "/>",
				Type: TokenHtmlSelfClosingTagClosingBracket,
			})
			continue
		}
		
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
				return out, fmt.Errorf(`SYNTAX ERROR: found TokenHtmlOpenTag that was split by " " but now has a length of 0: %s`, s)
			}
			for index, part := range parts {
				if index == 0 {
					out = append(out, Token{
						Lexeme: part,
						Type: TokenHtmlOpenTagName,
					})
					continue
				}
				if index != len(part)-1 {
					out = append(out, Token{
						Lexeme: " ",
						Type: TokenHtmlTagWhiteSpace,
					})
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

		if tok.Type == TokenHtmlTextNode {
			l := lexer.NewLexer(tok.Lexeme)
			for {
				if l.Done {
					break
				}
				l.MarkPos()
				l.WalkToUnescaped("%")
				if l.Peek(2, true) == "%s " || l.Peek(2, true) == "%t " || l.Peek(2, true) == "%d " {
					l.CollectFromMark()
					s := l.FlushBuffer()
					s = stur.RemoveLastChar(s)
					if s != "" {
						out = append(out, Token{
							Lexeme: s,
							Type: TokenHtmlRawText,
						})
					}
					l.MarkPos()
					l.Step()
					l.WalkToUnescaped("%")
					l.CollectFromMark()
					s = l.FlushBuffer()
					out = append(out, Token{
						Lexeme: s,
						Type: TokenHtmlTextArg,
					})
					l.Step()
					continue
				} else {
					l.CollectFromMark()
					s := l.FlushBuffer()
					out = append(out, Token{
						Lexeme: s,
						Type: TokenHtmlRawText,
					})
				}
				l.Step()
			}
			continue
		}

		out = append(out, tok)

	}
	return out, nil
}


func TokenizeStageThree(toks []Token) ([]Token, error) {
	out := []Token{}
	for _, tok := range toks {

		if tok.Type == TokenHtmlRawText {
			s := tok.Lexeme
			s = stur.Squeeze(s)
			if s == "::?" {
				out = append(out, Token{
					Lexeme: tok.Lexeme,
					Type: TokenHtmlElseSymbol,
				})
			} else {
				out = append(out, tok)
			}
			continue
		}

		if tok.Type == TokenHtmlAttribute {
			parts := strings.Split(tok.Lexeme, "=")
			if len(parts) != 2 {
				return out, fmt.Errorf(`SYNTAX ERROR: found an html attribute which was split by a " " but doesnt have 2 parts: %s`, tok.Lexeme)
			}
			attrName := parts[0]
			attrValue := strings.Join(parts[1:], "=")
			out = append(out, Token{
				Lexeme: attrName,
				Type: TokenHtmlAttributeName,
			})
			out = append(out, Token{
				Lexeme: "=",
				Type: TokenHtmlAttributeEqualSign,
			})
			out = append(out, Token{
				Lexeme: attrValue,
				Type: TokenHtmlAttributeValue,
			})
			continue
		}

		if tok.Type == TokenHtmlTextArg {
			s := tok.Lexeme
			if len(s) < 2 {
				return toks, fmt.Errorf(`SYNTAX ERROR: found a percentage arg with a length of less then two: %s`, s)
			}
			typeChar := string(s[1])
			switch typeChar {
			case "s": {
				out = append(out, Token{
					Lexeme: s,
					Type: TokenHtmlTextStringArg,
				})
				break
			}
			case "t": {
				out = append(out, Token{
					Lexeme: s,
					Type: TokenHtmlTextStringArg,
				})
				break
			}
			case "d": {
				out = append(out, Token{
					Lexeme: s,
					Type: TokenHtmlTextStringArg,
				})
				break
			}	
			default: {
				return toks, fmt.Errorf(`SYNTAX ERROR: found a percentage arg whose 2nd character was not a 't', 'd', or 's': %s`, s)
			}
			}
			continue
		}

		out = append(out, tok)
	}
	return out, nil
}

func TokenizeStageFour(toks []Token) ([]Token, error) {
	out := []Token{}
	for _, tok := range toks {

		if tok.Type == TokenHtmlAttributeValue {
			l := lexer.NewLexer(tok.Lexeme)
			for {
				if l.Done {
					break
				}
				l.MarkPos()
				l.WalkToUnescaped("%")
				if l.Peek(2, true) == "%s " || l.Peek(2, true) == "%t " || l.Peek(2, true) == "%d " {
					l.CollectFromMark()
					s := l.FlushBuffer()
					s = stur.RemoveLastChar(s)
					if s != "" {
						out = append(out, Token{
							Lexeme: s,
							Type: TokenHtmlAttributeValuePart,
						})
					}
					l.MarkPos()
					l.Step()
					l.WalkToUnescaped("%")
					l.CollectFromMark()
					s = l.FlushBuffer()
					if len(s) < 2 {
						return toks, fmt.Errorf(`SYNTAX ERROR: found an html attribute arg less than 2 characters long: %s`, s)
					}
					secondChar := string(s[1])
					if secondChar == "s" {
						out = append(out, Token{
							Lexeme: s,
							Type: TokenHtmlAttributeStringArg,
						})
					} else if secondChar == "d" {
						out = append(out, Token{
							Lexeme: s,
							Type: TokenHtmlAttributeBoolArg,
						})
					} else if secondChar == "t" {
						out = append(out, Token{
							Lexeme: s,
							Type: TokenHtmlAttributeIntArg,
						})
					} else {
						return toks, fmt.Errorf(`SYNTAX ERROR: found an html attribute arg without a 's', 't', or 'd' as the second character: %s`, s)
					}
					l.Step()
					continue
				} else {
					l.CollectFromMark()
					s := l.FlushBuffer()
					out = append(out, Token{
						Lexeme: s,
						Type: TokenHtmlAttributeValue,
					})
				}
				l.Step()
			}
			continue
		}
		
		out = append(out, tok)
	}
	return out, nil
}


func Construct(toks []Token) string {
	out := ""
	for _, tok := range toks {
		out = out + tok.Lexeme
		logi.Log(tok.String())
	}
	return out
}
