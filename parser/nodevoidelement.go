package parser

type NodeVoidElement struct {
	Info *NodeInfo	
}

func (n *NodeVoidElement) GetInfo() *NodeInfo {
	return n.Info
}

func NewNodeVoidElement(s string, t NodeType) *NodeVoidElement {
	info := NewNodeInfo(s, t)
	return &NodeVoidElement{
		Info: info,
	}
}