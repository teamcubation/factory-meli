package errors

import "errors"

var (
	ErrContentTooLong = errors.New("content must be 280 characters or less")
	ErrNotFound       = errors.New("the resource you tried to find doesn't exist")
)
