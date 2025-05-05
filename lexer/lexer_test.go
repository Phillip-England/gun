package lexer

import (
	"testing"
)


func TestLexerNavigation(t *testing.T) {
	input := []rune("<p>Hello, World!</p>\n<input type='text'>")
	l := NewLexer(input)
	// does stepping work as expected?
	l.Step()
	if l.Char() != "p" {
		t.Errorf(`expected lexer to be positioned at "p" but instead it is positioned at "%s"`, l.Char())
	}
	// what if we step back?
	l.StepBack()
	if l.Char() != "<" {
		t.Errorf(`expected lexer to be positioned at "<" but instead it is positioned at "%s"`, l.Char())
	}
	// what characters have been spent? (should be none)
	if len(l.Spent) != 0 {
		t.Errorf(`expected lexer to have 0 runes spent, but instead it had %d runes spent`, len(l.Spent))
	}
	// but if we step again we should have a len of 1
	l.Step()
	if len(l.Spent) != 1 {
		t.Errorf(`expected lexer to have 1 rune spent, but instead it had %d runes spent`, len(l.Spent))
	}
	// and the spent string should match
	if l.SpentString() != "<" {
		t.Errorf(`expected the lexers SpendString to equal "<" but it was "%s"`, l.SpentString())
	}
	// and if we step wildly
	l.Step()
	l.Step()
	l.Step()
	l.StepBack()
	l.StepBack()
	l.Step()
	if l.SpentString() != "<p>" {
		t.Errorf(`expected the lexers SpendString to equal "<p>" but it was "%s"`, l.SpentString())
	}
	// and if we collect from mark (the marked pos should be 0)
	// we should get <p>H
	if string(l.CollectFromMark()) != "<p>H" {
		t.Errorf(`expected CollectFromMark to output "<p>H" but it output %s instead`, string(l.CollectFromMark()))
	}
	// now lets mark and step
	l.Mark()
	l.Step()
	if string(l.CollectFromMark()) != "He" {
		t.Errorf(`expected CollectFromMark to output "He" but instead it output %s`, string(l.CollectFromMark()))
	}
	// if we step and collect again we should get Hel
	l.Step()
	if string(l.CollectFromMark()) != "Hel" {
		t.Errorf(`expected CollectFromMark to output "Hel" but instead it output %s`, string(l.CollectFromMark()))
	}
	// jump to the "<" in "</p>"
	l.WalkUntil('<')
	if string(l.CollectFromMark()) != "Hello, World!<" {
		t.Errorf(`expected CollectFromMark to output "Hello, World<" but instead it output %s`, string(l.CollectFromMark()))
	}
	// we should still be on the first line
	if l.Line != 1 {
		t.Errorf(`expected to be on the first line, but we are actually on line %d`, l.Line)
	}
	// but after jumping to the next "<" we should be on the second line
	l.WalkUntil('<')
	if l.Line != 2 {
		t.Errorf(`expected to be on second line but we are actually on line: %d`, l.Line)
	}
	// and when we walk to the end the SpentString should equal the input
	l.WalkToEnd()
	if l.SpentString() != string(input) {
		t.Errorf(`expected the SpentString to be equal to the %s, but it was equal to %s`, string(input), l.SpentString())
	}
	// we should be on line two column twenty
	if l.Line != 2 && l.Column != 20 {
		t.Errorf(`expected to be at line 2 and column 20 but found ourselves at line: %d and column: %d`, l.Line, l.Column)
	}
	// the current char should be ">"
	if l.Char() != ">" {
		t.Errorf(`expected the lexers char to be ">" but it was: %s`, l.Char())
	}
	// go to the start
	l.WalkBackToStart()
	// the position should be 0
	if l.Pos != 0 {
		t.Errorf(`walked to start and expected 0 position but we got %d`, l.Pos)
	}
	// count all the "<" runes
	for {
		if l.Terminated {
			break
		}
		if l.Char() == "<" {
			l.Count()
		}
		l.Step()
	}
	// we should have 3 of them
	if l.GetCount('<') != 3 {
		t.Errorf(`counted all "<" runes and expected 3 but we have %d`, l.GetCount('<'))
	}
	// what if we set the position to 100000?
	l.Pos = 100000000
	// lets try to step
	l.Step()
	// our position should be reset back to the len of our source - 1
	if l.Pos != len(l.Source) -1 {
		t.Errorf(`expected our position to be 39 but it is %d`, l.Pos)
	}
	// we should be on the last char ">"
	if l.Char() != ">" {
		t.Errorf(`explicitly set position past len of input and then requested current char, expect the last char ">" but instead got: %s`, l.Char())
	}
	// after we step back, our Position should get patched and equal 38
	l.StepBack()
	if l.Pos != 38 {
		t.Errorf(`expected Position to be 38 but it is %d`, l.Pos)
	}
	if l.Char() != "'" {
		t.Errorf(`expected Char to be ' but it is %s`, l.Char())
	}
	// and if we set ourselves back to -100000
	l.Pos = -1000000000
	l.Step()
	l.StepBack()
	// we should be on the first char again
	if l.Char() != "<" {
		t.Errorf(`expected char to be "<" but it is %s`, l.Char())
	}
	// and if we go to the end and mark
	l.WalkToEnd()
	l.Mark()
	// then walk back to the start
	l.WalkBackToStart()
	// our flushed buffer should equal the input
	if len(l.FlushFromMark()) != len(input) {
		t.Errorf(`FlushFromMark should be equal to the len of the input but its not input len is %d and FlushFromMark len is %d`, len(input), len(l.FlushFromMark()))
	}
}