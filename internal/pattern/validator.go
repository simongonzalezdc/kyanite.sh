package pattern

import (
	"fmt"
	"time"
)

// Validator validates patterns according to various rules
type Validator struct {
	rules []ValidationRule
}

// NewValidator creates a new validator with default rules
func NewValidator() *Validator {
	v := &Validator{
		rules: make([]ValidationRule, 0),
	}

	// Add default validation rules
	v.AddRule(&SyntaxValidationRule{})
	v.AddRule(&SemanticValidationRule{})
	v.AddRule(&PerformanceValidationRule{})
	v.AddRule(&MusicalValidationRule{})

	return v
}

// AddRule adds a validation rule to the validator
func (v *Validator) AddRule(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

// Validate validates a pattern and returns any errors
func (v *Validator) Validate(pattern *Pattern) []PatternError {
	var errors []PatternError

	for _, rule := range v.rules {
		ruleErrors := rule.Validate(pattern)
		errors = append(errors, ruleErrors...)
	}

	return errors
}

// ValidateProgram validates an entire program
func (v *Validator) ValidateProgram(program *Program) []PatternError {
	var errors []PatternError

	for _, stmt := range program.Statements {
		if patternStmt, ok := stmt.(*PatternStatement); ok {
			pattern := v.evaluatePattern(patternStmt.Expression)
			patternErrors := v.Validate(pattern)
			errors = append(errors, patternErrors...)
		}
	}

	return errors
}

// evaluatePattern evaluates an expression into a pattern
func (v *Validator) evaluatePattern(expr Expression) *Pattern {
	pattern := &Pattern{
		Values:   make([]PatternValue, 0),
		Metadata: make(map[string]interface{}),
	}

	switch e := expr.(type) {
	case *LiteralExpression:
		pattern.Values = append(pattern.Values, e.Value)
	case *ListExpression:
		for _, element := range e.Elements {
			if litExpr, ok := element.(*LiteralExpression); ok {
				pattern.Values = append(pattern.Values, litExpr.Value)
			}
		}
	case *SequenceExpression:
		for _, value := range e.Values {
			if litExpr, ok := value.(*LiteralExpression); ok {
				pattern.Values = append(pattern.Values, litExpr.Value)
			}
		}
	}

	return pattern
}

// SyntaxValidationRule validates the syntax of patterns
type SyntaxValidationRule struct{}

func (r *SyntaxValidationRule) Name() string {
	return "SyntaxValidation"
}

func (r *SyntaxValidationRule) Validate(pattern *Pattern) []PatternError {
	var errors []PatternError

	// Check if pattern has at least one value
	if len(pattern.Values) == 0 {
		errors = append(errors, PatternError{
			Message: "pattern must contain at least one value",
			Position: Position{},
		})
	}

	// Check for invalid pattern combinations
	for _, value := range pattern.Values {
		switch v := value.(type) {
		case NoteValue:
			if v.Octave < 0 || v.Octave > 9 {
				errors = append(errors, PatternError{
					Message: fmt.Sprintf("octave %d out of range (0-9)", v.Octave),
					Position: Position{},
				})
			}
		case SampleValue:
			if v.Name == "" {
				errors = append(errors, PatternError{
					Message: "sample name cannot be empty",
					Position: Position{},
				})
			}
		}
	}

	return errors
}

// SemanticValidationRule validates the semantics of patterns
type SemanticValidationRule struct{}

func (r *SemanticValidationRule) Name() string {
	return "SemanticValidation"
}

func (r *SemanticValidationRule) Validate(pattern *Pattern) []PatternError {
	var errors []PatternError

	// Check for logical consistency
	sampleCount := 0
	noteCount := 0
	restCount := 0

	for _, value := range pattern.Values {
		switch value.(type) {
		case SampleValue:
			sampleCount++
		case NoteValue:
			noteCount++
		case RestValue:
			restCount++
		}
	}

	// Warning if mixing samples and notes (not an error, just informational)
	if sampleCount > 0 && noteCount > 0 {
		errors = append(errors, PatternError{
			Message: "mixing samples and notes in the same pattern may have unexpected results",
			Position: Position{},
		})
	}

	// Check if pattern is all rests
	if restCount == len(pattern.Values) && len(pattern.Values) > 0 {
		errors = append(errors, PatternError{
			Message: "pattern consists entirely of rests",
			Position: Position{},
		})
	}

	return errors
}

// PerformanceValidationRule validates performance characteristics
type PerformanceValidationRule struct{}

func (r *PerformanceValidationRule) Name() string {
	return "PerformanceValidation"
}

func (r *PerformanceValidationRule) Validate(pattern *Pattern) []PatternError {
	var errors []PatternError

	// Check pattern complexity
	if len(pattern.Values) > 1000 {
		errors = append(errors, PatternError{
			Message: fmt.Sprintf("pattern is too complex (%d values, maximum 1000)", len(pattern.Values)),
			Position: Position{},
		})
	}

	// Check duration
	if pattern.Duration > 60*time.Second {
		errors = append(errors, PatternError{
			Message: "pattern duration exceeds 60 seconds",
			Position: Position{},
		})
	}

	return errors
}

// MusicalValidationRule validates musical properties
type MusicalValidationRule struct{}

func (r *MusicalValidationRule) Name() string {
	return "MusicalValidation"
}

func (r *MusicalValidationRule) Validate(pattern *Pattern) []PatternError {
	var errors []PatternError

	// Check for valid musical notes
	for _, value := range pattern.Values {
		if note, ok := value.(NoteValue); ok {
			// Validate note name
			validNotes := map[string]bool{
				"c": true, "d": true, "e": true, "f": true, 
				"g": true, "a": true, "b": true,
			}
			if !validNotes[note.Note] {
				errors = append(errors, PatternError{
					Message: fmt.Sprintf("invalid note name: %s", note.Note),
					Position: Position{},
				})
			}

			// Validate accidental
			if note.Accident != "" && note.Accident != "#" && note.Accident != "b" {
				errors = append(errors, PatternError{
					Message: fmt.Sprintf("invalid accidental: %s", note.Accident),
					Position: Position{},
				})
			}
		}
	}

	return errors
}

// PatternOptimizer optimizes patterns for better performance
type PatternOptimizer struct {
	enabled bool
}

// NewPatternOptimizer creates a new pattern optimizer
func NewPatternOptimizer() *PatternOptimizer {
	return &PatternOptimizer{enabled: true}
}

// Optimize optimizes a pattern
func (o *PatternOptimizer) Optimize(pattern *Pattern) *Pattern {
	if !o.enabled {
		return pattern
	}

	// Create a new optimized pattern
	optimized := &Pattern{
		Values:   make([]PatternValue, 0, len(pattern.Values)),
		Duration: pattern.Duration,
		Events:   make([]PatternEvent, 0, len(pattern.Events)),
		Metadata: make(map[string]interface{}),
	}

	// Copy metadata
	for k, v := range pattern.Metadata {
		optimized.Metadata[k] = v
	}

	// Optimize values
	optimized.Values = o.optimizeValues(pattern.Values)

	// Optimize events
	optimized.Events = o.optimizeEvents(pattern.Events)

	return optimized
}

// optimizeValues optimizes pattern values
func (o *PatternOptimizer) optimizeValues(values []PatternValue) []PatternValue {
	optimized := make([]PatternValue, 0, len(values))

	// Remove consecutive rests
	lastWasRest := false
	for _, value := range values {
		if _, isRest := value.(RestValue); isRest {
			if !lastWasRest {
				optimized = append(optimized, value)
				lastWasRest = true
			}
		} else {
			optimized = append(optimized, value)
			lastWasRest = false
		}
	}

	return optimized
}

// optimizeEvents optimizes pattern events
func (o *PatternOptimizer) optimizeEvents(events []PatternEvent) []PatternEvent {
	optimized := make([]PatternEvent, 0, len(events))

	// Remove duplicate events
	seen := make(map[string]bool)
	for _, event := range events {
		key := fmt.Sprintf("%v-%v-%v", event.Value, event.Time, event.Duration)
		if !seen[key] {
			optimized = append(optimized, event)
			seen[key] = true
		}
	}

	return optimized
}

// ValidationContext provides context for validation
type ValidationContext struct {
	AllowMixedTypes bool
	MaxPatternSize  int
	MaxDuration     time.Duration
	StrictMode      bool
}

// NewValidationContext creates a new validation context
func NewValidationContext() *ValidationContext {
	return &ValidationContext{
		AllowMixedTypes: false,
		MaxPatternSize:  1000,
		MaxDuration:     60 * time.Second,
		StrictMode:      false,
	}
}

// ContextAwareValidator validates patterns based on context
type ContextAwareValidator struct {
	context *ValidationContext
	rules   []ValidationRule
}

// NewContextAwareValidator creates a new context-aware validator
func NewContextAwareValidator(context *ValidationContext) *ContextAwareValidator {
	v := &ContextAwareValidator{
		context: context,
		rules:   make([]ValidationRule, 0),
	}

	// Add rules based on context
	if context.StrictMode {
		v.AddRule(&StrictValidationRule{})
	}

	return v
}

// AddRule adds a validation rule
func (v *ContextAwareValidator) AddRule(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

// Validate validates a pattern with context
func (v *ContextAwareValidator) Validate(pattern *Pattern) []PatternError {
	var errors []PatternError

	// Apply context-specific rules
	if !v.context.AllowMixedTypes {
		errors = append(errors, v.validateMixedTypes(pattern)...)
	}

	if len(pattern.Values) > v.context.MaxPatternSize {
		errors = append(errors, PatternError{
			Message: fmt.Sprintf("pattern size %d exceeds maximum %d", len(pattern.Values), v.context.MaxPatternSize),
			Position: Position{},
		})
	}

	if pattern.Duration > v.context.MaxDuration {
		errors = append(errors, PatternError{
			Message: fmt.Sprintf("pattern duration %v exceeds maximum %v", pattern.Duration, v.context.MaxDuration),
			Position: Position{},
		})
	}

	// Apply custom rules
	for _, rule := range v.rules {
		ruleErrors := rule.Validate(pattern)
		errors = append(errors, ruleErrors...)
	}

	return errors
}

// validateMixedTypes validates that patterns don't mix incompatible types
func (v *ContextAwareValidator) validateMixedTypes(pattern *Pattern) []PatternError {
	var errors []PatternError

	hasSamples := false
	hasNotes := false

	for _, value := range pattern.Values {
		switch value.(type) {
		case SampleValue:
			hasSamples = true
		case NoteValue:
			hasNotes = true
		}
	}

	if hasSamples && hasNotes {
		errors = append(errors, PatternError{
			Message: "pattern cannot mix samples and notes in strict mode",
			Position: Position{},
		})
	}

	return errors
}

// StrictValidationRule implements strict validation rules
type StrictValidationRule struct{}

func (r *StrictValidationRule) Name() string {
	return "StrictValidation"
}

func (r *StrictValidationRule) Validate(pattern *Pattern) []PatternError {
	var errors []PatternError

	// Strict validation rules
	for _, value := range pattern.Values {
		switch v := value.(type) {
		case SampleValue:
			// Validate sample format
			if len(v.Name) == 0 {
				errors = append(errors, PatternError{
					Message: "sample name cannot be empty",
					Position: Position{},
				})
			}
		case NoteValue:
			// Validate note range
			if v.Octave < 1 || v.Octave > 8 {
				errors = append(errors, PatternError{
					Message: fmt.Sprintf("octave %d out of range (1-8) in strict mode", v.Octave),
					Position: Position{},
				})
			}
		}
	}

	return errors
}