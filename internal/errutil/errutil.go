package errutil

import (
	"errors"
	"fmt"
	"strings"
)

// Wrap ensures a consistent wrapped error message for operational failures.
// It returns nil when the supplied error is nil to avoid redundant wrapping.
func Wrap(err error, operation string) error {
	if err == nil {
		return nil
	}

	op := strings.TrimSpace(operation)
	if op == "" {
		op = "operation"
	}

	return fmt.Errorf("%s failed: %w", op, err)
}

// WrapContext wraps an error with additional context while preserving the cause.
// When err is nil the function returns nil to maintain existing behaviour.
func WrapContext(err error, message string) error {
	if err == nil {
		return nil
	}

	msg := strings.TrimSpace(message)
	if msg == "" {
		msg = "operation failed"
	}

	return fmt.Errorf("%s: %w", msg, err)
}

// Wrapf formats a contextual message before wrapping the supplied error.
// It panics if the format string does not contain the required verbs.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return WrapContext(err, fmt.Sprintf(format, args...))
}

// Validation produces a validation error using the standardised message pattern.
func Validation(field string) error {
	fieldName := strings.TrimSpace(field)
	if fieldName == "" {
		fieldName = "field"
	}

	return fmt.Errorf("validation failed: %s", fieldName)
}

// Join normalises multi-error aggregation by removing nil entries before joining.
func Join(errs ...error) error {
	filtered := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			filtered = append(filtered, err)
		}
	}

	switch len(filtered) {
	case 0:
		return nil
	case 1:
		return filtered[0]
	default:
		return errors.Join(filtered...)
	}
}