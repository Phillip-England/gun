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

func GetAttributeSlice(n Node) ([]string, error) {
	tagname, err := GetTagName(n)
	if err != nil {
		return []string{}, err
	}
	toks, err := lexer.TokenizeHtml([]rune(n.GetInfo().Value))
	if err != nil {
		return []string{}, err
	}
	s := toks[0].GetLexeme()
	s = strings.Replace(s, "<", "", 1)
	s = stur.ReplaceLast(s, '>', "")
	s = strings.Replace(s, tagname, "", 1)
	parts := strings.Split(s, " ")
	filtered := []string{}
	for _, part := range parts {
		if stur.Squeeze(part) == "" {
			continue
		}
		if len(part) == 1 {
			continue
		}
		filtered = append(filtered, part)
	}
	return filtered, nil
}

func GetRawAttribute(n Node, attrName string) (string, error) {
	attrSlice, err := GetAttributeSlice(n)
	if err != nil {
		return "", err
	}
	for _, rawAttr := range attrSlice {
		if !strings.Contains(rawAttr, "=") {
			if attrName == rawAttr {
				return rawAttr, nil
			}
			continue
		}
		parts := strings.Split(rawAttr, "=")
		if len(parts) < 2 {
			continue
		}
		if parts[0] == attrName {
			return rawAttr, nil
		}
	}
	return "", nil
}

type Attribute struct {
	Name string
	Value string
}

func GetAttribute(n Node, attrName string) (Attribute, bool) {
	attr := &Attribute{}
	rawAttr, err := GetRawAttribute(n, attrName)
	if err != nil {
		return *attr, false
	}
	if !strings.Contains(rawAttr, "=") {
		if rawAttr != "" {
			attr.Name = rawAttr
			return *attr, true
		}	
	}
	rawAttr = strings.Replace(rawAttr, "=", " ", 1)
	parts := strings.Split(rawAttr, " ")
	name := ""
	value := ""
	if len(parts) > 1 {
		name = parts[0]
		value = strings.Join(parts[1:], " ")
	}
	attr.Name = name
	attr.Value = value
	return *attr, true
}

