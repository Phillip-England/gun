package lexer

type NodeInfo struct {
	Type     NodeType
	Children []Node
	String   string
}

func NewNodeInfo(t NodeType, s string) (NodeInfo, error) {
	info := &NodeInfo{}
	info.Type = t
	info.String = s
	return *info, nil
}
