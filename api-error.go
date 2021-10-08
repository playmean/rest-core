package restcore

import (
	"fmt"
	"runtime/debug"
)

type ApiError struct {
	code     string
	message  string
	stack    string
	original error
}

type ApiErrorOptions struct {
	Code     string
	Subcode  string
	Message  string
	Original error
}

func NewApiError(opts *ApiErrorOptions) error {
	e := &ApiError{
		code:     opts.Code,
		message:  opts.Message,
		stack:    string(debug.Stack()),
		original: opts.Original,
	}

	if opts.Subcode != "" {
		e.code += fmt.Sprintf("[%s]", opts.Subcode)
	}

	return e
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("[error] %s - %s", e.code, e.message)
}

func (e *ApiError) Code() string {
	return e.code
}

func (e ApiError) Message() string {
	return e.message
}

func (e ApiError) Stack() string {
	return e.stack
}

func (e ApiError) Original() error {
	return e.original
}
