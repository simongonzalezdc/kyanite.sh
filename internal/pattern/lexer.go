package pattern

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer tokenizes pattern source code into tokens
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
	line         int
	column       int
	tokens       []Token
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 1,
		tokens: make([]Token, 0),
	}
	l.readChar()
	return l
}

// Reset resets the lexer to the beginning state for re-tokenization
func (l *Lexer) Reset() {
	l.position = 0
	l.readPosition = 0
	l.line = 1
	l.column = 1
	l.readChar()
}

// readChar reads the next character from input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL represents EOF
	} else {
		l.ch, _ = utf8.DecodeRuneInString(l.input[l.readPosition:])
	}
	l.position = l.readPosition
	l.readPosition += utf8.RuneLen(l.ch)
	if l.ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing the position
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return ch
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	pos := Position{
		Line:   l.line,
		Column: l.column,
		Offset: l.position,
	}

	switch l.ch {
	case 0:
		tok = Token{Type: TokenEOF, Position: pos}
	case '+':
		tok = Token{Type: TokenPlus, Value: "+", Position: pos}
		l.readChar()
	case '-':
		tok = Token{Type: TokenMinus, Value: "-", Position: pos}
		l.readChar()
	case '*':
		tok = Token{Type: TokenMultiply, Value: "*", Position: pos}
		l.readChar()
	case '/':
		tok = Token{Type: TokenDivide, Value: "/", Position: pos}
		l.readChar()
	case '%':
		tok = Token{Type: TokenModulo, Value: "%", Position: pos}
		l.readChar()
	case ':':
		tok = Token{Type: TokenColon, Value: ":", Position: pos}
		l.readChar()
	case ',':
		tok = Token{Type: TokenComma, Value: ",", Position: pos}
		l.readChar()
	case '.':
		tok = Token{Type: TokenDot, Value: ".", Position: pos}
		l.readChar()
	case '?':
		tok = Token{Type: TokenQuestion, Value: "?", Position: pos}
		l.readChar()
	case '!':
		// Check for != operator
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			tok = Token{Type: TokenNotEq, Value: "!=", Position: pos}
		} else {
			tok = Token{Type: TokenExclaim, Value: "!", Position: pos}
			l.readChar()
		}
	case '@':
		tok = Token{Type: TokenAt, Value: "@", Position: pos}
		l.readChar()
	case '#':
		tok = Token{Type: TokenHash, Value: "#", Position: pos}
		l.readChar()
	case '$':
		tok = Token{Type: TokenDollar, Value: "$", Position: pos}
		l.readChar()
	case '^':
		tok = Token{Type: TokenCaret, Value: "^", Position: pos}
		l.readChar()
	case '&':
		tok = Token{Type: TokenAmpersand, Value: "&", Position: pos}
		l.readChar()
	case '|':
		tok = Token{Type: TokenPipe, Value: "|", Position: pos}
		l.readChar()
	case '~':
		tok = Token{Type: TokenTilde, Value: "~", Position: pos}
		l.readChar()
	case '=':
		// Check for == operator
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			tok = Token{Type: TokenEquals, Value: "==", Position: pos}
		} else {
			tok = Token{Type: TokenEquals, Value: "=", Position: pos}
			l.readChar()
		}
	case '<':
		// Check for <= operator
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			tok = Token{Type: TokenLessEq, Value: "<=", Position: pos}
		} else if l.peekChar() == '>' {
			// Empty <>, treat as less than
			tok = Token{Type: TokenLess, Value: "<", Position: pos}
			l.readChar()
		} else {
			// Look ahead to see if this looks like a sample or note
			nextPos := l.readPosition
			foundClosing := false
			content := ""

			// Look for a closing > within a reasonable distance
			for i := 0; i < 20 && nextPos+i < len(l.input); i++ {
				if l.input[nextPos+i] == '>' {
					foundClosing = true
					content = l.input[nextPos : nextPos+i]
					break
				}
			}

			// If we found a closing > and the content looks like a sample or note
			if foundClosing && (len(content) > 0) {
				// Check if it looks like a sample (contains colon) or note (starts with a-g)
				if strings.Contains(content, ":") ||
					(len(content) > 0 && content[0] >= 'a' && content[0] <= 'g') ||
					(len(content) > 0 && content[0] >= 'A' && content[0] <= 'G') {
					// Try to read as sample or note
					tok = l.readSampleOrNote()
					// If it didn't work, treat as less than
					if tok.Type == TokenEOF && tok.Value == "" {
						tok = Token{Type: TokenLess, Value: "<", Position: pos}
						l.readChar()
					}
				} else {
					// Doesn't look like a sample/note, treat as less than
					tok = Token{Type: TokenLess, Value: "<", Position: pos}
					l.readChar()
				}
			} else {
				// No closing > found or content is empty, treat as less than
				tok = Token{Type: TokenLess, Value: "<", Position: pos}
				l.readChar()
			}
		}
	case '>':
		// Check for >= operator
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			tok = Token{Type: TokenGreaterEq, Value: ">=", Position: pos}
		} else {
			tok = Token{Type: TokenGreater, Value: ">", Position: pos}
			l.readChar()
		}
	case '(':
		tok = Token{Type: TokenLParen, Value: "(", Position: pos}
		l.readChar()
	case ')':
		tok = Token{Type: TokenRParen, Value: ")", Position: pos}
		l.readChar()
	case '[':
		tok = Token{Type: TokenLBracket, Value: "[", Position: pos}
		l.readChar()
	case ']':
		tok = Token{Type: TokenRBracket, Value: "]", Position: pos}
		l.readChar()
	case '{':
		tok = Token{Type: TokenLBrace, Value: "{", Position: pos}
		l.readChar()
	case '}':
		tok = Token{Type: TokenRBrace, Value: "}", Position: pos}
		l.readChar()
	case '"':
		tok.Type = TokenString
		tok.Value = l.readString()
		tok.Position = pos
	default:
		if isDigit(l.ch) {
			tok.Type = TokenNumber
			tok.Value = l.readNumber()
			tok.Position = pos
		} else if isLetter(l.ch) {
			ident := l.readIdentifier()
			tok.Type = LookupIdent(ident)
			tok.Value = ident
			tok.Position = pos
		} else if l.ch == ';' {
			tok = Token{Type: TokenSemicolon, Value: ";", Position: pos}
			l.readChar()
		} else {
			tok = Token{Type: TokenEOF, Value: "", Position: pos}
			l.readChar()
		}
	}

	return tok
}

// skipWhitespace skips over whitespace characters
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (integer or float)
func (l *Lexer) readNumber() string {
	position := l.position
	hasDot := false

	for isDigit(l.ch) || (l.ch == '.' && !hasDot) {
		if l.ch == '.' {
			hasDot = true
		}
		l.readChar()
	}

	return l.input[position:l.position]
}

// readString reads a string literal
func (l *Lexer) readString() string {
	position := l.position + 1 // skip opening quote
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	value := l.input[position:l.position]
	l.readChar() // skip closing quote
	return value
}

// readSampleOrNote reads a sample or note (enclosed in < >)
func (l *Lexer) readSampleOrNote() Token {
	pos := Position{
		Line:   l.line,
		Column: l.column,
		Offset: l.position,
	}

	l.readChar() // skip '<'
	position := l.position

	for l.ch != '>' && l.ch != 0 {
		l.readChar()
	}

	value := l.input[position:l.position]

	// Check if we found a closing '>'
	if l.ch == '>' {
		l.readChar() // skip '>'

		// Determine if it's a sample or note based on content
		if strings.Contains(value, ":") {
			return Token{Type: TokenSample, Value: value, Position: pos}
		}
		return Token{Type: TokenNote, Value: value, Position: pos}
	}

	// No closing '>', return EOF token to indicate failure
	return Token{Type: TokenEOF, Value: "", Position: pos}
}

// Tokenize tokenizes the entire input and returns all tokens
func (l *Lexer) Tokenize() []Token {
	var tokens []Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}

	return tokens
}

// isDigit checks if a character is a digit
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// isLetter checks if a character is a letter
func isLetter(ch rune) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

// LookupIdent looks up an identifier to see if it's a keyword
func LookupIdent(ident string) TokenType {
	// Check for keywords first
	keywords := map[string]TokenType{
		"pattern":     TokenPattern,
		"sample":      TokenSample,
		"note":        TokenNote,
		"rest":        TokenRest,
		"speed":       TokenSpeed,
		"volume":      TokenVolume,
		"pan":         TokenPan,
		"delay":       TokenDelay,
		"reverb":      TokenReverb,
		"filter":      TokenFilter,
		"loop":        TokenLoop,
		"rand":        TokenRand,
		"irand":       TokenIrand,
		"choose":      TokenChoose,
		"degenerate":  TokenDegenerate,
		"sometimes":   TokenSometimes,
		"often":       TokenOften,
		"rarely":      TokenRarely,
		"never":       TokenNever,
		"superimpose": TokenSuperimpose,
		"chop":        TokenChop,
		"gap":         TokenGap,
		"off":         TokenOff,
		"jux":         TokenJux,
		"stack":       TokenStack,
		"or":          TokenOr,
		"and":         TokenAnd,
		"seq":         TokenSeq,
		"fast":        TokenFast,
		"slow":        TokenSlow,
		"degradeBy":   TokenDegradeBy,
		"gain":        TokenGain,
		"room":        TokenRoom,
		"size":        TokenSize,
		"density":     TokenDensity,
	}

	if tokType, ok := keywords[ident]; ok {
		return tokType
	}

	return TokenIdentifier
}

// Pre-compiled regex patterns for performance
var (
	// Pattern for matching numbers (integers and floats)
	numberPattern = regexp.MustCompile(`^\d+(\.\d+)?`)

	// Pattern for matching identifiers
	identifierPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)

	// Pattern for matching samples (e.g., "kick:1" or "snare:2")
	samplePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+:[0-9]+`)

	// Pattern for matching notes (e.g., "c4", "d#5", "eb3")
	notePattern = regexp.MustCompile(`^[a-gA-G][#b]?\d?`)
)

// FastTokenize provides a regex-based tokenization for better performance
func FastTokenize(input string) []Token {
	var tokens []Token
	line := 1
	column := 1
	offset := 0

	for offset < len(input) {
		// Skip whitespace
		if input[offset] == ' ' || input[offset] == '\t' {
			column++
			offset++
			continue
		} else if input[offset] == '\n' {
			line++
			column = 1
			offset++
			continue
		}

		pos := Position{Line: line, Column: column, Offset: offset}
		remaining := input[offset:]

		// Try to match patterns in order of specificity
		if remaining[0] == '<' {
			// Handle <sample:1> or <note> syntax
			endIndex := strings.Index(remaining, ">")
			if endIndex != -1 {
				content := remaining[1:endIndex]
				var tokType TokenType
				if strings.Contains(content, ":") {
					tokType = TokenSample
				} else {
					tokType = TokenNote
				}
				tokens = append(tokens, Token{Type: tokType, Value: content, Position: pos})
				offset += endIndex + 1
				column += endIndex + 1
				continue
			}
		}

		if match := numberPattern.FindString(remaining); match != "" {
			tokens = append(tokens, Token{Type: TokenNumber, Value: match, Position: pos})
			offset += len(match)
			column += len(match)
		} else if match := identifierPattern.FindString(remaining); match != "" {
			tokType := LookupIdent(match)
			tokens = append(tokens, Token{Type: tokType, Value: match, Position: pos})
			offset += len(match)
			column += len(match)
		} else {
			// Check for multi-character operators first
			if offset+1 < len(input) {
				twoChar := input[offset : offset+2]
				switch twoChar {
				case "==":
					tokens = append(tokens, Token{Type: TokenEquals, Value: "==", Position: pos})
					offset += 2
					column += 2
					continue
				case "!=":
					tokens = append(tokens, Token{Type: TokenNotEq, Value: "!=", Position: pos})
					offset += 2
					column += 2
					continue
				case "<=":
					tokens = append(tokens, Token{Type: TokenLessEq, Value: "<=", Position: pos})
					offset += 2
					column += 2
					continue
				case ">=":
					tokens = append(tokens, Token{Type: TokenGreaterEq, Value: ">=", Position: pos})
					offset += 2
					column += 2
					continue
				}
			}

			// Single character tokens
			ch := input[offset]
			var tokType TokenType
			switch ch {
			case '+':
				tokType = TokenPlus
			case '-':
				tokType = TokenMinus
			case '*':
				tokType = TokenMultiply
			case '/':
				tokType = TokenDivide
			case '%':
				tokType = TokenModulo
			case ':':
				tokType = TokenColon
			case ',':
				tokType = TokenComma
			case '.':
				tokType = TokenDot
			case '?':
				tokType = TokenQuestion
			case '!':
				tokType = TokenExclaim
			case '@':
				tokType = TokenAt
			case '#':
				tokType = TokenHash
			case '$':
				tokType = TokenDollar
			case '^':
				tokType = TokenCaret
			case '&':
				tokType = TokenAmpersand
			case '|':
				tokType = TokenPipe
			case '~':
				tokType = TokenTilde
			case '(':
				tokType = TokenLParen
			case ')':
				tokType = TokenRParen
			case '[':
				tokType = TokenLBracket
			case ']':
				tokType = TokenRBracket
			case '{':
				tokType = TokenLBrace
			case '}':
				tokType = TokenRBrace
			case ';':
				tokType = TokenSemicolon
			default:
				tokType = TokenEOF
			}
			tokens = append(tokens, Token{Type: tokType, Value: string(ch), Position: pos})
			offset++
			column++
		}
	}

	// Add EOF token
	tokens = append(tokens, Token{Type: TokenEOF, Value: "", Position: Position{Line: line, Column: column, Offset: offset}})
	return tokens
}
