package pattern

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

// SafeParserTest provides crash-safe testing for pattern parsing
type SafeParserTest struct {
	timeout     time.Duration
	maxMemoryMB uint64
}

// NewSafeParserTest creates a new safe parser test instance
func NewSafeParserTest() *SafeParserTest {
	return &SafeParserTest{
		timeout:     5 * time.Second,
		maxMemoryMB: 100, // 100MB limit
	}
}

// TestLexer_NextToken_Safe tests lexer with safety bounds
func TestLexer_NextToken_Safe(t *testing.T) {
	input := `<kick:1> <snare:2> <hihat:3> ~`

	tests := []struct {
		expectedType  TokenType
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
		// Add timeout protection
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		done := make(chan Token, 1)
		go func() {
			done <- l.NextToken()
		}()

		select {
		case tok := <-done:
			if tok.Type != tt.expectedType {
				t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q",
					i, tt.expectedType, tok.Type)
			}

			if tok.Value != tt.expectedValue {
				t.Errorf("tests[%d] - value wrong. expected=%q, got=%q",
					i, tt.expectedValue, tok.Value)
			}
		case <-ctx.Done():
			t.Errorf("tests[%d] - lexer timed out", i)
		}
	}
}

// TestLexer_Notes_Safe tests note lexing with safety bounds
func TestLexer_Notes_Safe(t *testing.T) {
	input := `<c4> <d#5> <eb3> ~`

	tests := []struct {
		expectedType  TokenType
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
		// Add timeout protection
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		done := make(chan Token, 1)
		go func() {
			done <- l.NextToken()
		}()

		select {
		case tok := <-done:
			if tok.Type != tt.expectedType {
				t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q",
					i, tt.expectedType, tok.Type)
			}

			if tok.Value != tt.expectedValue {
				t.Errorf("tests[%d] - value wrong. expected=%q, got=%q",
					i, tt.expectedValue, tok.Value)
			}
		case <-ctx.Done():
			t.Errorf("tests[%d] - lexer timed out", i)
		}
	}
}

// TestParser_SimplePattern_Safe tests simple pattern parsing with safety
func TestParser_SimplePattern_Safe(t *testing.T) {
	input := `[<kick:1> <snare:2> <hihat:3> ~]`

	// Add memory and timeout protection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Monitor memory usage
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	initialMemory := memStats.Alloc

	done := make(chan struct {
		program *Program
		errors  []PatternError
	}, 1)

	go func() {
		program, errors := ParsePattern(input)
		done <- struct {
			program *Program
			errors  []PatternError
		}{program, errors}
	}()

	select {
	case result := <-done:
		// Check for excessive memory usage
		runtime.ReadMemStats(&memStats)
		memoryUsed := memStats.Alloc - initialMemory
		if memoryUsed > 50*1024*1024 { // 50MB limit for this test
			t.Errorf("test used too much memory: %d bytes", memoryUsed)
		}

		if len(result.errors) != 0 {
			t.Errorf("parser returned %d errors: %v", len(result.errors), result.errors)
		}

		if result.program == nil {
			t.Fatal("program should not be nil")
		}

		if len(result.program.Statements) != 1 {
			t.Fatalf("program has not 1 statements. got=%d", len(result.program.Statements))
		}

		stmt, ok := result.program.Statements[0].(*ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
				result.program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ListExpression)
		if !ok {
			t.Fatalf("exp not *ListExpression. got=%T", stmt.Expression)
		}

		if len(exp.Elements) != 4 {
			t.Fatalf("list has not 4 elements. got=%d", len(exp.Elements))
		}

	case <-ctx.Done():
		t.Fatal("parser timed out")
	}
}

// TestParser_ErrorHandling_Safe tests error handling with safety bounds
func TestParser_ErrorHandling_Safe(t *testing.T) {
	input := `[<kick:1> <snare:2>` // Missing closing bracket

	// Add timeout protection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan struct {
		program *Program
		errors  []PatternError
	}, 1)

	go func() {
		program, errors := ParsePattern(input)
		done <- struct {
			program *Program
			errors  []PatternError
		}{program, errors}
	}()

	select {
	case result := <-done:
		// Should have at least one error for malformed input
		if len(result.errors) == 0 {
			t.Error("expected parser errors for malformed input, but got none")
		}

		// Program should still be returned even with errors
		if result.program == nil {
			t.Error("program should not be nil even with errors")
		}
	case <-ctx.Done():
		t.Fatal("parser timed out")
	}
}

// TestValidator_EmptyPattern_Safe tests validation with safety
func TestValidator_EmptyPattern_Safe(t *testing.T) {
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

// TestValidator_InvalidNote_Safe tests invalid note validation with safety
func TestValidator_InvalidNote_Safe(t *testing.T) {
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

// TestPerformanceMetrics_Safe tests performance with safety bounds
func TestPerformanceMetrics_Safe(t *testing.T) {
	input := `<kick:1> <snare:2> <hihat:3> <kick:1> <snare:2> <hihat:3>`

	// Add timeout protection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	done := make(chan struct {
		program *Program
		errors  []PatternError
		metrics PerformanceMetrics
	}, 1)

	go func() {
		program, errors, metrics := ParseWithMetrics(input)
		done <- struct {
			program *Program
			errors  []PatternError
			metrics PerformanceMetrics
		}{program, errors, metrics}
	}()

	select {
	case result := <-done:
		if len(result.errors) != 0 {
			t.Errorf("parser returned %d errors", len(result.errors))
		}

		if result.program == nil {
			t.Error("program should not be nil")
		}

		if result.metrics.TotalTime <= 0 {
			t.Error("total time should be positive")
		}

		if result.metrics.TokenCount <= 0 {
			t.Error("token count should be positive")
		}

		if result.metrics.NodeCount <= 0 {
			t.Error("node count should be positive")
		}

		// Performance requirement: parsing should be under 100ms
		if result.metrics.ParseTime > 100*time.Millisecond {
			t.Errorf("parsing took too long: %v (should be under 100ms)", result.metrics.ParseTime)
		}
	case <-ctx.Done():
		t.Fatal("parser timed out")
	}
}

// TestFastTokenize_Safe tests fast tokenization with safety bounds
func TestFastTokenize_Safe(t *testing.T) {
	input := `kick:1 snare:2 hihat:3 1 2 3`

	// Add timeout protection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan []Token, 1)
	go func() {
		tokens := FastTokenize(input)
		done <- tokens
	}()

	select {
	case tokens := <-done:
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
	case <-ctx.Done():
		t.Fatal("fast tokenization timed out")
	}
}

// TestOptimizer_ConsecutiveRests_Safe tests optimization with safety
func TestOptimizer_ConsecutiveRests_Safe(t *testing.T) {
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

	// Add timeout protection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan *Pattern, 1)
	go func() {
		optimizer := NewPatternOptimizer()
		optimized := optimizer.Optimize(pattern)
		done <- optimized
	}()

	select {
	case optimized := <-done:
		if optimized == nil {
			t.Fatal("optimized pattern should not be nil")
		}

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
	case <-ctx.Done():
		t.Fatal("optimizer timed out")
	}
}

// TestStressTest_Safe performs stress testing with safety bounds
func TestStressTest_Safe(t *testing.T) {
	// Test with a large but reasonable input
	input := `<kick:1> <snare:2> <hihat:3> <kick:1> <snare:2> <hihat:3>`

	// Repeat the pattern to create a larger input
	largeInput := ""
	for i := 0; i < 100; i++ {
		largeInput += input + " "
	}

	// Add timeout and memory protection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Monitor memory usage
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	initialMemory := memStats.Alloc

	done := make(chan struct {
		program *Program
		errors  []PatternError
	}, 1)

	go func() {
		program, errors := ParsePattern(largeInput)
		done <- struct {
			program *Program
			errors  []PatternError
		}{program, errors}
	}()

	select {
	case result := <-done:
		// Check memory usage
		runtime.ReadMemStats(&memStats)
		memoryUsed := memStats.Alloc - initialMemory
		if memoryUsed > 100*1024*1024 { // 100MB limit
			t.Errorf("stress test used too much memory: %d bytes", memoryUsed)
		}

		if len(result.errors) != 0 {
			t.Errorf("parser returned %d errors in stress test", len(result.errors))
		}

		if result.program == nil {
			t.Error("program should not be nil in stress test")
		}

	case <-ctx.Done():
		t.Fatal("stress test timed out")
	}
}

// TestMalformedInput_Safe tests various malformed inputs with safety
func TestMalformedInput_Safe(t *testing.T) {
	malformedInputs := []string{
		`<kick:1> <snare:2`,      // Missing closing bracket
		`<unclosed`,              // Unclosed sample/note
		`[<kick:1> <snare:2>`,    // Missing closing list bracket
		`(<kick:1> + <snare:2>)`, // Missing closing paren
		`{<kick:1> <snare:2>}`,   // Missing closing brace
		`<invalid note>`,         // Invalid note format
		`<kick>`,                 // Missing sample number
		`[<kick:1> <snare:2>] <kick:3> <snare:4>`, // Mixed formats
	}

	for i, input := range malformedInputs {
		t.Run(fmt.Sprintf("malformed_%d", i), func(t *testing.T) {
			// Add timeout protection for each test case
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			done := make(chan struct {
				program *Program
				errors  []PatternError
			}, 1)

			go func() {
				program, errors := ParsePattern(input)
				done <- struct {
					program *Program
					errors  []PatternError
				}{program, errors}
			}()

			select {
			case result := <-done:
				// For malformed input, we expect either errors or nil program
				// The important thing is that it doesn't crash or hang
				if result.program == nil && len(result.errors) == 0 {
					t.Errorf("malformed input %d: got neither program nor errors", i)
				}
			case <-ctx.Done():
				t.Errorf("malformed input %d timed out", i)
			}
		})
	}
}

// BenchmarkParser_SimplePattern_Safe provides safe benchmarking
func BenchmarkParser_SimplePattern_Safe(b *testing.B) {
	input := `<kick:1> <snare:2> <hihat:3> ~`

	// Limit benchmark time to prevent infinite loops
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		select {
		case <-ctx.Done():
			b.Fatalf("benchmark timed out")
		default:
			_, _ = ParsePattern(input)
		}
	}
}

// BenchmarkLexer_Tokenize_Safe provides safe lexer benchmarking
func BenchmarkLexer_Tokenize_Safe(b *testing.B) {
	input := `<kick:1> <snare:2> <hihat:3> ~ 1 2 3 4 5 6 7 8`

	// Limit benchmark time
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		select {
		case <-ctx.Done():
			b.Fatalf("benchmark timed out")
		default:
			lexer := NewLexer(input)
			_ = lexer.Tokenize()
		}
	}
}

// BenchmarkFastTokenize_Safe provides safe fast tokenization benchmarking
func BenchmarkFastTokenize_Safe(b *testing.B) {
	input := `<kick:1> <snare:2> <hihat:3> ~ 1 2 3 4 5 6 7 8`

	// Limit benchmark time
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		select {
		case <-ctx.Done():
			b.Fatalf("benchmark timed out")
		default:
			_ = FastTokenize(input)
		}
	}
}
