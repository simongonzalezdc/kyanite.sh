package pattern

import (
	"time"
)

// TokenType represents the type of a token in the pattern language
type TokenType int

const (
	// Literals
	TokenEOF TokenType = iota
	TokenNumber
	TokenString
	TokenIdentifier

	// Operators
	TokenPlus      // +
	TokenMinus     // -
	TokenMultiply  // *
	TokenDivide    // /
	TokenModulo    // %
	TokenSlash     // / (alias for Divide)
	TokenAsterisk  // * (alias for Multiply)
	TokenColon     // :
	TokenComma     // ,
	TokenDot       // .
	TokenQuestion  // ?
	TokenExclaim   // !
	TokenAt        // @
	TokenHash      // #
	TokenDollar    // $
	TokenPercent   // %
	TokenCaret     // ^
	TokenAmpersand // &
	TokenPipe      // |
	TokenTilde     // ~
	TokenEquals    // =
	TokenNotEq     // !=
	TokenLess      // <
	TokenLessEq    // <=
	TokenGreater   // >
	TokenGreaterEq // >=
	TokenSemicolon // ;

	// Brackets
	TokenLParen   // (
	TokenRParen   // )
	TokenLBracket // [
	TokenRBracket // ]
	TokenLBrace   // {
	TokenRBrace   // }

	// Keywords
	TokenPattern
	TokenSample
	TokenNote
	TokenRest
	TokenSpeed
	TokenVolume
	TokenPan
	TokenDelay
	TokenReverb
	TokenFilter
	TokenLoop
	TokenRand
	TokenIrand
	TokenChoose
	TokenDegenerate
	TokenSometimes
	TokenOften
	TokenRarely
	TokenNever
	TokenSuperimpose
	TokenChop
	TokenGap
	TokenOff
	TokenJux
	TokenStack
	TokenOr
	TokenAnd
	TokenSeq
	TokenFast
	TokenSlow
	TokenDegradeBy
	TokenGain
	TokenRoom
	TokenSize
	TokenDensity
)

// Token represents a lexical token in the pattern language
type Token struct {
	Type     TokenType
	Value    string
	Position Position
}

// Position represents a position in source code
type Position struct {
	Line   int
	Column int
	Offset int
}

// PatternValue represents a value in a pattern
type PatternValue interface {
	IsPatternValue()
}

// NumberValue represents a numeric value
type NumberValue struct {
	Value float64
}

func (n NumberValue) IsPatternValue() {}

// StringValue represents a string value
type StringValue struct {
	Value string
}

func (s StringValue) IsPatternValue() {}

// SampleValue represents a sample reference
type SampleValue struct {
	Name string
}

func (s SampleValue) IsPatternValue() {}

// NoteValue represents a musical note
type NoteValue struct {
	Note     string
	Octave   int
	Accident string
}

func (n NoteValue) IsPatternValue() {}

// RestValue represents a rest/silence
type RestValue struct{}

func (r RestValue) IsPatternValue() {}

// Pattern represents a musical pattern
type Pattern struct {
	Values   []PatternValue
	Duration time.Duration
	Events   []PatternEvent
	Metadata map[string]interface{}
}

// PatternEvent represents a single event in a pattern
type PatternEvent struct {
	Value    PatternValue
	Time     time.Duration
	Duration time.Duration
	Params   map[string]PatternValue
}

// String returns the string representation of a TokenType
func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenNumber:
		return "NUMBER"
	case TokenString:
		return "STRING"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenMultiply:
		return "*"
	case TokenDivide:
		return "/"
	case TokenModulo:
		return "%"
	case TokenSlash:
		return "/"
	case TokenAsterisk:
		return "*"
	case TokenColon:
		return ":"
	case TokenComma:
		return ","
	case TokenDot:
		return "."
	case TokenQuestion:
		return "?"
	case TokenExclaim:
		return "!"
	case TokenAt:
		return "@"
	case TokenHash:
		return "#"
	case TokenDollar:
		return "$"
	case TokenPercent:
		return "%"
	case TokenCaret:
		return "^"
	case TokenAmpersand:
		return "&"
	case TokenPipe:
		return "|"
	case TokenTilde:
		return "~"
	case TokenEquals:
		return "="
	case TokenNotEq:
		return "!="
	case TokenLess:
		return "<"
	case TokenLessEq:
		return "<="
	case TokenGreater:
		return ">"
	case TokenGreaterEq:
		return ">="
	case TokenSemicolon:
		return ";"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenLBracket:
		return "["
	case TokenRBracket:
		return "]"
	case TokenLBrace:
		return "{"
	case TokenRBrace:
		return "}"
	case TokenPattern:
		return "pattern"
	case TokenSample:
		return "sample"
	case TokenNote:
		return "note"
	case TokenRest:
		return "rest"
	case TokenSpeed:
		return "speed"
	case TokenVolume:
		return "volume"
	case TokenPan:
		return "pan"
	case TokenDelay:
		return "delay"
	case TokenReverb:
		return "reverb"
	case TokenFilter:
		return "filter"
	case TokenLoop:
		return "loop"
	case TokenRand:
		return "rand"
	case TokenIrand:
		return "irand"
	case TokenChoose:
		return "choose"
	case TokenDegenerate:
		return "degenerate"
	case TokenSometimes:
		return "sometimes"
	case TokenOften:
		return "often"
	case TokenRarely:
		return "rarely"
	case TokenNever:
		return "never"
	case TokenSuperimpose:
		return "superimpose"
	case TokenChop:
		return "chop"
	case TokenGap:
		return "gap"
	case TokenOff:
		return "off"
	case TokenJux:
		return "jux"
	case TokenStack:
		return "stack"
	case TokenOr:
		return "or"
	case TokenAnd:
		return "and"
	case TokenSeq:
		return "seq"
	case TokenFast:
		return "fast"
	case TokenSlow:
		return "slow"
	case TokenDegradeBy:
		return "degradeBy"
	case TokenGain:
		return "gain"
	case TokenRoom:
		return "room"
	case TokenSize:
		return "size"
	case TokenDensity:
		return "density"
	default:
		return "UNKNOWN"
	}
}

// PatternState represents the current state of pattern evaluation
type PatternState struct {
	CurrentTime time.Duration
	Cycle       int
	Phase       float64
}

// PatternFunction represents a function that can be applied to patterns
type PatternFunction interface {
	Apply(pattern *Pattern, state PatternState) *Pattern
	Name() string
}

// PatternError represents an error in pattern parsing or evaluation
type PatternError struct {
	Message  string
	Position Position
	Cause    error
}

func (e PatternError) Error() string {
	if e.Position.Line > 0 {
		return e.Message + " at line " + string(rune(e.Position.Line)) + ", column " + string(rune(e.Position.Column))
	}
	return e.Message
}

// ParseResult represents the result of parsing a pattern
type ParseResult struct {
	Pattern *Pattern
	Errors  []PatternError
	Success bool
}

// ValidationRule represents a rule for validating patterns
type ValidationRule interface {
	Validate(pattern *Pattern) []PatternError
	Name() string
}

// PerformanceMetrics tracks parsing performance
type PerformanceMetrics struct {
	LexerTime      time.Duration
	ParseTime      time.Duration
	ValidationTime time.Duration
	TotalTime      time.Duration
	TokenCount     int
	NodeCount      int
}
