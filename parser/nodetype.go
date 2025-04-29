package parser

type NodeType string

const (
	Root NodeType = "Root"
	Void NodeType = "Void"
	Normal NodeType = "Normal"
	Text NodeType = "Text"
)