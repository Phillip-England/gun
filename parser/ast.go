package parser

import (
	"github.com/phillip-england/gtml/token"
)

type Document struct {
	Info *NodeInfo	
}

func (elm *Document) GetInfo() *NodeInfo {
	return elm.Info
}

func NewAst(toks []token.Token) (Node, error) {
	var doc Node
	doc = &Document{
		Info: NewNodeInfo("", Root),
	}
	doc, err := firstPass(doc, toks)
	if err != nil {
		return doc, err
	}
	return doc, nil
}


// FIRST PASS IS RUNNING FOREVER
// NEED TO CREATE GOOD DOM HERE

func firstPass(n Node, toks []token.Token) (Node, error) {
	switch n.GetInfo().Type {
	case Normal:
		innerToks, err := token.ShedOuterHtml(toks)
		if err != nil {
			return n, err
		}
		for i := 0; i < len(innerToks); {
			tok := innerToks[i]
			switch tok.GetType() {
			case token.HtmlOpen:
				_, endTagI, err := token.GetClosingTag(tok, i, innerToks)
				if err != nil {
					return n, err
				}
				child, err := firstPass(
					NewNodeNormal(token.Construct(innerToks[i:endTagI+1]), Normal),
					innerToks[i:endTagI+1],
				)
				if err != nil {
					return n, err
				}
				AppendChild(n, child)
				i = endTagI + 1
				continue
			case token.HtmlVoid:
				child, err := firstPass(NewNodeVoid(tok.GetLexeme(), Void), []token.Token{tok})
				if err != nil {
					return n, err
				}
				AppendChild(n, child)
				i++
				continue
			// case lexer.EmptySpace:
			// 	AppendTextNode(n, tok.GetLexeme())
			// 	i++
			// 	continue
			case token.Text:
				AppendTextNode(n, tok.GetLexeme())
				i++
				continue
			default:
				// Handle other tokens if necessary
				i++
			}
		}



	case Root:
		isSelfContained, err := token.IsSelfContained(toks)
		if err != nil {
			return n, err
		}
		if isSelfContained {
			child, err := firstPass(NewNodeNormal(token.Construct(toks), Normal), toks)
			if err != nil {
				return n, err
			}
			AppendChild(n, child)
			break
		}
		for i := 0; i < len(toks); {
			tok := toks[i]
			switch tok.GetType() {
			case token.HtmlOpen:
				_, endTagI, err := token.GetClosingTag(tok, i, toks)
				if err != nil {
					return n, err
				}
				child, err := firstPass(
					NewNodeNormal(token.Construct(toks[i:endTagI+1]), Normal),
					toks[i:endTagI+1],
				)
				if err != nil {
					return n, err
				}
				AppendChild(n, child)
				i = endTagI + 1
				continue
			case token.HtmlVoid:
				child, err := firstPass(NewNodeVoid(tok.GetLexeme(), Void), []token.Token{tok})
				if err != nil {
					return n, err
				}
				AppendChild(n, child)
				i++
				continue
			// case lexer.EmptySpace:
			// 	AppendTextNode(n, tok.GetLexeme())
			// 	i++
			// 	continue
			case token.Text:
				AppendTextNode(n, tok.GetLexeme())
				i++
				continue
			default:
				// Handle other tokens if necessary
				i++
			}
		}
	}
	return n, nil
}


func Walk(n Node, cb func(Node) error) error {
	if err := cb(n); err != nil {
		return err
	}
	for _, child := range n.GetInfo().Children {
		if err := Walk(child, cb); err != nil {
			return err
		}
	}
	return nil
}