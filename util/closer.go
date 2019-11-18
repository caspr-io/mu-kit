package util

import "io"

type MultiCloser []io.Closer

func (mc *MultiCloser) Add(c io.Closer) {
	if c != nil {
		*mc = append(*mc, c)
	}
}

func (mc *MultiCloser) Close() error {
	errorCollector := new(ErrorCollector)

	for i := len(*mc) - 1; i >= 0; i-- {
		c := (*mc)[i]
		if err := c.Close(); err != nil {
			errorCollector.Collect(err)
		}
	}

	if errorCollector.HasErrors() {
		return errorCollector
	}

	return nil
}
