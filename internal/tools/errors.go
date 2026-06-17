// Package tools provides MCP tool definitions and handlers for the Moonraker
// API. Each tool wraps a single Moonraker endpoint.
package tools

import "github.com/cockroachdb/errors"

// ErrValidation indicates invalid parameters provided by the caller.
var ErrValidation = errors.New("validation error")

// ErrMoonraker indicates a failure talking to the Moonraker API.
var ErrMoonraker = errors.New("moonraker request error")

// ErrFieldRequired indicates a required parameter was missing or invalid.
var ErrFieldRequired = errors.New("required parameter missing")

// ErrMutuallyExclusive indicates two parameters were set that cannot combine.
var ErrMutuallyExclusive = errors.New("parameters are mutually exclusive")

// ErrOutOfRange indicates a numeric parameter fell outside its allowed range.
var ErrOutOfRange = errors.New("parameter out of range")

// validationErr marks an error as a validation error.
func validationErr(err error) error {
	//nolint:wrapcheck // Mark adds a sentinel category; the caller supplies the message.
	return errors.Mark(err, ErrValidation)
}

// moonrakerErr wraps a message and underlying error as a Moonraker API error.
func moonrakerErr(msg string, err error) error {
	//nolint:wrapcheck // Mark adds a sentinel category on top of Wrap which adds context.
	return errors.Mark(errors.Wrap(err, msg), ErrMoonraker)
}

// requireString returns a validation error when value is empty.
func requireString(field, value string) error {
	if value == "" {
		return validationErr(errors.Wrapf(ErrFieldRequired, "%s is required", field))
	}

	return nil
}

// requirePositive returns a validation error when value is not a positive int.
func requirePositive(field string, value int) error {
	if value <= 0 {
		return validationErr(errors.Wrapf(ErrFieldRequired, "%s must be a positive integer", field))
	}

	return nil
}

// requirePresent returns a validation error when count is not positive; it is
// used to assert that a map or slice parameter is non-empty.
func requirePresent(field string, count int) error {
	if count <= 0 {
		return validationErr(errors.Wrapf(ErrFieldRequired, "%s is required", field))
	}

	return nil
}

// mutuallyExclusive returns a validation error stating that two parameters
// cannot be combined.
func mutuallyExclusive(first, second string) error {
	return validationErr(errors.Wrapf(ErrMutuallyExclusive, "set %s or %s, not both", first, second))
}

// requireRange returns a validation error when value is outside [low, high].
func requireRange(field string, value, low, high int) error {
	if value < low || value > high {
		return validationErr(errors.Wrapf(ErrOutOfRange, "%s must be between %d and %d", field, low, high))
	}

	return nil
}
