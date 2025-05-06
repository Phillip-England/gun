package lexer

import (
	"testing"
)

func TestHtmlToken(t *testing.T) {
	form := []rune(`
		<form>
			<h1>Login Form</h1>
			<ul>
				<li>
					<input type='text' name='username'>			
				</li>
				<li>
					<input type='password' name='password'>
				</li>
			</ul>
		</form>
	`)
	toks, err := TokenizeHtml(form)
	if err != nil {
		panic(err)
	}
	// 13 tokens, right?
	if len(toks) != 13 {
		t.Errorf(`expected 13 toks but found %d`, len(toks))
	}
	// if we iterate through each tag name and collect, we should get something like this:
	expectedTagNames := []string{
		"form",    // <form>
		"h1",      // <h1>
		"",        // Login Form (text)
		"h1",      // </h1>
		"ul",      // <ul>
		"li",      // <li>
		"input",   // <input type='text' name='username'>
		"li",      // </li>
		"li",      // <li>
		"input",   // <input type='password' name='password'>
		"li",      // </li>
		"ul",      // </ul>
		"form",    // </form>
	}
	// then if we collect the actual tag names
	actualTagNames := []string{}
	for i, tok := range toks {
		actualTagNames = append(actualTagNames, GetTagName(tok))
		currentName := GetTagName(tok)
		expectedName := expectedTagNames[i]
		if currentName != expectedName {
			t.Errorf(`all tag names should match but here they dont expected tag name: %s actual tag name: %s`, currentName, expectedName)
		}
	}
	if len(expectedTagNames) != len(actualTagNames) {
		t.Errorf(`lens do not match and they should len(expectedTagNames) %d len(actualTagNames) %d`, len(expectedTagNames), len(actualTagNames))
	}
	// a slice of tokens is considered to be self contained it if does not contain any outlying stragling html bits
	
}