package process

import "errors"

var (
	errInnerDBIsNotOpened    = errors.New("inner DB is not opened")
	errInnerDBIsNotClosed    = errors.New("inner DB is not closed")
	errNilDirectoriesHandler = errors.New("nil directories handler instance")
	errNilDBWrapper          = errors.New("nil DB wrapper instance")
)
