package parser

type NodeNoraml struct {
	Info *NodeInfo	
}

func (n *NodeNoraml) GetInfo() *NodeInfo {
	return n.Info
}

func NewNodeNoraml(s string, t NodeType) *NodeNoraml {
	info := NewNodeInfo(s, t)
	return &NodeNoraml{
		Info: info,
	}
}