package token

import "fmt"

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