package token

type TokenType string

// A Token is made up of a TokenType and a Literal which is the actual value of that token
type Token struct {
	Type    TokenType
	Literal string
}

// Constant variables to define Keywords and Operaters of Monkey Language
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // Identifier string for - add, foo, bar, x, y, ...
	INT   = "INT"   // 12345

	// Operaters
	ASSIGN   = "="
	EQ       = "=="
	NOT_EQ   = "!="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	STRING  = "STRING"
	COMMENT = "COMMENT"
)

// Map to store language specific keywords
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// Returns TokenType given ident string - keyword if present in map else IDENT to indicate user-defined identifier
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
