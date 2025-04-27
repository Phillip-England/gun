package lexer

import "strings"

type Lexer struct {
	Source  string
	Current string
	Pos     int
	Buffer  []string
	Done    bool
	Mark    int
}

// NewLexer creates a new Lexer instance from the given source string.
func NewLexer(source string) *Lexer {
	l := &Lexer{}
	l.Source = source
	l.Pos = 0
	l.Buffer = []string{}
	l.Done = false
	l.Mark = 0
	if len(source) > 0 {
		l.Current = string(source[0])
	} else {
		l.Current = ""
		l.Done = true
	}
	return l
}

// Step moves the cursor forward by one character.
func (l *Lexer) Step() {
	l.Pos += 1
	if l.Pos > len(l.Source)-1 {
		l.Done = true
		return
	}
	ch := string(l.Source[l.Pos])
	l.Current = ch
}

// WalkTo steps forward until the current character matches the target character.
func (l *Lexer) WalkTo(target string) {
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

// Char returns the current character under the cursor.
func (l *Lexer) Char() string {
	return l.Current
}

// Push adds the current character to the buffer if it's not empty.
func (l *Lexer) Push() {
	if l.Current != "" {
		l.Buffer = append(l.Buffer, l.Current)
	}
}

// Grow advances the cursor by the length of the provided string.
func (l *Lexer) Grow(s string) {
	l.Pos += len(s)
	if l.Pos >= len(l.Source) {
		l.Pos = len(l.Source) - 1
		l.Current = ""
		l.Done = true
		return
	}
	l.Current = string(l.Source[l.Pos])
	l.Done = false
}

// MarkPos saves the current cursor position to Mark.
func (l *Lexer) MarkPos() {
	l.Mark = l.Pos
}

// ClearMark resets the Mark back to 0.
func (l *Lexer) ClearMark() {
	l.Mark = 0
}

// CollectFromMark collects all characters from Mark to the current position into the buffer.
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
	for _, ch := range substr {
		l.Buffer = append(l.Buffer, string(ch))
	}
}

// Rewind moves the cursor back to the last marked position.
func (l *Lexer) Rewind() {
	l.Pos = l.Mark
	l.Mark = 0
	if l.Pos >= 0 && l.Pos < len(l.Source) {
		l.Current = string(l.Source[l.Pos])
	} else {
		l.Current = ""
		l.Done = true
	}
}

// SkipWhitespace advances the cursor while it's on whitespace characters (space, tab, newline).
func (l *Lexer) SkipWhitespace() {
	for {
		if l.Done {
			return
		}
		if l.Char() != " " && l.Char() != "\t" && l.Char() != "\n" {
			return
		}
		l.Step()
	}
}

// Peek looks ahead (or behind) by a certain number of characters, optionally returning a substring.
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
		return l.Source[start : end+1]
	}
	return string(l.Source[target])
}

// FlushBuffer returns the contents of the buffer as a string and clears the buffer.
func (l *Lexer) FlushBuffer() string {
	var b strings.Builder
	for _, s := range l.Buffer {
		b.WriteString(s)
	}
	l.Buffer = []string{}
	return b.String()
}

// StepBack moves the cursor backward by one character.
func (l *Lexer) StepBack() {
	if l.Pos <= 0 {
		l.Pos = 0
		l.Current = ""
		l.Done = true
		return
	}
	l.Pos -= 1
	l.Current = string(l.Source[l.Pos])
	l.Done = false
}

// WalkBackTo steps backward until the current character matches the target character.
func (l *Lexer) WalkBackTo(target string) {
	for {
		if l.Pos <= 0 {
			l.Pos = 0
			l.Current = ""
			l.Done = true
			return
		}
		if l.Current == target {
			return
		}
		l.StepBack()
	}
}

// WalkToWithQuoteSkip steps forward until the target character is found outside of quotes.
func (l *Lexer) WalkToWithQuoteSkip(target string) {
	inQuote := false
	quoteChar := ""

	for {
		if l.Done {
			return
		}
		if (l.Char() == `"` || l.Char() == `'`) && l.Peek(-1, false) != `\` {
			if !inQuote {
				inQuote = true
				quoteChar = l.Char()
			} else if l.Char() == quoteChar {
				inQuote = false
				quoteChar = ""
			}
		}
		if l.Char() == target && !inQuote {
			return
		}
		l.Step()
	}
}

// FlushSplitWithStringPreserve flushes the buffer and splits the result
// by the given delimiter, but ignores delimiters inside quotes.
func (l *Lexer) FlushSplitWithStringPreserve(delim string) []string {
	text := l.FlushBuffer()
	var parts []string
	var b strings.Builder

	inQuote := false
	quoteChar := ""
	i := 0
	for i < len(text) {
		ch := string(text[i])
		if (ch == `"` || ch == `'`) && (i == 0 || string(text[i-1]) != `\`) {
			if !inQuote {
				inQuote = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuote = false
				quoteChar = ""
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

// WalkToUnescaped steps forward until the target character is found unescaped (not preceded by a backslash).
func (l *Lexer) WalkToUnescaped(target string) {
	for {
		if l.Done {
			return
		}
		// Check if current char matches and is not escaped
		if l.Current == target && l.Peek(-1, false) != `\` {
			return
		}
		l.Step()
	}
}