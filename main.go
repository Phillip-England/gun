package main

import (
	"github.com/phillip-england/gun/logi"
	"github.com/phillip-england/gun/token"
)

func main() {

	logi.Clear()

	toks, err := token.TokenizeHtml(`
		<h1 name='name' age='25'>%s name%</h1>
		<p>%s age%</p>
		<input type='text'/>
		<ul _for="friend in user.Friend Friend[]">
			<li>
				<p raw>%s friend.Name%</p>
				<div _if="friend.Age > 21">
					<p>you all can drink together</p>
					::?
					<p>you all cannot drink together</p>
				</div>
			</li>
		</ul>
	`)
	if err != nil {
		panic(err)
	}

	toks, err = token.Deconstruct(toks)
	if err != nil {
		panic(err)
	}

	for _, tok := range toks {
		logi.Log(tok.Type)
		logi.Log(tok.Lexeme)


	}

}
