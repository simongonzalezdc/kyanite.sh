package pattern

import (
	"testing"
	"time"
)

func TestLexer_NextToken(t *testing.T) {
	input := `<kick:1> <snare:2> <hihat:3> ~`

	tests := []struct {
		expectedType TokenType
		expectedValue string
	}{
		{TokenSample, "kick:1"},
		{TokenSample, "snare:2"},
		{TokenSample, "hihat:3"},
		{TokenTilde, "~"},
		{TokenEOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Errorf("tests[%d] - value wrong. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}

func TestLexer_Notes(t *testing.T) {
	input := `<c4> <d#5> <eb3> ~`

	tests := []struct {
		expectedType TokenType
		expectedValue string
	}{
		{TokenNote, "c4"},
		{TokenNote, "d#5"},
		{TokenNote, "eb3"},
		{TokenTilde, "~"},
		{TokenEOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Errorf("tests[%d] - value wrong. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}

func TestLexer_Operators(t *testing.T) {
	input := "+ - * / % : , . ? ! @ # $ ^ & | ~ ( ) [ ] { } = != < <= > >= ;"

	tests := []struct {
		expectedType TokenType
		expectedValue string
	}{
		{TokenPlus, "+"},
		{TokenMinus, "-"},
		{TokenMultiply, "*"},
		{TokenDivide, "/"},
		{TokenModulo, "%"},
		{TokenColon, ":"},
		{TokenComma, ","},
		{TokenDot, "."},
		{TokenQuestion, "?"},
		{TokenExclaim, "!"},
		{TokenAt, "@"},
		{TokenHash, "#"},
		{TokenDollar, "$"},
		{TokenCaret, "^"},
		{TokenAmpersand, "&"},
		{TokenPipe, "|"},
		{TokenTilde, "~"},
		{TokenLParen, "("},
		{TokenRParen, ")"},
		{TokenLBracket, "["},
		{TokenRBracket, "]"},
		{TokenLBrace, "{"},
		{TokenRBrace, "}"},
		{TokenEquals, "="},
		{TokenNotEq, "!="},
		{TokenLess, "<"},
		{TokenLessEq, "<="},
		{TokenGreater, ">"},
		{TokenGreaterEq, ">="},
		{TokenSemicolon, ";"},
		{TokenEOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Errorf("tests[%d] - value wrong. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}

func TestParser_SimplePattern(t *testing.T) {
	input := `[<kick:1> <snare:2> <hihat:3> ~]`

	program, errors := ParsePattern(input)
	if len(errors) != 0 {
		t.Errorf("parser returned %d errors", len(errors))
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ListExpression)
	if !ok {
		t.Fatalf("exp not *ListExpression. got=%T", stmt.Expression)
	}

	if len(exp.Elements) != 4 {
		t.Fatalf("list has not 4 elements. got=%d", len(exp.Elements))
	}

	// Check first element (kick)
	kickLit, ok := exp.Elements[0].(*LiteralExpression)
	if !ok {
		t.Fatalf("element not *LiteralExpression. got=%T", exp.Elements[0])
	}

	kickSample, ok := kickLit.Value.(SampleValue)
	if !ok {
		t.Fatalf("value not SampleValue. got=%T", kickLit.Value)
	}

	if kickSample.Name != "kick:1" {
		t.Errorf("sample name wrong. expected='kick:1', got='%s'", kickSample.Name)
	}
}

func TestParser_NotePattern(t *testing.T) {
	input := `[<c4> <d#5> <eb3> ~]`

	program, errors := ParsePattern(input)
	if len(errors) != 0 {
		t.Errorf("parser returned %d errors", len(errors))
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ListExpression)
	if !ok {
		t.Fatalf("exp not *ListExpression. got=%T", stmt.Expression)
	}

	if len(exp.Elements) != 4 {
		t.Fatalf("list has not 4 elements. got=%d", len(exp.Elements))
	}

	// Check first element (c4)
	noteLit, ok := exp.Elements[0].(*LiteralExpression)
	if !ok {
		t.Fatalf("element not *LiteralExpression. got=%T", exp.Elements[0])
	}

	note, ok := noteLit.Value.(NoteValue)
	if !ok {
		t.Fatalf("value not NoteValue. got=%T", noteLit.Value)
	}

	if note.Note != "c" {
		t.Errorf("note name wrong. expected='c', got='%s'", note.Note)
	}

	if note.Octave != 4 {
		t.Errorf("octave wrong. expected=4, got=%d", note.Octave)
	}

	// Check second element (d#5)
	noteLit2, ok := exp.Elements[1].(*LiteralExpression)
	if !ok {
		t.Fatalf("element not *LiteralExpression. got=%T", exp.Elements[1])
	}

	note2, ok := noteLit2.Value.(NoteValue)
	if !ok {
		t.Fatalf("value not NoteValue. got=%T", noteLit2.Value)
	}

	if note2.Note != "d" {
		t.Errorf("note name wrong. expected='d', got='%s'", note2.Note)
	}

	if note2.Accident != "#" {
		t.Errorf("accidental wrong. expected='#', got='%s'", note2.Accident)
	}

	if note2.Octave != 5 {
		t.Errorf("octave wrong. expected=5, got=%d", note2.Octave)
	}
}

func TestParser_FunctionCall(t *testing.T) {
	input := `[<kick:1> <snare:2>].fast(2)`

	program, errors := ParsePattern(input)
	if len(errors) != 0 {
		t.Errorf("parser returned %d errors", len(errors))
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
			program.Statements[0])
	}

	modExpr, ok := stmt.Expression.(*ModifierExpression)
	if !ok {
		t.Fatalf("exp not *ModifierExpression. got=%T", stmt.Expression)
	}

	if modExpr.Modifier != "fast" {
		t.Errorf("modifier name wrong. expected='fast', got='%s'", modExpr.Modifier)
	}

	if len(modExpr.Arguments) != 1 {
		t.Fatalf("modifier has not 1 arguments. got=%d", len(modExpr.Arguments))
	}

	arg, ok := modExpr.Arguments[0].(*LiteralExpression)
	if !ok {
		t.Fatalf("argument not *LiteralExpression. got=%T", modExpr.Arguments[0])
	}

	num, ok := arg.Value.(NumberValue)
	if !ok {
		t.Fatalf("argument value not NumberValue. got=%T", arg.Value)
	}

	if num.Value != 2 {
		t.Errorf("argument value wrong. expected=2, got=%f", num.Value)
	}
}

func TestParser_BinaryExpression(t *testing.T) {
	input := `1 + 2 * 3`

	program, errors := ParsePattern(input)
	if len(errors) != 0 {
		t.Errorf("parser returned %d errors", len(errors))
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
			program.Statements[0])
	}

	binExpr, ok := stmt.Expression.(*BinaryExpression)
	if !ok {
		t.Fatalf("exp not *BinaryExpression. got=%T", stmt.Expression)
	}

	if binExpr.Operator != TokenPlus {
		t.Errorf("operator wrong. expected='+', got='%s'", binExpr.Operator)
	}

	// Check left side
	leftLit, ok := binExpr.Left.(*LiteralExpression)
	if !ok {
		t.Fatalf("left not *LiteralExpression. got=%T", binExpr.Left)
	}

	leftNum, ok := leftLit.Value.(NumberValue)
	if !ok {
		t.Fatalf("left value not NumberValue. got=%T", leftLit.Value)
	}

	if leftNum.Value != 1 {
		t.Errorf("left value wrong. expected=1, got=%f", leftNum.Value)
	}

	// Check right side (should be another binary expression)
	rightBin, ok := binExpr.Right.(*BinaryExpression)
	if !ok {
		t.Fatalf("right not *BinaryExpression. got=%T", binExpr.Right)
	}

	if rightBin.Operator != TokenMultiply {
		t.Errorf("right operator wrong. expected='*', got='%s'", rightBin.Operator)
	}
}

func TestParser_GroupExpression(t *testing.T) {
	input := `(1 + 2) * 3`

	program, errors := ParsePattern(input)
	if len(errors) != 0 {
		t.Errorf("parser returned %d errors", len(errors))
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
			program.Statements[0])
	}

	binExpr, ok := stmt.Expression.(*BinaryExpression)
	if !ok {
		t.Fatalf("exp not *BinaryExpression. got=%T", stmt.Expression)
	}

	if binExpr.Operator != TokenMultiply {
		t.Errorf("operator wrong. expected='*', got='%s'", binExpr.Operator)
	}

	// Check left side (should be group expression)
	leftGroup, ok := binExpr.Left.(*GroupExpression)
	if !ok {
		t.Fatalf("left not *GroupExpression. got=%T", binExpr.Left)
	}

	// Check inside group
	innerBin, ok := leftGroup.Expression.(*BinaryExpression)
	if !ok {
		t.Fatalf("group expression not *BinaryExpression. got=%T", leftGroup.Expression)
	}

	if innerBin.Operator != TokenPlus {
		t.Errorf("inner operator wrong. expected='+', got='%s'", innerBin.Operator)
	}
}

func TestParser_ErrorHandling(t *testing.T) {
	input := `[<kick:1> <snare:2>` // Missing closing bracket

	program, errors := ParsePattern(input)
	
	// Should have at least one error
	if len(errors) == 0 {
		t.Errorf("expected parser errors, but got none")
	}

	// Program should still be returned but may be incomplete
	if program == nil {
		t.Error("program should not be nil even with errors")
	}
}

func TestValidator_EmptyPattern(t *testing.T) {
	pattern := &Pattern{
		Values:   make([]PatternValue, 0),
		Metadata: make(map[string]interface{}),
	}

	validator := NewValidator()
	errors := validator.Validate(pattern)

	if len(errors) == 0 {
		t.Error("expected validation error for empty pattern")
	}

	expectedError := "pattern must contain at least one value"
	found := false
	for _, err := range errors {
		if err.Message == expectedError {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected error '%s' not found", expectedError)
	}
}

func TestValidator_InvalidNote(t *testing.T) {
	pattern := &Pattern{
		Values: []PatternValue{
			NoteValue{Note: "h", Octave: 4}, // Invalid note
		},
		Metadata: make(map[string]interface{}),
	}

	validator := NewValidator()
	errors := validator.Validate(pattern)

	if len(errors) == 0 {
		t.Error("expected validation error for invalid note")
	}

	expectedError := "invalid note name: h"
	found := false
	for _, err := range errors {
		if err.Message == expectedError {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected error '%s' not found", expectedError)
	}
}

func TestValidator_OutOfRangeOctave(t *testing.T) {
	pattern := &Pattern{
		Values: []PatternValue{
			NoteValue{Note: "c", Octave: 15}, // Out of range octave
		},
		Metadata: make(map[string]interface{}),
	}

	validator := NewValidator()
	errors := validator.Validate(pattern)

	if len(errors) == 0 {
		t.Error("expected validation error for out of range octave")
	}

	expectedError := "octave 15 out of range (0-9)"
	found := false
	for _, err := range errors {
		if err.Message == expectedError {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected error '%s' not found", expectedError)
	}
}

func TestPerformanceMetrics(t *testing.T) {
	input := `<kick:1> <snare:2> <hihat:3> <kick:1> <snare:2> <hihat:3>`

	program, errors, metrics := ParseWithMetrics(input)

	if len(errors) != 0 {
		t.Errorf("parser returned %d errors", len(errors))
	}

	if program == nil {
		t.Error("program should not be nil")
	}

	if metrics.TotalTime <= 0 {
		t.Error("total time should be positive")
	}

	if metrics.TokenCount <= 0 {
		t.Error("token count should be positive")
	}

	if metrics.NodeCount <= 0 {
		t.Error("node count should be positive")
	}

	// Performance requirement: parsing should be under 100ms
	if metrics.ParseTime > 100*time.Millisecond {
		t.Errorf("parsing took too long: %v (should be under 100ms)", metrics.ParseTime)
	}
}

func TestFastTokenize(t *testing.T) {
	input := `kick:1 snare:2 hihat:3 1 2 3`

	tokens := FastTokenize(input)

	if len(tokens) == 0 {
		t.Error("no tokens returned")
	}

	// Check that we have at least some tokens
	if len(tokens) < 6 {
		t.Errorf("expected at least 6 tokens, got %d", len(tokens))
	}

	// Check first few tokens
	expectedTypes := []TokenType{
		TokenSample, TokenSample, TokenSample,
		TokenNumber, TokenNumber, TokenNumber,
	}

	for i, expectedType := range expectedTypes {
		if i >= len(tokens) {
			t.Errorf("missing token at index %d", i)
			continue
		}

		if tokens[i].Type != expectedType {
			t.Errorf("token %d type wrong. expected=%q, got=%q",
				i, expectedType, tokens[i].Type)
		}
	}
}

func TestOptimizer_ConsecutiveRests(t *testing.T) {
	pattern := &Pattern{
		Values: []PatternValue{
			SampleValue{Name: "kick:1"},
			RestValue{},
			RestValue{}, // This should be removed
			RestValue{}, // This should be removed
			SampleValue{Name: "snare:2"},
		},
		Metadata: make(map[string]interface{}),
	}

	optimizer := NewPatternOptimizer()
	optimized := optimizer.Optimize(pattern)

	if len(optimized.Values) != 3 {
		t.Errorf("expected 3 values after optimization, got %d", len(optimized.Values))
	}

	// Check that consecutive rests were collapsed to one
	restCount := 0
	for _, value := range optimized.Values {
		if _, isRest := value.(RestValue); isRest {
			restCount++
		}
	}

	if restCount != 1 {
		t.Errorf("expected 1 rest after optimization, got %d", restCount)
	}
}

func BenchmarkParser_SimplePattern(b *testing.B) {
	input := `<kick:1> <snare:2> <hihat:3> ~`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParsePattern(input)
	}
}

func BenchmarkLexer_Tokenize(b *testing.B) {
	input := `<kick:1> <snare:2> <hihat:3> ~ 1 2 3 4 5 6 7 8`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		_ = lexer.Tokenize()
	}
}

func BenchmarkFastTokenize(b *testing.B) {
	input := `<kick:1> <snare:2> <hihat:3> ~ 1 2 3 4 5 6 7 8`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FastTokenize(input)
	}
}