package parser

import (
	"fmt"

	"github.com/phillip-england/gtml/lexer"
)

type Document struct {
	Info *NodeInfo	
}

func (elm *Document) GetInfo() *NodeInfo {
	return elm.Info
}


func makeAst(elm Node, toks []lexer.Token) (Node, error) {
	isSelfContained, err := lexer.IsSelfContained(toks)
	if err != nil {
		return elm, err
	}
	if isSelfContained {
		if len(toks) == 1 {
			return elm, nil
		}
		child := NewNodeNoramlElement(lexer.Construct(toks), Normal)
		AppendChild(elm, child)
	}
	innerToks, err := lexer.ShedOuterHtml(toks)
	if err != nil {
		return elm, err
	}
	for i, tok := range innerToks {
		if tok.GetType() == string(lexer.HtmlVoid) {
			child := NewNodeVoidElement(tok.GetLexeme(), Void)
			AppendChild(elm, child)
			continue
		}
		if tok.GetType() == string(lexer.Text) {
			child := NewNodeText(tok.GetLexeme(), Text)
			AppendChild(elm, child)
			continue
		}
		if tok.GetType() == string(lexer.HtmlOpen) {
			_, endTagIndex, err := lexer.GetClosingTag(tok, i, innerToks)
			if err != nil {
				return elm, err
			}
			elmToks := innerToks[i:endTagIndex+1]
			child := NewNodeNoramlElement(lexer.Construct(elmToks), Normal)
			childElm, ok := any(child).(Node)
			if !ok {
				return childElm, fmt.Errorf("failed to assert Document to Element")
			}
			AppendChild(elm, childElm)
			continue
		}
		
	}
	return elm, nil
}


func NewAst(toks []lexer.Token) (Node, error) {
	toks = lexer.RemoveEmptySpace(toks)
	isSelfContainer, err := lexer.IsSelfContained(toks)
	if err != nil {
		return nil, err
	}
	if !isSelfContainer {
		spanStart := lexer.HtmlToken{
			Lexeme: "<span>",
			Type:   lexer.HtmlOpen,
		}
		spanEnd := lexer.HtmlToken{
			Lexeme: "</span>",
			Type:   lexer.HtmlClose,
		}
		toks = append(toks, spanEnd)
		toks = append([]lexer.Token{spanStart}, toks...)
	}
	doc := &Document{
		Info: NewNodeInfo("", Root),
	}
	docElm, ok := any(doc).(Node)
	if !ok {
		return docElm, fmt.Errorf("failed to assert Document to Element")
	}
	docElm, err = makeAst(docElm, toks)
	if err != nil {
		return docElm, err
	}
	return docElm, nil
}


func WalkNodes(node Node, fn func(i int, n Node) error) error {
	for i, child := range node.GetInfo().Children {
		err := fn(i, child)
		if err != nil {
			return err
		}
		WalkNodes(child, fn)
	}
	return nil
}