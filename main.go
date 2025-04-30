package main

import (
	"fmt"

	"github.com/phillip-england/gtml/lexer"
	"github.com/phillip-england/gtml/logi"
	"github.com/phillip-england/gtml/parser"
)

func main() {

	logi.Clear()

	toks, err := lexer.TokenizeHtml([]rune(`
	
	
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

	fmt.Println(len(ast.GetInfo().Children))
	parser.WalkNodes(ast, func(i int, n parser.Node) error {
		logi.Log(n.GetInfo().Value)
		return nil
	})



}
