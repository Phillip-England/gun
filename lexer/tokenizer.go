package lexer

type LexerState int

const (
	LexerStateBlank     LexerState = iota
	LexerStateHtmlStart LexerState = iota
	LexerStateHtmlAttrs LexerState = iota
)

func Tokenize(s string) ([]Node, error) {
	nodes := []Node{}
	buf := ""
	state := LexerStateBlank
	for _, char := range s {
		ch := string(char)

	}
	return nodes, nil
}
