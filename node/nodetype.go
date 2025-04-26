package lexer

type NodeType int

const (
	NodeTypeRoot      NodeType = iota
	NodeTypeHtmlStart NodeType = iota
	NodeTypeArg       NodeType = iota
)
