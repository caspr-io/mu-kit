package util

import (
	"fmt"
)

type ErrorCollector []error

func (c *ErrorCollector) Collect(e error) {
	if e == nil {
		return
	}

	switch e := e.(type) {
	case *ErrorCollector:
		if e.HasErrors() {
			for _, err := range *e { //nolint:gosimple (the simplification results in the Collector itself being appended, instead of appending its errors)
				*c = append(*c, err)
			}
		}
	default:
		*c = append(*c, e)
	}
}

func (c *ErrorCollector) Error() (err string) {
	err = "Collected errors:\n"
	for i, e := range *c {
		err += fmt.Sprintf("\tError %d: %s\n", i, e.Error())
	}

	return err
}

func (c *ErrorCollector) HasErrors() bool {
	return len(*c) != 0
}

type TryCatch struct {
	err error
}

func Try() *TryCatch {
	return &TryCatch{nil}
}

func (cb *TryCatch) Try(f func() error) *TryCatch {
	if cb.err != nil {
		return cb
	}

	cb.err = f()

	return cb
}

func (cb *TryCatch) Caught() bool {
	return cb.err != nil
}

func (cb *TryCatch) Error() error {
	return cb.err
}
