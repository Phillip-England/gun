package lexer


type Lexer struct {
	Pos int
	String string
	Char string
	Terminated bool
}

func NewLexer(runes []rune) (*Lexer) {
	s := string(runes)
	l := &Lexer{}
	l.Pos = 0
	l.String = s
	l.Terminated = false
	if len(s) == 0 {
		l.Char = ""
		l.Terminated = true
	} else {
		l.Char = string(s[0])
	}
	return l
}

func (l *Lexer) Step() {
	if l.Terminated {
		return
	}
	if l.Pos + 1 > len(l.String)-1 {
		l.Terminated = true
		return
	}
	l.Pos += 1
	l.Char = string(l.String[l.Pos])
}

func (l *Lexer) SkipWhiteSpace() {
	for {
		if l.Terminated {
			break
		}
		if l.Char == " " || l.Char == "\t" || l.Char == "\n" {
			l.Step()
			continue
		}
		break
	}
}