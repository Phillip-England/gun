package token

type TokenType string

const (

	// stage 1 tokens
	TokenHtmlOpenTag TokenType = "TokenHtmlOpenTag"
	TokenHtmlCloseTag TokenType = "TokenHtmlCloseTag"
	TokenHtmlSelfClosingTag TokenType = "TokenHtmlSelfClosingTag"
	TokenHtmlTextNode TokenType = "TokenHtmlTextNode"
	TokenHtmlTextWhiteSpace TokenType = "TokenHtmlTextWhiteSpace"
	TokenEndOfInput TokenType = "TokenEndOfInput"
	
	// stage 2 tokens
	TokenHtmlOpenTagOpeningBracket TokenType = "TokenHtmlOpenTagOpeningBracket"
	TokenHtmlOpenTagClosingBracket TokenType = "TokenHtmlOpenTagClosingBracket"
	TokenHtmlOpenTagName TokenType = "TokenHtmlOpenTagName"
	TokenHtmlAttribute TokenType = "TokenHtmlAttrbute"
	TokenHtmlBooleanAttribute TokenType = "TokenHtmlBooleanAttribute"
	TokenHtmlRawText TokenType = "TokenHtmlRawText"
	TokenHtmlTextArg TokenType = "TokenHtmlTextArg"
	TokenHtmlCloseTagOpeningBracket TokenType = "TokenHtmlCloseTagOpeningBracket"
	TokenHtmlCloseTagClosingBracket TokenType = "TokenHtmlCloseTagClosingBracket"
	TokenHtmlCloseTagName TokenType = "TokenHtmlCloseTagName"
	TokenHtmlSelfClosingTagOpeningBracket TokenType = "TokenHtmlSelfClosingTagOpeningBracket"
	TokenHtmlSelfClosingTagClosingBracket TokenType = "TokenHtmlSelfClosingTagClosingBracket"
	TokenHtmlSelfClosingTagName TokenType = "TokenHtmlSelfClosingTagName"
	TokenHtmlTagWhiteSpace TokenType = "TokenHtmlTagSpace"

	// stage 3 tokens
	TokenHtmlAttributeName TokenType = "TokenHtmlAttributeName"
	TokenHtmlAttributeEqualSign TokenType = "TokenHtmlAttributeEqualSign"
	TokenHtmlAttributeValue TokenType = "TokenHtmlAttributeValue"
	TokenHtmlTextStringArg TokenType = "TokenHtmlTextStringArg"
	TokenHtmlTextIntArg TokenType = "TokenHtmlTextIntArg"
	TokenHtmlTextBoolArg TokenType = "TokenHtmlTextBooleanArg"
	TokenHtmlElseSymbol TokenType = "TokenHtmlElseSymbol"

	// stage 4 tokens
	TokenHtmlAttributeValuePart TokenType = "TokenHtmlAttributeValuePart"
	TokenHtmlAttributeStringArg TokenType = "TokenHtmlAttributeStringArg"
	TokenHtmlAttributeBoolArg TokenType = "TokenHtmlAttributeBoolArg"
	TokenHtmlAttributeIntArg TokenType = "TokenHtmlAttributeIntArg"



)
