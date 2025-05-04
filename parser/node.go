package parser

import (
	"strings"

	"github.com/phillip-england/gtml/stur"
)

type Node interface {
	GetInfo() *NodeInfo
}

func AppendChild(parent Node, child Node) {
	parent.GetInfo().Children = append(parent.GetInfo().Children, child)
}

func AppendTextNode(parent Node, text string) {
	parent.GetInfo().TextContent = parent.GetInfo().TextContent+text
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


