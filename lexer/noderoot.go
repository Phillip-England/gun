package lexer

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type NodeRoot struct {
	Info NodeInfo
}

func NewNodeRoot(s string) (NodeRoot, error) {
	n := &NodeRoot{}
	info, err := NewNodeInfo(NodeTypeRoot, s)
	if err != nil {
		return *n, err
	}
	n.Info = info

	_, err = goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return *n, err
	}

	return *n, nil
}

func (n NodeRoot) GetInfo() NodeInfo {
	return n.Info
}
