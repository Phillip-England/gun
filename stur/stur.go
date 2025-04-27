package stur

import "strings"

func EnforceWhitelist(s string, check ...string) bool {
	for _, allowed := range check {
		if s == allowed {
			return true
		}
	}
	return false
}

func Squeeze(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, "\t", ""), "\n", ""), " ", "")
}

func SplitWithStringPreserve(s, delim string) []string {
	var parts []string
	var b strings.Builder

	inQuote := false
	quoteChar := ""
	i := 0
	for i < len(s) {
		ch := string(s[i])
		if (ch == `"` || ch == `'`) && (i == 0 || string(s[i-1]) != `\`) {
			if !inQuote {
				inQuote = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuote = false
				quoteChar = ""
			}
		}
		if !inQuote && strings.HasPrefix(s[i:], delim) {
			parts = append(parts, b.String())
			b.Reset()
			i += len(delim)
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	if b.Len() > 0 {
		parts = append(parts, b.String())
	}
	return parts
}

func RemoveLastChar(s string) string {
	if len(s) == 0 {
		return s
	}
	return s[:len(s)-1]
}

func LastChar(s string) string {
	if len(s) == 0 {
		return ""
	}
	return string(s[len(s)-1])
}

func StartsWith(s, p string) bool {
	if len(p) > len(s) {
			return false
	}
	return s[:len(p)] == p
}