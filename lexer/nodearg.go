package lexer

type NodeArg struct {
	Info NodeInfo
}

func NewNodeArg(s string) (NodeArg, error) {
	n := &NodeArg{}
	info, err := NewNodeInfo(NodeTypeArg, s)
	if err != nil {
		return *n, err
	}
	n.Info = info
	return *n, nil
}

func (n NodeArg) GetInfo() NodeInfo {
	return n.Info
}
