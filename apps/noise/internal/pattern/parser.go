package pattern

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Parser implements a recursive descent parser for the pattern language
type Parser struct {
	lexer          *Lexer
	currentToken   Token
	peekToken      Token
	errors         []PatternError
	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

// NewParser creates a new parser
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: make([]PatternError, 0),
	}

	// Register prefix parse functions
	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(TokenNumber, p.parseNumberLiteral)
	p.registerPrefix(TokenString, p.parseStringLiteral)
	p.registerPrefix(TokenIdentifier, p.parseIdentifier)
	p.registerPrefix(TokenSample, p.parseSampleLiteral)
	p.registerPrefix(TokenNote, p.parseNoteLiteral)
	p.registerPrefix(TokenRest, p.parseRestLiteral)
	p.registerPrefix(TokenLParen, p.parseGroupExpression)
	p.registerPrefix(TokenLBracket, p.parseListExpression)
	p.registerPrefix(TokenLBrace, p.parseStructExpression)
	p.registerPrefix(TokenMinus, p.parseUnaryExpression)
	p.registerPrefix(TokenTilde, p.parseUnaryExpression)
	p.registerPrefix(TokenQuestion, p.parseUnaryExpression)
	p.registerPrefix(TokenExclaim, p.parseUnaryExpression)

	// Register infix parse functions
	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(TokenPlus, p.parseBinaryExpression)
	p.registerInfix(TokenMinus, p.parseBinaryExpression)
	p.registerInfix(TokenMultiply, p.parseBinaryExpression)
	p.registerInfix(TokenDivide, p.parseBinaryExpression)
	p.registerInfix(TokenModulo, p.parseBinaryExpression)
	p.registerInfix(TokenColon, p.parseParameterExpression)
	p.registerInfix(TokenDot, p.parseModifierExpression)
	p.registerInfix(TokenLBracket, p.parseIndexExpression)

	// Read two tokens to set currentToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// registerPrefix registers a prefix parse function for a token type
func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix registers an infix parse function for a token type
func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// nextToken advances to the next token
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	next := p.lexer.NextToken()
	p.peekToken = next
}

// Parse parses the input and returns a program
func (p *Parser) Parse() *Program {
	program := &Program{
		Statements: make([]Statement, 0),
	}

	for p.currentToken.Type != TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses a statement
func (p *Parser) parseStatement() Statement {
	switch p.currentToken.Type {
	case TokenIdentifier:
		if p.peekToken.Type == TokenColon || p.peekToken.Type == TokenEquals {
			return p.parsePatternStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parsePatternStatement parses a pattern definition statement
func (p *Parser) parsePatternStatement() *PatternStatement {
	stmt := &PatternStatement{
		Name:     p.currentToken.Value,
		Position: p.currentToken.Position,
	}

	// Skip identifier
	p.nextToken()

	// Skip colon or equals
	if p.currentToken.Type == TokenColon || p.currentToken.Type == TokenEquals {
		p.nextToken()
	}

	// Parse expression
	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

// parseExpressionStatement parses an expression statement
func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{
		Expression: p.parseExpression(LOWEST),
	}

	if p.peekToken.Type == TokenSemicolon {
		p.nextToken()
	}

	return stmt
}

// ExpressionStatement represents an expression as a statement
type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) node()          {}
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) Pos() Position {
	if es.Expression != nil {
		return es.Expression.Pos()
	}
	return Position{}
}

// Precedence levels
const (
	_ int = iota
	LOWEST
	LOGICAL_OR
	LOGICAL_AND
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
	MODIFIER
)

// precedences maps token types to their precedence levels
var precedences = map[TokenType]int{
	TokenOr:        LOGICAL_OR,
	TokenAnd:       LOGICAL_AND,
	TokenEquals:    EQUALS,
	TokenNotEq:     EQUALS,
	TokenLess:      LESSGREATER,
	TokenLessEq:    LESSGREATER,
	TokenGreater:   LESSGREATER,
	TokenGreaterEq: LESSGREATER,
	TokenPlus:      SUM,
	TokenMinus:     SUM,
	TokenSlash:     PRODUCT,
	TokenAsterisk:  PRODUCT,
	TokenModulo:    PRODUCT,
	TokenLParen:    CALL,
	TokenLBracket:  INDEX,
	TokenDot:       MODIFIER,
}

// parseExpression parses an expression with the given precedence
func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(TokenSemicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// parseNumberLiteral parses a number literal
func (p *Parser) parseNumberLiteral() Expression {
	value, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		p.errors = append(p.errors, PatternError{
			Message:  fmt.Sprintf("could not parse %q as number", p.currentToken.Value),
			Position: p.currentToken.Position,
		})
		return nil
	}

	return &LiteralExpression{
		Value:    NumberValue{Value: value},
		Position: p.currentToken.Position,
	}
}

// parseStringLiteral parses a string literal
func (p *Parser) parseStringLiteral() Expression {
	return &LiteralExpression{
		Value:    StringValue{Value: p.currentToken.Value},
		Position: p.currentToken.Position,
	}
}

// parseIdentifier parses an identifier or function call
func (p *Parser) parseIdentifier() Expression {
	// Check if this is a function call
	if p.peekToken.Type == TokenLParen {
		return p.parseFunctionCall()
	}

	// Check if this is a keyword that should be a literal
	switch p.currentToken.Value {
	case "true":
		return &LiteralExpression{
			Value:    NumberValue{Value: 1},
			Position: p.currentToken.Position,
		}
	case "false":
		return &LiteralExpression{
			Value:    NumberValue{Value: 0},
			Position: p.currentToken.Position,
		}
	}

	// Regular identifier
	return &LiteralExpression{
		Value:    StringValue{Value: p.currentToken.Value},
		Position: p.currentToken.Position,
	}
}

// parseFunctionCall parses a function call
func (p *Parser) parseFunctionCall() Expression {
	expr := &FunctionCallExpression{
		Name:     p.currentToken.Value,
		Position: p.currentToken.Position,
	}

	p.nextToken() // skip identifier

	if !p.expectPeek(TokenLParen) {
		return nil
	}

	if p.peekTokenIs(TokenRParen) {
		p.nextToken() // consume ')'
		return expr
	}

	p.nextToken()
	expr.Arguments = p.parseExpressionList(TokenRParen)

	return expr
}

// parseSampleLiteral parses a sample literal
func (p *Parser) parseSampleLiteral() Expression {
	return &LiteralExpression{
		Value:    SampleValue{Name: p.currentToken.Value},
		Position: p.currentToken.Position,
	}
}

// parseNoteLiteral parses a note literal
func (p *Parser) parseNoteLiteral() Expression {
	// Parse note value (e.g., "c4", "d#5", "eb3")
	noteStr := p.currentToken.Value
	if len(noteStr) == 0 {
		p.errors = append(p.errors, PatternError{
			Message:  "empty note value",
			Position: p.currentToken.Position,
		})
		return nil
	}

	// Extract note name
	note := strings.ToLower(string(noteStr[0]))
	if note < "a" || note > "g" {
		p.errors = append(p.errors, PatternError{
			Message:  fmt.Sprintf("invalid note: %s", note),
			Position: p.currentToken.Position,
		})
		return nil
	}

	// Extract accidentals and octave
	accident := ""
	octave := 4 // default octave

	i := 1
	for i < len(noteStr) && (noteStr[i] == '#' || noteStr[i] == 'b') {
		accident += string(noteStr[i])
		i++
	}

	if i < len(noteStr) {
		o, err := strconv.Atoi(noteStr[i:])
		if err == nil {
			octave = o
		}
	}

	return &LiteralExpression{
		Value:    NoteValue{Note: note, Accident: accident, Octave: octave},
		Position: p.currentToken.Position,
	}
}

// parseRestLiteral parses a rest literal
func (p *Parser) parseRestLiteral() Expression {
	return &LiteralExpression{
		Value:    RestValue{},
		Position: p.currentToken.Position,
	}
}

// parseGroupExpression parses a grouped expression
func (p *Parser) parseGroupExpression() Expression {
	p.nextToken() // skip '('

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(TokenRParen) {
		return nil
	}

	return &GroupExpression{
		Expression: exp,
		Position:   p.currentToken.Position,
	}
}

// parseListExpression parses a list expression
func (p *Parser) parseListExpression() Expression {
	list := &ListExpression{
		Position: p.currentToken.Position,
	}

	p.nextToken() // skip '['

	if p.currentToken.Type != TokenRBracket {
		list.Elements = p.parseExpressionList(TokenRBracket)
	}

	return list
}

// parseStructExpression parses a struct expression
func (p *Parser) parseStructExpression() Expression {
	structExpr := &StructExpression{
		Position: p.currentToken.Position,
	}

	p.nextToken() // skip '{'

	if p.currentToken.Type != TokenRBrace {
		structExpr.Fields = p.parseParameterList(TokenRBrace)
	}

	return structExpr
}

// parseExpressionList parses a list of expressions until the given end token
func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := make([]Expression, 0)

	if p.currentTokenIs(end) {
		return list
	}

	list = append(list, p.parseExpression(LOWEST))

	for !p.peekTokenIs(end) {

		if p.peekTokenIs(TokenComma) {
			p.nextToken() // consume comma
			p.nextToken() // move to next element
		} else {
			p.nextToken()
		}

		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return list
	}

	return list
}

// parseParameterList parses a list of parameters until the given end token
func (p *Parser) parseParameterList(end TokenType) []*ParameterExpression {
	list := make([]*ParameterExpression, 0)

	if p.currentTokenIs(end) {
		return list
	}

	param := p.parseParameter()
	if param != nil {
		list = append(list, param)
	}

	for !p.currentTokenIs(end) {
		if !p.expectPeek(TokenComma) {
			return list
		}
		p.nextToken()
		param := p.parseParameter()
		if param != nil {
			list = append(list, param)
		}
	}

	if !p.expectPeek(end) {
		return list
	}

	return list
}

// parseParameter parses a parameter expression
func (p *Parser) parseParameter() *ParameterExpression {
	param := &ParameterExpression{
		Name:     p.currentToken.Value,
		Position: p.currentToken.Position,
	}

	p.nextToken() // skip name

	if !p.expectPeek(TokenColon) {
		return nil
	}

	p.nextToken() // skip ':'

	param.Value = p.parseExpression(LOWEST)

	return param
}

// parseUnaryExpression parses a unary expression
func (p *Parser) parseUnaryExpression() Expression {
	expression := &UnaryExpression{
		Operator: p.currentToken.Type,
		Position: p.currentToken.Position,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseBinaryExpression parses a binary expression
func (p *Parser) parseBinaryExpression(left Expression) Expression {
	expression := &BinaryExpression{
		Left:     left,
		Operator: p.currentToken.Type,
		Position: p.currentToken.Position,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseParameterExpression parses a parameter expression (infix)
func (p *Parser) parseParameterExpression(left Expression) Expression {
	if ident, ok := left.(*LiteralExpression); ok {
		if strVal, ok := ident.Value.(StringValue); ok {
			return &ParameterExpression{
				Name:     strVal.Value,
				Value:    p.parseExpression(LOWEST),
				Position: left.Pos(),
			}
		}
	}

	p.errors = append(p.errors, PatternError{
		Message:  "invalid parameter expression",
		Position: p.currentToken.Position,
	})
	return nil
}

// parseModifierExpression parses a modifier expression
func (p *Parser) parseModifierExpression(left Expression) Expression {
	expr := &ModifierExpression{
		Base:     left,
		Position: p.currentToken.Position,
	}

	p.nextToken() // skip '.'

	if !isIdentifierLikeToken(p.currentToken.Type) {
		p.errors = append(p.errors, PatternError{
			Message:  "expected modifier name after '.'",
			Position: p.currentToken.Position,
		})
		return nil
	}

	expr.Modifier = p.currentToken.Value

	// Check for arguments
	if p.peekToken.Type == TokenLParen {
		if !p.expectPeek(TokenLParen) {
			return nil
		}

		if p.peekTokenIs(TokenRParen) {
			p.nextToken() // consume ')'
			return expr
		}

		p.nextToken()
		expr.Arguments = p.parseExpressionList(TokenRParen)
	}

	return expr
}

// parseIndexExpression parses an index expression
func (p *Parser) parseIndexExpression(left Expression) Expression {
	expr := &IndexExpression{
		Base:     left,
		Position: p.currentToken.Position,
	}

	p.nextToken() // skip '['

	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(TokenRBracket) {
		return nil
	}

	return expr
}

// Helper functions

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, PatternError{
		Message:  msg,
		Position: p.currentToken.Position,
	})
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, PatternError{
		Message:  msg,
		Position: p.peekToken.Position,
	})
}

// Errors returns the parser errors
func (p *Parser) Errors() []PatternError {
	return p.errors
}

// ParsePattern is a convenience function to parse a pattern string
func ParsePattern(input string) (*Program, []PatternError) {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.Parse()
	return program, parser.Errors()
}

// ParseWithMetrics parses a pattern and returns performance metrics
func ParseWithMetrics(input string) (*Program, []PatternError, PerformanceMetrics) {
	start := time.Now()

	lexerStart := time.Now()
	lex := NewLexer(input)
	tokens := lex.Tokenize()
	lexerTime := time.Since(lexerStart)

	parserStart := time.Now()
	// Reset lexer to beginning state for parsing
	lex.Reset()
	parser := NewParser(lex)
	program := parser.Parse()
	parserTime := time.Since(parserStart)
	if parserTime <= 0 {
		parserTime = time.Nanosecond
	}

	totalTime := time.Since(start)
	if totalTime <= 0 {
		totalTime = time.Nanosecond
	}

	metrics := PerformanceMetrics{
		LexerTime:  lexerTime,
		ParseTime:  parserTime,
		TotalTime:  totalTime,
		TokenCount: len(tokens),
		NodeCount:  countNodes(program),
	}

	return program, parser.Errors(), metrics
}

// countNodes counts the number of nodes in an AST
func isIdentifierLikeToken(t TokenType) bool {
	if t == TokenIdentifier {
		return true
	}
	if t >= TokenPattern && t <= TokenDensity {
		return true
	}
	return false
}

func countNodes(node Node) int {
	if node == nil {
		return 0
	}

	count := 1

	switch n := node.(type) {
	case *Program:
		for _, stmt := range n.Statements {
			count += countNodes(stmt)
		}
	case *PatternStatement:
		count += countNodes(n.Expression)
	case *SequenceExpression:
		for _, value := range n.Values {
			count += countNodes(value)
		}
	case *FunctionCallExpression:
		for _, arg := range n.Arguments {
			count += countNodes(arg)
		}
	case *BinaryExpression:
		count += countNodes(n.Left)
		count += countNodes(n.Right)
	case *UnaryExpression:
		count += countNodes(n.Right)
	case *GroupExpression:
		count += countNodes(n.Expression)
	case *ParameterExpression:
		count += countNodes(n.Value)
	case *ModifierExpression:
		count += countNodes(n.Base)
		for _, arg := range n.Arguments {
			count += countNodes(arg)
		}
	case *TimeExpression:
		count += countNodes(n.Value)
	case *RangeExpression:
		count += countNodes(n.Start)
		count += countNodes(n.End)
		if n.Step != nil {
			count += countNodes(n.Step)
		}
	case *ListExpression:
		for _, element := range n.Elements {
			count += countNodes(element)
		}
	case *StructExpression:
		for _, field := range n.Fields {
			count += countNodes(field)
		}
	case *IndexExpression:
		count += countNodes(n.Base)
		count += countNodes(n.Index)
	case *PropertyExpression:
		count += countNodes(n.Base)
	}

	return count
}
