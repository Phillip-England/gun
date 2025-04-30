package parser

import (
	"fmt"
	"strings"

	"github.com/phillip-england/gtml/lexer"
	"github.com/phillip-england/gtml/stur"
)

type Node interface {
	GetInfo() *NodeInfo
}

func AppendChild(parent Node, child Node) {
	parent.GetInfo().Children = append(parent.GetInfo().Children, child)
}

func GetTagName(n Node) (string, error) {
	if n.GetInfo().Type != Normal && n.GetInfo().Type != Void {
		return "", fmt.Errorf("a tag name can only be extracted from a node of type Normal or Void but you attempted on type: %s", n.GetInfo().Type)
	}
	s := n.GetInfo().Value
	runes := []rune(s)
	toks, err := lexer.TokenizeHtml(runes)
	if err != nil {
		return "", err
	}
	name, err := lexer.GetTagName(toks[0])
	if err != nil {
		return "", err
	}
	return name, nil
}

func GetAttributes(n Node) ([]Attribute) {
	val := n.GetInfo().Value
	val = strings.Replace(val, "<", "", 1)
	val = strings.TrimSpace(val)
	bits := strings.Split(val, " ")
	if len(bits) == 1 || len(bits) == 0 {
		return []Attribute{}
	}
	potentialAttrs := stur.SplitWithStringPreserve(strings.Join(bits, " "), ">")
	potentialAttrs = strings.Split(potentialAttrs[0], " ")
	if len(potentialAttrs) == 1 || len(potentialAttrs) == 0 {
		return []Attribute{}
	}
	potentialAttrs = potentialAttrs[1:]

	attrs := []Attribute{}
	for _, attr := range potentialAttrs {
		if !strings.Contains(attr, "=") {
			continue
		}
		parts := strings.Split(attr, "=")
		if len(parts) == 2 {
			name := parts[0]
			val := parts[1]
			attrs = append(attrs, Attribute{
				Name: name,
				Value: val,
			})
		}
		if len(parts) == 1 {
			attrs = append(attrs, Attribute{
				Name: attr,
				Value: "",
			})
		}
	}
	return attrs
}

type Attribute struct {
	Name string
	Value string
}

func GetAttribute(n Node, attrName string) (Attribute, bool) {
	attrs := GetAttributes(n)
	for _, attr := range attrs {
		if attr.Name == attrName {
			return attr, true
		}
	}
	return Attribute{}, false
}


