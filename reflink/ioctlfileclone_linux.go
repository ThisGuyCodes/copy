//go:build linux

package reflink

import (
	"os"

	"golang.org/x/sys/unix"
)

func ioctlFileClone(from *os.File, toDir *os.File, toName string) error {
	fromFD := int(from.Fd())
	toDirFD := int(toDir.Fd())

	toFD, err := unix.Openat(toDirFD, toName, unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL, 0644)
	if err != nil {
		return err
	}

	doDeferClose := true
	defer func() {
		if doDeferClose {
			unix.Close(toFD)
		}
	}()

	err = unix.IoctlFileClone(fromFD, toFD)
	if err != nil {
		return err
	}

	doDeferClose = false
	return unix.Close(toFD)
}
