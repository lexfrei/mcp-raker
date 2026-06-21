// Package moonraker is a client for the Moonraker HTTP API (the web server that
// fronts the Klipper 3D printer firmware). It talks to the documented REST
// endpoints, unwraps the {"result": ...} envelope, and supports API-key,
// pre-obtained Bearer, and username/password (JWT) authentication.
package moonraker

import "github.com/cockroachdb/errors"

// ErrNoCredentials indicates a JWT login was required but no username/password
// was configured.
var ErrNoCredentials = errors.New("no credentials configured")

// ErrLoginFailed indicates the /access/login endpoint rejected the credentials.
var ErrLoginFailed = errors.New("login failed")

// ErrNotAuthenticated indicates Moonraker rejected the request because the
// API key or token is missing, invalid, or expired.
var ErrNotAuthenticated = errors.New("not authenticated")

// ErrAPI indicates Moonraker returned an error response or an unexpected body.
var ErrAPI = errors.New("moonraker API error")

// ErrNotFound indicates Moonraker returned HTTP 404. It wraps ErrAPI, so it
// satisfies errors.Is(err, ErrAPI) too. Callers can use it to degrade gracefully
// when an optional component (a sensor, power device, or WLED strip) is simply
// not configured.
var ErrNotFound = errors.Wrap(ErrAPI, "resource not found")

// apiErr wraps err with a message and marks it as an API failure.
func apiErr(err error, format string, args ...any) error {
	//nolint:wrapcheck // Mark only adds a sentinel category; Wrapf already added the message.
	return errors.Mark(errors.Wrapf(err, format, args...), ErrAPI)
}

// loginErr wraps err with a message and marks it as a login failure.
func loginErr(err error, msg string) error {
	//nolint:wrapcheck // Mark only adds a sentinel category; Wrap already added the message.
	return errors.Mark(errors.Wrap(err, msg), ErrLoginFailed)
}
