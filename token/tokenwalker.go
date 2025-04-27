package token

import (
	"fmt"
	"slices"
)

type TokenWalker struct {
	Tokens  []Token
	Current Token
	Pos     int
	Done    bool
	Marked  int
	Buffer  []Token
}

func NewTokenWalker(toks []Token) (*TokenWalker, error) {
	tw := &TokenWalker{}
	tw.Tokens = toks
	tw.Pos = 0
	tw.Done = false
	tw.Marked = 0
	tw.Buffer = []Token{}

	if len(toks) > 0 {
		tw.Current = toks[0]
	} else {
		tw.Done = true
	}

	return tw, nil
}

func (tw *TokenWalker) MarkPos() {
	tw.Marked = tw.Pos
}

func (tw *TokenWalker) FlushFromMarkedPos() []Token {
	if tw.Marked < 0 || tw.Marked >= len(tw.Tokens) {
		return nil
	}
	if tw.Pos < tw.Marked {
		return nil
	}
	collected := tw.Tokens[tw.Marked : tw.Pos+1] // include current token
	tw.Buffer = []Token{}
	return collected
}

func (tw *TokenWalker) Step() {
	if tw.Pos+1 > len(tw.Tokens)-1 {
		tw.Done = true
		return
	}
	tw.Pos += 1
	tw.Current = tw.Tokens[tw.Pos]
}

func (tw *TokenWalker) StepBack() {
	if tw.Pos-1 < 0 {
		// Already at the beginning; can't step back.
		tw.Done = true
		return
	}
	tw.Pos -= 1
	tw.Current = tw.Tokens[tw.Pos]
	tw.Done = false // If you had marked it done stepping forward, stepping back should re-enable it
}


func (tw *TokenWalker) Draw() Token {
	return tw.Current
}

func (tw *TokenWalker) Type() TokenType {
	return tw.Current.Type
}

func (tw *TokenWalker) WalkTo(search TokenType) {
	for {
		if tw.Done {
			return
		}
		if tw.Type() == search {
			return
		}
		tw.Step()
	}
}

func (tw *TokenWalker) WalkUntil(searches []TokenType) {
	for !tw.Done {
		if slices.Contains(searches, tw.Type()) {
			return
		}
		tw.Step()
	}
}

// WalkToEnd steps forward until the end of the tokens.
func (tw *TokenWalker) WalkToEnd() {
	for {
		if tw.Done {
			return
		}
		tw.Step()
	}
}

func (tw *TokenWalker) WalkBack(target TokenType) {
	for {
		if tw.Pos <= 0 {
			tw.Done = true
			return
		}
		if tw.Type() == target {
			return
		}
		tw.StepBack()
	}
}

func (tw *TokenWalker) WalkBackUntil(targets []TokenType) {
	for {
		if tw.Pos <= 0 {
			tw.Done = true
			return
		}
		if slices.Contains(targets, tw.Type()) {
			return
		}
		tw.StepBack()
	}
}


func (tw *TokenWalker) SplitIntoElementStrings() ([]string, error) {
	out := []string{}
	for {
		if tw.Done {
			break
		}
		if tw.Matches([]TokenType{
			TokenHtmlOpenTagOpeningBracket,
			TokenHtmlSelfClosingTagOpeningBracket,
		}) {
			tw.MarkPos()
			err := tw.WalkToElementClose()
			if err != nil {
				return out, err
			}
			out = append(out, Construct(tw.FlushFromMarkedPos()))
		}
		tw.Step()
	}
	return out, nil
}

func (tw *TokenWalker) SplitIntoChildTokenMatrix() ([][]Token, error) {
	tokMatrix := [][]Token{}
	elms, err := tw.SplitIntoElementStrings()
	if err != nil {
		return tokMatrix, err
	}
	for _, elm := range elms {
		toks, err := Tokenize(elm)
		if err != nil {
			return tokMatrix, err
		}
		tokMatrix = append(tokMatrix, toks)
	}
	return tokMatrix, nil
}

func (tw *TokenWalker) WalkToElementClose() (error) {
	if !tw.Matches([]TokenType{
		TokenHtmlOpenTagOpeningBracket,
		TokenHtmlSelfClosingTagOpeningBracket,
	}) {
		return fmt.Errorf(`TOKENWALKER ERROR: attempted to walk to element close, but that function may only be called if your cursor is position directly at an elements starting bracket, you were at token of type: %s`, tw.Type())
	}
	bracketType := tw.Type()
	tw.Step()
	if !tw.Matches([]TokenType{
		TokenHtmlOpenTagName,
		TokenHtmlSelfClosingTagName,
	}) {
		return fmt.Errorf(`TOKENWALKER ERROR: found ourselves at an opening tag bracket, but then took a step, and failed to find a tag name: %s`, tw.Type())
	}
	if bracketType == TokenHtmlSelfClosingTagOpeningBracket {
		tw.WalkUntil([]TokenType{TokenHtmlSelfClosingTagClosingBracket})
		return nil
	}
	startTagName := tw.Draw().Lexeme


	count := 0
	for {
		// if tw.Done {
		// 	return fmt.Errorf(`TOKENWALKER ERROR: was searching for the closing element with the tagname: %s, but failed to locate it`, startTagName)
		// }
		if tw.Type() == TokenHtmlOpenTagName && tw.Draw().Lexeme == startTagName {
			count+=1
		}
		if tw.Type() == TokenHtmlCloseTagName && tw.Draw().Lexeme == startTagName {
			count-=1
		}
		if count == 0 && startTagName == tw.Draw().Lexeme {
			tw.Step()
			break
		}
		tw.Step()
	}
	return nil
}


func (tw *TokenWalker) Peek(n int) (Token, bool) {
	target := tw.Pos + n
	if target < 0 || target >= len(tw.Tokens) {
		var empty Token
		return empty, false
	}
	return tw.Tokens[target], true
}


func (tw *TokenWalker) Expect(expected TokenType) error {
	if tw.Done {
		return fmt.Errorf("expect: reached end of tokens, expected %v", expected)
	}
	if tw.Type() != expected {
		return fmt.Errorf("expect: expected %v but found %v", expected, tw.Type())
	}
	return nil
}

func (tw *TokenWalker) Matches(types []TokenType) bool {
	for _, t := range types {
		if tw.Type() == t {
			return true
		}
	}
	return false
}