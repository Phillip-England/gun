package lexer

type NodeType int

const (
	NodeTypeRoot NodeType = iota
	NodeTypeArg  NodeType = iota
)
