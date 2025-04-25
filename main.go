package main

import (
	"github.com/phillip-england/gun/lexer"
	"github.com/phillip-england/gun/logi"
)

func main() {

	nodes, err := lexer.Tokenize(`
		<h1>%s name%</h1>
		<p>%s age%</p>
		<ul _for="friend in user.Friend Friend[]">
			<li>
				<p>%s friend.Name%</p>
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

	for _, n := range nodes {
		logi.Log(n.GetInfo().String)
	}

}
