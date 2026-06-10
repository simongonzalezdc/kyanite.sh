package errors

import "github.com/Kyanite/noise/internal/logging"

// Result is a generic result type that can hold either a value or an error
type Result[T any] struct {
	value T
	err   error
}

// Ok creates a successful result
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, err: nil}
}

// Err creates an error result
func Err[T any](err error) Result[T] {
	var zero T
	return Result[T]{value: zero, err: err}
}

// IsOk returns true if the result is successful
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the result is an error
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Unwrap returns the value or panics if error
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		logging.Errorf("Panic: called Unwrap on error result: %v", r.err)
		panic("called Unwrap on error result")
	}
	return r.value
}

// UnwrapOr returns the value or a default if error
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

// UnwrapOrElse returns the value or calls the provided function if error
func (r Result[T]) UnwrapOrElse(f func() T) T {
	if r.err != nil {
		return f()
	}
	return r.value
}

// UnwrapOrZero returns the value or the zero value of type T if error
func (r Result[T]) UnwrapOrZero() T {
	if r.err != nil {
		var zero T
		return zero
	}
	return r.value
}

// Expect returns the value or panics with message if error
func (r Result[T]) Expect(message string) T {
	if r.err != nil {
		logging.Errorf("Panic: %s: %v", message, r.err)
		panic(message + ": " + r.err.Error())
	}
	return r.value
}

// Error returns the error or nil
func (r Result[T]) Error() error {
	return r.err
}

// Map transforms the value if successful
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.err != nil {
		return r
	}
	return Ok(f(r.value))
}

// MapErr transforms the error if present
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.err == nil {
		return r
	}
	return Err[T](f(r.err))
}

// AndThen chains results
func (r Result[T]) AndThen(f func(T) Result[T]) Result[T] {
	if r.err != nil {
		return r
	}
	return f(r.value)
}

// Option is a generic optional type
type Option[T any] struct {
	value   T
	present bool
}

// Some creates an option with a value
func Some[T any](value T) Option[T] {
	return Option[T]{value: value, present: true}
}

// None creates an empty option
func None[T any]() Option[T] {
	var zero T
	return Option[T]{value: zero, present: false}
}

// IsSome returns true if the option has a value
func (o Option[T]) IsSome() bool {
	return o.present
}

// IsNone returns true if the option is empty
func (o Option[T]) IsNone() bool {
	return !o.present
}

// Unwrap returns the value or panics if none
func (o Option[T]) Unwrap() T {
	if !o.present {
		panic("called Unwrap on None")
	}
	return o.value
}

// UnwrapOr returns the value or a default if none
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if !o.present {
		return defaultValue
	}
	return o.value
}

// UnwrapOrElse returns the value or calls the provided function if none
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if !o.present {
		return f()
	}
	return o.value
}

// UnwrapOrZero returns the value or the zero value of type T if none
func (o Option[T]) UnwrapOrZero() T {
	if !o.present {
		var zero T
		return zero
	}
	return o.value
}

// Map transforms the value if present
func (o Option[T]) Map(f func(T) T) Option[T] {
	if !o.present {
		return o
	}
	return Some(f(o.value))
}
