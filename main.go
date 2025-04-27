package main

import (
	"github.com/phillip-england/gtml/logi"
	"github.com/phillip-england/gtml/token"
)

func main() {

	logi.Clear()

	_, err := token.Tokenize([]rune(`
	
		<div>

		<div>
	

	`))
	if err != nil {
	panic(err)
	}



}







