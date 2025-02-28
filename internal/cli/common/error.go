package common

import (
	"errors"
	"strings"
)

var ErrInvalidElement = errors.New("invalid element")

type ValidateError struct {
	errors []string
}

func NewValidateError() *ValidateError {
	return &ValidateError{}
}

func (e *ValidateError) AddError(msg string) {
	e.errors = append(e.errors, msg)
}

func (e *ValidateError) HasErrors() bool {
	return len(e.errors) > 0
}

func (e *ValidateError) Error() string {
	return strings.Join(e.errors, "\n")
}

type FormError struct {
	text string
}

func NewFormError(text string) *FormError {
	return &FormError{
		text: text,
	}
}

func (e *FormError) Error() string {
	return e.text
}
