package lexer

type NodeHtmlStart struct {
	Info NodeInfo
}

func NewNodeHtmlStart(s string) (NodeHtmlStart, error) {
	n := &NodeHtmlStart{}
	info, err := NewNodeInfo(NodeTypeHtmlStart, s)
	if err != nil {
		return *n, err
	}
	n.Info = info
	return *n, nil
}

func (n NodeHtmlStart) GetInfo() NodeInfo {
	return n.Info
}
