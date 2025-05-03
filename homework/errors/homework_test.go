package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.errs)))

	for _, err := range e.errs {
		builder.WriteString(fmt.Sprintf("\t* %s", err.Error()))
	}
	builder.WriteString("\n")

	return builder.String()
}

func Append(err error, errs ...error) *MultiError {
	multierr, ok := err.(*MultiError)
	if !ok {
		return &MultiError{
			errs: errs,
		}
	}
	multierr.errs = append(multierr.errs, errs...)
	return multierr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
