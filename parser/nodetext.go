package parser

type NodeText struct {
	Info *NodeInfo	
}

func (n *NodeText) GetInfo() *NodeInfo {
	return n.Info
}

func NewNodeText(s string, t NodeType) *NodeText {
	info := NewNodeInfo(s, t)
	return &NodeText{
		Info: info,
	}
}