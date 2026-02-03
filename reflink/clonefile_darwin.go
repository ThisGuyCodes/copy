//go:build darwin

package reflink

import (
	"os"

	"golang.org/x/sys/unix"
)

func clonefile(from *os.File, toDir *os.File, toName string) error {
	fromFD := int(from.Fd())
	toDirFD := int(toDir.Fd())
	return unix.Fclonefileat(fromFD, toDirFD, toName, unix.CLONE_NOFOLLOW|unix.CLONE_NOOWNERCOPY)
}
