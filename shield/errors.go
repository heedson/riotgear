package shield

import "github.com/pkg/errors"

// An ErrorCode is a code to help allow internal errors to be translated to gRPC errors.
type ErrorCode uint32

const (
	// Unknown is for when the reason for the error is unknown.
	Unknown ErrorCode = 0
	// Internal is for when the reason for the error is an internal service issue.
	Internal ErrorCode = 1
	// InvalidArgument is for when the provided argument has an invalid value.
	InvalidArgument ErrorCode = 2
)

// Error returns a Gear wrapped error with a stack trace.
func Error(errorCode ErrorCode, message string) error {
	switch errorCode {
	case Internal:
		return ErrInternal{Err: errors.New(message)}
	case InvalidArgument:
		return ErrInvalidArgument{Err: errors.New(message)}
	default:
		return ErrUnknown{Err: errors.New(message)}
	}
}

// Errorf returns a Gear wrapped error with a stack trace.
func Errorf(errorCode ErrorCode, format string, args ...interface{}) error {
	switch errorCode {
	case Internal:
		return ErrInternal{Err: errors.Errorf(format, args...)}
	case InvalidArgument:
		return ErrInvalidArgument{Err: errors.Errorf(format, args...)}
	default:
		return ErrUnknown{Err: errors.Errorf(format, args...)}
	}
}

// Wrap returns a Gear wrapped error with a stack trace.
func Wrap(errorCode ErrorCode, err error) error {
	switch errorCode {
	case Internal:
		return ErrInternal{Err: errors.WithStack(err)}
	case InvalidArgument:
		return ErrInvalidArgument{Err: errors.WithStack(err)}
	default:
		return ErrUnknown{Err: errors.WithStack(err)}
	}
}

// ErrUnknown is returned when an errors reason is unknown.
type ErrUnknown struct {
	Err error
}

func (e ErrUnknown) Error() string {
	return e.Err.Error()
}

// ErrInternal is returned when an internal error occurs.
type ErrInternal struct {
	Err error
}

func (e ErrInternal) Error() string {
	return e.Err.Error()
}

// ErrInvalidArgument is returned when a supplied argument has an invalid value.
type ErrInvalidArgument struct {
	Err error
}

func (e ErrInvalidArgument) Error() string {
	return e.Err.Error()
}
