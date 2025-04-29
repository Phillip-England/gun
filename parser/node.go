package parser

type Node interface {
	GetInfo() *NodeInfo
}

func AppendChild(parent Node, child Node) {
	parent.GetInfo().Children = append(parent.GetInfo().Children, child)
}


