package node

type NodeType string

const (
	NodeTypeRoot     NodeType = "NodeTypeRoot"
	NodeTypeText     NodeType = "NodeTypeText"
	NodeTypeArg      NodeType = "NodeTypeArg"
	NodeTypeFor      NodeType = "NodeTypeFor"
	NodeTypeIf       NodeType = "NodeTypeIf"
	NodeTypeElse     NodeType = "NodeTypeElse"
	NodeTypeIfHtml   NodeType = "NodeTypeIfHtml"
	NodeTypeElseHtml NodeType = "NodeTypeElseHtml"
)

type Node struct {
	Value string
	Type NodeType
	Children []*Node
}

func NewNode(value string) (*Node) {
	n := &Node{}
	n.Value = value
	return n
} 

func (n *Node) Append(child *Node) {
	n.Children = append(n.Children, child)
}

