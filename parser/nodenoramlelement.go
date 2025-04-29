package parser

type NodeNoramlElement struct {
	Info *NodeInfo	
}

func (n *NodeNoramlElement) GetInfo() *NodeInfo {
	return n.Info
}

func NewNodeNoramlElement(s string, t NodeType) *NodeNoramlElement {
	info := NewNodeInfo(s, t)
	return &NodeNoramlElement{
		Info: info,
	}
}