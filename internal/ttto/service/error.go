package service

import (
	"encoding/json"
)

// ErrValidation used to wrap validation error
// to hide validation logic from server layer
type ErrValidation struct {
	err error
}

func newErrValidation(err error) error {
	return &ErrValidation{err: err}
}

func (e *ErrValidation) Error() string {
	return e.err.Error()
}

func (e *ErrValidation) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.err)
}

func (e *ErrValidation) Unwrap() error {
	return e.err
}
