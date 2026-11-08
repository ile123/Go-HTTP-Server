package constants

type ErrorKind string

const (
	ErrMissingField ErrorKind = "x-missing-field"
	ErrInvalidValue ErrorKind = "x-invalid-value"
	ErrTooLong      ErrorKind = "x-value-too-long"
)

type ResponseError struct {
	Kind  ErrorKind
	Value string
}
