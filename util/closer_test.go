package util

import "testing"

import "gotest.tools/v3/assert"

import "fmt"

type TestCloser struct {
	e error
}

func (t *TestCloser) Close() error {
	return t.e
}

func TestMultiCloserShouldReturnNilErrorWhenNoErrorsClosing(t *testing.T) {
	multi := new(MultiCloser)
	multi.Add(&TestCloser{e: nil})

	assert.NilError(t, multi.Close())
}

func TestMultiCloserShouldReturnErrorWhenErrorDuringClosing(t *testing.T) {
	multi := new(MultiCloser)
	multi.Add(&TestCloser{e: fmt.Errorf("An Error")})

	assert.ErrorContains(t, multi.Close(), "An Error")
}
