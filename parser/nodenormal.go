package parser

type NodeNormal struct {
	Info *NodeInfo	
}

func (n *NodeNormal) GetInfo() *NodeInfo {
	return n.Info
}

func NewNodeNormal(s string, t NodeType) Node {
	info := NewNodeInfo(s, t)
	return &NodeNormal{
		Info: info,
	}
}