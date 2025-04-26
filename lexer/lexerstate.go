package lexer

type LexerState int

const (
	LexerStateInit LexerState = iota
	LexerStateHtmlStartTagOpen
	LexerStateHtmlStartTagClose
)
