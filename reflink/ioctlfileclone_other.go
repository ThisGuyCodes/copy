//go:build !linux

package reflink

import (
	"os"
)

func ioctlFileClone(from *os.File, toDir *os.File, toName string) error {
	return ErrNotOnPlatform
}
