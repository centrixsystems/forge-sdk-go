package forge

import "fmt"

// ServerError is returned when the server responds with a 4xx/5xx status.
type ServerError struct {
	StatusCode int
	Message    string
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("forge: server error (%d): %s", e.StatusCode, e.Message)
}

// ConnectionError is returned when the HTTP request fails.
type ConnectionError struct {
	Cause error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("forge: connection error: %v", e.Cause)
}

func (e *ConnectionError) Unwrap() error {
	return e.Cause
}
