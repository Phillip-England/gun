package main

import (
	"fmt"

	"github.com/phillip-england/gun/lexer"
)

func main() {

	n, err := lexer.NewNodeRoot(`
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

	fmt.Println(n)

}
