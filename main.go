package main

import (
	"github.com/phillip-england/gun/logi"
	"github.com/phillip-england/gun/token"
)

func main() {

	logi.Clear()

	toks, err := token.Tokenize(`
		<h1 name='%s name%' age='25%'>
			<h1>hello, %s name%!</h1>
		</h1>
		<p>wtf is up! %s age%</p>
		<input type='text'>
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

	logi.Log(token.Construct(toks))


}