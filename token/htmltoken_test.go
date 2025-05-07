package token

import (
	"fmt"
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
		"form",  // <form>
		"h1",    // <h1>
		"",      // Login Form (text)
		"h1",    // </h1>
		"ul",    // <ul>
		"li",    // <li>
		"input", // <input type='text' name='username'>
		"li",    // </li>
		"li",    // <li>
		"input", // <input type='password' name='password'>
		"li",    // </li>
		"ul",    // </ul>
		"form",  // </form>
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
	isSelfContained, err := IsSelfContained(toks)
	if err != nil {
		panic(err)
	}
	if !isSelfContained {
		t.Errorf(`the login form should be a self-contained element but it is not`)
	}
	// lets lose our outer html and see what our first tag is
	inner, err := ShedOuterHtml(toks)
	if err != nil {
		panic(err)
	}
	tagname := GetTagName(inner[0])
	if tagname != "h1" {
		t.Errorf(`expected first token to be an <h1> but instead it was <%s>`, GetTagName(inner[0]))
	}
	// since the <h1> is not self-contained (as it has free-floating siblings)
	// ShedOuterHtml should not modify it
	preShedStr := Construct(inner)
	inner1, err := ShedOuterHtml(inner)
	if err != nil {
		panic(err)
	}
	postShedStr := Construct(inner1)
	if preShedStr != postShedStr {
		t.Errorf(`expected shedding to not alter the <h1> in the login form but it did`)
	}
	// lets see if we can locate the closing tag
	h1 := inner1[0]
	if GetTagName(h1) != "h1" {
		t.Errorf(`expected first tag in inner1 to be a <h1>`)
	}
	closingh1, _, err := GetClosingTag(h1, 0, inner1)
	if err != nil {
		panic(err)
	}
	if closingh1.GetLexeme() != "</h1>" {
		t.Errorf(`expected </h1> but found %s`, closingh1.GetLexeme())
	}
	// and can we just extract the element outright?
	elmStr, err := ExtractFullElement(h1, 0, inner1)
	if err != nil {
		panic(err)
	}
	if elmStr != "<h1>Login Form</h1>" {
		t.Errorf(`expected ExtractFullElement to pull the <h1> out but instead we got %s`, elmStr)
	}
	// okay what if we start messing with our input, how do things work?
	empty := []rune("")
	toks, err = TokenizeHtml(empty)
	if err != nil {
		panic(err)
	}
	fmt.Println(toks)

}
