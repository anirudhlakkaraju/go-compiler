package lexer

import (
	"go-compiler/src/monkey/token"
)

// Lexer is synonymous with Tokenizer. Given a string input it reads the tokens.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (points to NEXT char after current)
	ch           byte // current char under examination
}

// Returns Lexer for input string. This Lexer can read the input string's tokens
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Initialize ch, position and readPosition
	return l
}

// Reads next char of input string
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// Returns the Token Type and Literal of the char ch under examination
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			tok.Type = token.COMMENT
			tok.Literal = l.readComment()
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok // Returning early since ch is advanced in l.readIdentifier()
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok // Returning early since ch is advanced in l.readNumber()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// Returns new token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Returns identifier string starting from l.position
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Returns true if ch is Letter, included '_' to accomodate identifiers like 'foo_bar'
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Moves position to skip whitespace
func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// Returns true if ch is Digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Returns number string from l.position
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Returns char in next position
func (l *Lexer) peekChar() byte {
	if l.readPosition < len(l.input) {
		return l.input[l.readPosition]
	} else {
		return 0
	}
}

// readString returns the characters within quotes
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()

		// TODO: Add support for char escaping

		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// readComment returns the comment text
func (l *Lexer) readComment() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\n' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
