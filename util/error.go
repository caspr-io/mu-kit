package util

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
