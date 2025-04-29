package lexer

// Lexer represents a simple rune-based lexer for walking over text input.
type Lexer struct {
	Pos         int            // Current position in the source.
	Source      []rune          // Source runes being lexed.
	Current     rune            // Current rune being analyzed.
	Terminated  bool            // Whether lexer has finished walking through the source.
	Buffer      []rune          // Optional buffer to collect runes manually.
	MarkedPos   int             // Marked position for later collection.
	State       string
	CharCounter map[rune]int    // 🔥 Track rune counts.
	Line        int             // 🔥 Current line number (starts at 1).
	Column      int             // 🔥 Current column number (starts at 1, relative to last newline).
}


// NewLexer creates and initializes a new Lexer from the given rune slice.
func NewLexer(runes []rune) *Lexer {
	l := &Lexer{}
	l.Pos = 0
	l.Source = runes
	l.Terminated = false
	l.MarkedPos = 0
	l.CharCounter = make(map[rune]int)
	l.Line = 1   // 🔥 Start line at 1
	l.Column = 1 // 🔥 Start column at 1
	if len(runes) == 0 {
		l.Current = 0
		l.Terminated = true
	} else {
		l.Current = l.Source[0]
	}
	return l
}



// Step advances the lexer to the next rune in the source,
// updating line and column tracking.
func (l *Lexer) Step() {
	if l.Terminated {
		return
	}
	if l.Pos+1 > len(l.Source)-1 {
		l.Terminated = true
		return
	}
	l.Pos++
	l.Current = l.Source[l.Pos]

	// 🔥 Update line/column numbers:
	if l.Current == '\n' {
		l.Line++
		l.Column = 1
	} else {
		l.Column++
	}
}


// SkipWhiteSpace skips over spaces, tabs, and newlines.
func (l *Lexer) SkipWhiteSpace() {
	for {
		if l.Terminated {
			break
		}
		if l.Current == ' ' || l.Current == '\t' || l.Current == '\n' {
			l.Step()
			continue
		}
		break
	}
}

// Mark saves the current lexer position into MarkedPos.
// Used to later collect a slice of runes from that position.
func (l *Lexer) Mark() {
	l.MarkedPos = l.Pos
}

// CollectFromMark collects and returns all runes from the last MarkedPos
// up to and including the current position.
func (l *Lexer) CollectFromMark() []rune {
	if l.MarkedPos < 0 || l.MarkedPos > len(l.Source) {
		return nil
	}
	if l.MarkedPos > l.Pos {
		return nil // Prevent invalid slicing if MarkedPos somehow moved past Pos
	}
	return l.Source[l.MarkedPos : l.Pos+1]
}

// Push appends the current rune into the Buffer.
func (l *Lexer) Push() {
	l.Buffer = append(l.Buffer, l.Current)
}

// WalkToEnd moves the lexer position to the end of the source
// and marks the lexer as terminated.
func (l *Lexer) WalkToEnd() {
	for {
		if l.Terminated {
			break
		}
		l.Step()
	}
}

// StepBack moves the lexer one rune backward,
// updating line and column tracking accordingly.
func (l *Lexer) StepBack() {
	if l.Pos <= 0 {
		return
	}
	l.Pos--
	l.Current = l.Source[l.Pos]
	l.Terminated = false

	// 🔥 Update line/column numbers when stepping back:
	if l.Current == '\n' {
		l.Line--
		// Recompute column (walk backward to find last '\n' or start)
		l.Column = 1
		for i := l.Pos - 1; i >= 0; i-- {
			if l.Source[i] == '\n' {
				break
			}
			l.Column++
		}
	} else {
		l.Column--
		if l.Column < 1 {
			l.Column = 1
		}
	}
}


// JumpToMark repositions the lexer back to the last marked position.
func (l *Lexer) JumpToMark() {
	if l.MarkedPos >= 0 && l.MarkedPos < len(l.Source) {
		l.Pos = l.MarkedPos
		l.Current = l.Source[l.Pos]
		l.Terminated = false
	}
}

// WalkUntil walks forward until the target rune is found.
// If found, returns true. If not found, lexer just terminates naturally and returns false.
// Does NOT touch the mark.
func (l *Lexer) WalkUntil(target rune) bool {
	for {
		if l.Terminated {
			return false
		}
		if l.Current == target {
			return true
		}
		l.Step()
	}
}

// WalkBackUntil walks backward until the target rune is found.
// If found, returns true. If not found, lexer just stops at start and returns false.
// Does NOT touch the mark.
func (l *Lexer) WalkBackUntil(target rune) bool {
	for {
		if l.Pos <= 0 {
			return false
		}
		if l.Current == target {
			return true
		}
		l.StepBack()
	}
}

func (l *Lexer) Char() string {
	return string(l.Current)
}

// Flush returns the entire Buffer and clears it.
func (l *Lexer) Flush() []rune {
	out := l.Buffer
	l.Buffer = []rune{}
	return out
}

// FlushFromMark collects runes from MarkedPos up to and including the current Pos,
// clears the Buffer, and returns the collected slice.
func (l *Lexer) FlushFromMark() []rune {
	collected := l.CollectFromMark()
	l.Buffer = []rune{}
	return collected
}


func (l *Lexer) CharIs(char string) (bool) {
	return l.Char() == char	
}

// Peek returns the rune that is offset spaces away from the current position.
// Positive offset looks forward, negative offset looks backward.
// If the calculated position is out of bounds, returns 0.
func (l *Lexer) Peek(offset int) rune {
	targetPos := l.Pos + offset
	if targetPos < 0 || targetPos >= len(l.Source) {
		return 0 // You could choose another sentinel if you want
	}
	return l.Source[targetPos]
}


// WalkUntilSkipQuotes walks forward until it finds the target rune,
// skipping over any instances of the target that occur inside single ('') or double ("") quotes.
// Returns true if the target was found outside of quotes, false otherwise.
func (l *Lexer) WalkUntilSkipQuotes(target rune) bool {
	inSingleQuote := false
	inDoubleQuote := false

	for {
		if l.Terminated {
			return false
		}

		// If outside any quotes and we find the target, success
		if !inSingleQuote && !inDoubleQuote && l.Current == target {
			return true
		}

		// Handle entering or exiting quotes
		if l.Current == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		} else if l.Current == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		}

		l.Step()
	}
}

// Count records the current rune into the CharCounter map.
func (l *Lexer) Count() {
	l.CharCounter[l.Current]++
}


// GetCount returns the number of times the given rune has appeared so far.
func (l *Lexer) GetCount(r rune) int {
	return l.CharCounter[r]
}

// ResetCount clears the CharCounter map, resetting all recorded rune counts.
func (l *Lexer) ResetCount() {
	l.CharCounter = make(map[rune]int)
}

// IsEscaped checks if the current character is escaped by a backslash.
// It looks at the previous character to see if it's a backslash.
// Note: This doesn't account for double-escaped characters.
func (l *Lexer) IsEscaped() bool {
	// If we're at the beginning of the source, the character can't be escaped
	if l.Pos <= 0 {
			return false
	}
	
	// Check if the previous character is a backslash
	return l.Source[l.Pos-1] == '\\'
}