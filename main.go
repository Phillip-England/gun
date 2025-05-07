package main

import (
	"fmt"

	"github.com/phillip-england/gtml/logi"
	"github.com/phillip-england/gtml/parser"
	"github.com/phillip-england/gtml/token"
)

func main() {

	logi.Clear()

	toks, err := token.TokenizeHtml([]rune(`
	
	
		<input type='text'>

		<div>
			<h1>Hello, %s name%!</h1>
			<p>I am %s age% years old. How old are you?</p>
		</div>


	`))
	if err != nil {
		panic(err)
	}



	ast, err := parser.NewAst(toks)
	if err != nil {
		panic(err)
	}


	for _, tok := range toks {
		fmt.Println(tok.GetLexeme(), tok.GetLine(), tok.GetColumn())
	}

	err = parser.Walk(ast, func(n parser.Node) error {
		logi.Log(n.GetInfo().Value)
		logi.Log(n.GetInfo().TextContent)
		logi.Log(parser.GetAttribute(n, "type"))
		return nil
	})
	if err != nil {
		panic(err)
	}

	

}
