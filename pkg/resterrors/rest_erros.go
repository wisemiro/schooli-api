package resterrors

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	ErrorInvalidJSONBody      string = "invalid json body"
	ErrorParsingQuery         string = "error parsing query"
	ErrorPasswordMisMatch     string = "password and password re-enter did not match"
	ErrorAccountNotFound      string = "account not found"
	ErrorItemNotFound         string = "item not found"
	ErrorInvalidToken         string = "invalid token"
	ErrorInvalidEmailPassword string = "invalid email or password"
	ErrorProcessingRequest    string = "error processing request"
	ErrorEmailExists          string = "email already exists"
	ErrorPhoneExists          string = "phone number exists"
)

type RestErr struct {
	TimeStamp time.Time `json:"time_stamp"`
	Message   string    `json:"message"`
	Status    int       `json:"status"`
	Error     string    `json:"error"`
}

func NewError(msg string) error {
	return errors.New(msg)
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		TimeStamp: time.Now(),
		Message:   message,
		Status:    http.StatusBadRequest,
		Error:     "bad_request",
	}
}
func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		TimeStamp: time.Now(),
		Message:   message,
		Status:    http.StatusInternalServerError,
		Error:     "internal_server_error",
	}
}
func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		TimeStamp: time.Now(),
		Message:   message,
		Status:    http.StatusNotFound,
		Error:     "not_found",
	}
}

// Error represents an error that could be wrapping another error, it includes a code for determining what
// triggered the error.
type Error struct {
	orig error
	msg  string
	code ECode
}

// ErrorCode defines supported error codes.
type ECode uint

const (
	ECodeUnknown ECode = iota
	ECodeNotFound
	ECodeInvalidArgument
)

// WrapErrorf returns a wrapped error.
func WrapErrorf(orig error, code ECode, format string, a ...any) error {
	return &Error{
		code: code,
		orig: orig,
		msg:  fmt.Sprintf(format, a...),
	}
}

// NewErrorf instantiates a new error.
func NewErrorf(code ECode, format string, a ...interface{}) error {
	return WrapErrorf(nil, code, format, a...)
}

// Error returns the message, when wrapping errors the wrapped error is returned.
func (e *Error) Error() string {
	if e.orig != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.orig)
	}

	return e.msg
}

// Unwrap returns the wrapped error, if any.
func (e *Error) Unwrap() error {
	return e.orig
}

// Code returns the code representing this error.
func (e *Error) Code() ECode {
	return e.code
}
