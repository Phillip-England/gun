package lexer

type NodeHtmlStart struct {
	Info NodeInfo
}

func NewNodeHtmlStart() (NodeHtmlStart, error) {
	n := &NodeHtmlStart{}
	info, err := NewNodeInfo(NodeTypeHtmlStart, "<")
	if err != nil {
		return *n, err
	}
	n.Info = info
	return *n, nil
}

func (n NodeHtmlStart) GetInfo() NodeInfo {
	return n.Info
}
