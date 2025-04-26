package token

type TokenType string

const (

	// deconstructable tokens
	TokenHtmlOpenTag TokenType = "TokenHtmlOpenTag"
	TokenHtmlCloseTag TokenType = "TokenHtmlCloseTag"
	TokenHtmlSelfClosingTag TokenType = "TokenHtmlSelfClosingTag"
	TokenHtmlText TokenType = "TokenHtmlText"
	TokenHtmlWhiteSpace TokenType = "TokenHtmlWhiteSpace"
	TokenHtmlAttribute TokenType = "TokenHtmlAttrbute"
	TokenEndOfFile TokenType = "TokenEndOfFile"

	// smallest-level tokens
	TokenHtmlOpenTagName TokenType = "TokenHtmlOpenTagName"
	TokenHtmlOpenTagOpeningBracket TokenType = "TokenHtmlOpenTagOpeningBracket"
	TokenHtmlOpenTagClosingBracket TokenType = "TokenHtmlOpenTagClosingBracket"
	TokenHtmlBooleanAttribute TokenType = "TokenHtmlBooleanAttribute"
	

)
