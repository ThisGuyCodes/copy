package reflink

import (
	"errors"
)

var ErrNotOnPlatform = errors.New("this function is not available on this platform")

type ErrCanNotReflink struct {
	wrapped error
}

func newErrCanNotReflink(underlying error) ErrCanNotReflink {
	return ErrCanNotReflink{wrapped: underlying}
}

func (nr ErrCanNotReflink) Error() string {
	return "Reflink doesn't work here"
}

func (nr ErrCanNotReflink) Unwrap() error {
	return nr.wrapped
}

func (nr ErrCanNotReflink) Is(err error) bool {
	switch err.(type) {
	case ErrCanNotReflink:
		return true
	default:
		return false
	}
}
