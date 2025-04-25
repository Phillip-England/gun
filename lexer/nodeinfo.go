package lexer

type NodeInfo struct {
	Type     NodeType
	Children []Node
}

func NewNodeInfo(t NodeType) (NodeInfo, error) {
	info := &NodeInfo{}
	info.Type = t
	return *info, nil
}
