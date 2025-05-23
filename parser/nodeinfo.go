package parser

type NodeInfo struct {
	Value string
	Children []Node
	Type NodeType
	TextContent string
}

func NewNodeInfo(val string, t NodeType) *NodeInfo {
	return &NodeInfo{
		Value: val,
		Children: make([]Node, 0),
		Type: t,
	}
}