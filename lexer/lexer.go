package lexer

import "strings"

type Lexer struct {
	Source  []byte
	Current byte
	Pos     int
	Buffer  []byte
	Done    bool
	Mark    int
}

// NewLexer creates a new Lexer instance from the given source string.
func NewLexer(source string) *Lexer {
	l := &Lexer{
		Source: []byte(source),
		Pos:    0,
		Buffer: []byte{},
		Done:   false,
		Mark:   0,
	}
	if len(source) > 0 {
		l.Current = l.Source[0]
	} else {
		l.Current = 0
		l.Done = true
	}
	return l
}

func (l *Lexer) Step() {
	l.Pos++
	if l.Pos > len(l.Source)-1 {
		l.Done = true
		return
	}
	l.Current = l.Source[l.Pos]
}

func (l *Lexer) WalkTo(target byte) {
	for {
		if l.Done {
			return
		}
		if l.Current == target {
			return
		}
		l.Step()
	}
}

func (l *Lexer) Char() byte {
	return l.Current
}

func (l *Lexer) Push() {
	l.Buffer = append(l.Buffer, l.Current)
}

func (l *Lexer) Grow(s string) {
	l.Pos += len(s)
	if l.Pos >= len(l.Source) {
		l.Pos = len(l.Source) - 1
		l.Current = 0
		l.Done = true
		return
	}
	l.Current = l.Source[l.Pos]
	l.Done = false
}

func (l *Lexer) MarkPos() {
	l.Mark = l.Pos
}

func (l *Lexer) ClearMark() {
	l.Mark = 0
}

func (l *Lexer) CollectFromMark() {
	start := l.Mark
	end := l.Pos
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if end >= len(l.Source) {
		end = len(l.Source) - 1
	}
	substr := l.Source[start : end+1]
	l.Buffer = append(l.Buffer, substr...)
}

func (l *Lexer) Rewind() {
	l.Pos = l.Mark
	l.Mark = 0
	if l.Pos >= 0 && l.Pos < len(l.Source) {
		l.Current = l.Source[l.Pos]
	} else {
		l.Current = 0
		l.Done = true
	}
}

func (l *Lexer) SkipWhitespace() {
	for {
		if l.Done {
			return
		}
		if l.Char() != ' ' && l.Char() != '\t' && l.Char() != '\n' {
			return
		}
		l.Step()
	}
}

func (l *Lexer) Peek(by int, asSubstring bool) string {
	if len(l.Source) == 0 {
		return ""
	}
	target := l.Pos + by
	if target < 0 {
		target = 0
	}
	if target >= len(l.Source) {
		target = len(l.Source) - 1
	}
	if asSubstring {
		start := l.Pos
		end := target
		if start > end {
			start, end = end, start
		}
		if end >= len(l.Source) {
			end = len(l.Source) - 1
		}
		return string(l.Source[start : end+1])
	}
	return string(l.Source[target])
}

func (l *Lexer) FlushBuffer() string {
	result := string(l.Buffer)
	l.Buffer = []byte{}
	return result
}

func (l *Lexer) StepBack() {
	if l.Pos <= 0 {
		l.Pos = 0
		l.Current = 0
		l.Done = true
		return
	}
	l.Pos--
	l.Current = l.Source[l.Pos]
	l.Done = false
}

func (l *Lexer) WalkBackTo(target byte) {
	for {
		if l.Pos <= 0 {
			l.Pos = 0
			l.Current = 0
			l.Done = true
			return
		}
		if l.Current == target {
			return
		}
		l.StepBack()
	}
}

func (l *Lexer) WalkToWithQuoteSkip(target byte) {
	inQuote := false
	var quoteChar byte

	for {
		if l.Done {
			return
		}
		if (l.Char() == '"' || l.Char() == '\'') && l.Peek(-1, false) != `\` {
			if !inQuote {
				inQuote = true
				quoteChar = l.Char()
			} else if l.Char() == quoteChar {
				inQuote = false
				quoteChar = 0
			}
		}
		if l.Char() == target && !inQuote {
			return
		}
		l.Step()
	}
}

func (l *Lexer) FlushSplitWithStringPreserve(delim string) []string {
	text := l.FlushBuffer()
	var parts []string
	var b strings.Builder

	inQuote := false
	var quoteChar rune
	i := 0
	for i < len(text) {
		ch := rune(text[i])
		if (ch == '"' || ch == '\'') && (i == 0 || rune(text[i-1]) != '\\') {
			if !inQuote {
				inQuote = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuote = false
				quoteChar = 0
			}
		}
		if !inQuote && strings.HasPrefix(text[i:], delim) {
			parts = append(parts, b.String())
			b.Reset()
			i += len(delim)
			continue
		}
		b.WriteByte(text[i])
		i++
	}
	if b.Len() > 0 {
		parts = append(parts, b.String())
	}
	return parts
}

func (l *Lexer) WalkToUnescaped(target byte) {
	for {
		if l.Done {
			return
		}
		if l.Current == target && l.Peek(-1, false) != `\` {
			return
		}
		l.Step()
	}
}
