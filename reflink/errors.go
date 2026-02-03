package reflink

import (
	"errors"
)

var (
	ErrNotOnPlatform = errors.New("this function is not available on this platform")
)
