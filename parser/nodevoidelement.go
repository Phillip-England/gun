package parser

type NodeVoid struct {
	Info *NodeInfo	
}

func (n *NodeVoid) GetInfo() *NodeInfo {
	return n.Info
}

func NewNodeVoid(s string, t NodeType) *NodeVoid {
	info := NewNodeInfo(s, t)
	return &NodeVoid{
		Info: info,
	}
}