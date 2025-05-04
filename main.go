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

	count := 0
	err = parser.Walk(ast, func(n parser.Node) error {
		count+=1
		return nil
	})
	fmt.Println(count)
	if err != nil {
		panic(err)
	}

	

}
