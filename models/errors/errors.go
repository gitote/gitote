package errors

import "errors"

var InternalServerError = errors.New("internal server error")

// New is a wrapper of real errors.New function.
func New(text string) error {
	return errors.New(text)
}
