package common

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation fail")
)

type entityNotFound struct {
	Entity string
	Err    error
}

func (e *entityNotFound) Error() string {
	return fmt.Sprintf("%s, reason: %s", e.Entity, e.Err.Error())
}

type businessFail struct {
	Message string
	Err     error
}

func (e *businessFail) Error() string {
	if e.Err == nil {
		e.Err = ErrValidation
	}
	return fmt.Sprintf("%s: %s", e.Err.Error(), e.Message)
}

func NotFound(entity string, cause error) error {
	return &entityNotFound{entity, cause}
}

func BusinessFailure(msg string, cause error) error {
	return &businessFail{msg, cause}
}
