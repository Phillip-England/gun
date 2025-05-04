package lexer

import "github.com/phillip-england/gtml/logi"

type Token interface {
	GetLexeme() string
	GetType() HtmlTokenType
	GetLine() int
	GetColumn() int
}

func LogTokens(toks []Token) {
	logi.Log(Construct(toks))
	for _, tok := range toks {
		logi.Log(tok.GetType(), tok.GetLexeme())
	}
}

func Construct(toks[]Token) string {
	out := ""
	for _, tok := range toks {
		out += tok.GetLexeme()
	}
	return out
}

