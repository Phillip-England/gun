package lexer

// Lexer represents a simple rune-based lexer for walking over text input.
type Lexer struct {
	Pos         int          // Current position in the source.
	Source      []rune       // Source runes being lexed.
	Current     rune         // Current rune being analyzed.
	Terminated  bool         // Whether lexer has finished walking through the source.
	Buffer      []rune       // Optional buffer to collect runes manually.
	MarkedPos   int          // Marked position for later collection.
	State       string
	CharCounter map[rune]int // Track rune counts.
	Spent []rune
	Line int
	Column int
}

// NewLexer creates and initializes a new Lexer from the given rune slice.
func NewLexer(runes []rune) *Lexer {
	l := &Lexer{
		Pos:         0,
		Source:      runes,
		MarkedPos:   0,
		CharCounter: make(map[rune]int),
		Spent: []rune{},
	}
	if len(runes) == 0 {
		l.Current = 0
		l.Terminated = true
	} else {
		l.Current = l.Source[0]
	}
	return l
}

// Step advances the lexer to the next rune.
func (l *Lexer) Step() {
	if l.Pos < 0 {
		l.Pos = 0
	}
	if l.Terminated {
		l.Pos = len(l.Source)-1
		return
	}
	if l.Pos+1 >= len(l.Source) {
		l.Terminated = true
		l.Pos = len(l.Source)-1
		return
	}
	l.Pos++
	l.Current = l.Source[l.Pos]
	l.updateSpent()
	l.updateLineAndColumn()
}


// StepBack moves the lexer one rune backward.
func (l *Lexer) StepBack() {
	if l.Pos <= 0 {
		l.Pos = 0
		return
	}
	if l.Pos > len(l.Source)-1 {
		l.Pos = len(l.Source)-1
	}
	l.Pos--
	l.Current = l.Source[l.Pos]
	l.Terminated = false
	l.updateSpent()
	l.updateLineAndColumn()
}


// Mark saves the current position.
func (l *Lexer) Mark() {
	l.MarkedPos = l.Pos
}

// JumpToMark repositions the lexer back to the last marked position.
func (l *Lexer) JumpToMark() {
	if l.MarkedPos >= 0 && l.MarkedPos < len(l.Source) {
		l.Pos = l.MarkedPos
		l.Current = l.Source[l.Pos]
		l.Terminated = false
		l.updateSpent()
		l.updateLineAndColumn()
	}
}


// CollectFromMark returns all runes from MarkedPos up to current Pos.
func (l *Lexer) CollectFromMark() []rune {
	if l.MarkedPos < 0  || l.Pos >= len(l.Source) {
		return nil
	}
	if l.Pos+1 > l.MarkedPos {
		return l.Source[l.MarkedPos : l.Pos+1]
	} else {
		return l.Source[l.Pos : l.MarkedPos+1]
	}
}

// Push adds the current rune to the buffer.
func (l *Lexer) Push() {
	l.Buffer = append(l.Buffer, l.Current)
}

// Flush clears and returns the buffer.
func (l *Lexer) Flush() []rune {
	out := l.Buffer
	l.Buffer = []rune{}
	return out
}

// FlushFromMark collects runes from mark and clears the buffer.
func (l *Lexer) FlushFromMark() []rune {
	collected := l.CollectFromMark()
	l.Buffer = []rune{}
	return collected
}

// WalkToEnd steps until termination.
func (l *Lexer) WalkToEnd() {
	for !l.Terminated {
		l.Step()
	}
}

// WalkUntil stops when target rune is found.
func (l *Lexer) WalkUntil(target rune) bool {
	for !l.Terminated {
		l.Step()
		if l.Current == target {
			return true
		}
	}
	return false
}

// WalkBackUntil steps back until target rune is found.
func (l *Lexer) WalkBackUntil(target rune) bool {
	for l.Pos > 0 {
		if l.Current == target {
			return true
		}
		l.StepBack()
	}
	return false
}

// Char returns the current rune as string.
func (l *Lexer) Char() string {
	return string(l.Current)
}

// CharIs checks if current rune matches.
func (l *Lexer) CharIs(char string) bool {
	return l.Char() == char
}

// Peek looks ahead or behind without moving.
func (l *Lexer) Peek(offset int) rune {
	targetPos := l.Pos + offset
	if targetPos < 0 || targetPos >= len(l.Source) {
		return 0
	}
	return l.Source[targetPos]
}

// SkipWhiteSpace steps through space, tab, newline.
func (l *Lexer) SkipWhiteSpace() {
	for !l.Terminated && (l.Current == ' ' || l.Current == '\t' || l.Current == '\n') {
		l.Step()
	}
}

// WalkUntilSkipQuotes skips quoted targets.
func (l *Lexer) WalkUntilSkipQuotes(target rune) bool {
	inSingleQuote := false
	inDoubleQuote := false

	for !l.Terminated {
		if !inSingleQuote && !inDoubleQuote && l.Current == target {
			return true
		}
		if l.Current == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		} else if l.Current == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		}
		l.Step()
	}
	return false
}

// Count tracks current rune.
func (l *Lexer) Count() {
	l.CharCounter[l.Current]++
}

// GetCount returns count of a rune.
func (l *Lexer) GetCount(r rune) int {
	return l.CharCounter[r]
}

// ResetCount clears all rune counts.
func (l *Lexer) ResetCount() {
	l.CharCounter = make(map[rune]int)
}

// IsEscaped checks if current rune is escaped.
func (l *Lexer) IsEscaped() bool {
	return l.Pos > 0 && l.Source[l.Pos-1] == '\\'
}


// updateSpent recalculates the runes from start to current Pos.
func (l *Lexer) updateSpent() {
	if l.Pos == len(l.Source)-1 {
		l.Spent = l.Source
		return
	}
	if l.Pos >= 0 && l.Pos < len(l.Source) {
		l.Spent = l.Source[0 : l.Pos]
	} else if l.Pos >= len(l.Source) {
		l.Spent = l.Source
	} else {
		l.Spent = []rune{}
	}
}

func (l *Lexer) updateLineAndColumn() {
	l.Line = 1
	l.Column = 1
	for i := 0; i < l.Pos; i++ {
		if l.Source[i] == '\n' {
			l.Line++
			l.Column = 1
		} else {
			l.Column++
		}
	}
}

// SpentString returns the spent runes as a string.
func (l *Lexer) SpentString() string {
	return string(l.Spent)
}

// StepBackToStart resets the lexer to the start of the source.
func (l *Lexer) WalkBackToStart() {
	l.Pos = 0
	l.Terminated = len(l.Source) == 0
	if !l.Terminated {
		l.Current = l.Source[0]
	}
	l.updateSpent()
	l.updateLineAndColumn()
}


