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




func NewAst(toks []lexer.Token) (Node, error) {
	doc, toks, err := initDoc(toks)
	if err != nil {
		return doc, err
	}
	docElm, ok := any(doc).(Node)
	if !ok {
		return docElm, fmt.Errorf("failed to assert Document to Element")
	}

	docElm, err = firstPass(docElm, toks)
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

func initDoc(toks []lexer.Token) (*Document, []lexer.Token, error) {
	toks = lexer.RemoveEmptySpace(toks)
	isSelfContainer, err := lexer.IsSelfContained(toks)
	if err != nil {
		return nil, toks, err
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
		Info: NewNodeInfo(lexer.Construct(toks), Root),
	}
	return doc, toks, nil
}

func firstPass(elm Node, toks []lexer.Token) (Node, error) {
	isSelfContained, err := lexer.IsSelfContained(toks)
	if err != nil {
		return elm, err
	}
	if isSelfContained {
		if len(toks) == 1 {
			return elm, nil
		}
		child := NewNodeNoraml(lexer.Construct(toks), Normal)
		AppendChild(elm, child)
	}
	innerToks, err := lexer.ShedOuterHtml(toks)
	if err != nil {
		return elm, err
	}
	for i, tok := range innerToks {
		if tok.GetType() == string(lexer.HtmlVoid) {
			child := NewNodeVoid(tok.GetLexeme(), Void)
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
			child := NewNodeNoraml(lexer.Construct(elmToks), Normal)
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