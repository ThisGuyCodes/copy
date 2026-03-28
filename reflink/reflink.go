package reflink

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func Reflink(from *os.File, toDir *os.File, toName string) error {
	// indirection is so we can add other platform specific options later
	err := clonefile(from, toDir, toName)
	if err == nil || err != ErrNotOnPlatform {
		return err
	}

	return ioctlFileClone(from, toDir, toName)
}

func ReflinkOrCopy(from *os.File, toDir *os.File, toName string) (bool, error) {
	wasReflinked := true
	err := Reflink(from, toDir, toName)
	if err != nil {
		wasReflinked = false
	}
	if err == nil || err != ErrNotOnPlatform {
		return wasReflinked, err
	}

	fromPerms, err := from.Stat()
	if err != nil {
		return wasReflinked, err
	}

	toFile, err := createFile(toDir, toName, fromPerms.Mode())
	if err != nil {
		return wasReflinked, err
	}

	doDeferClose := true
	defer func() {
		if doDeferClose {
			toFile.Close() // nolint:errcheck
		}
	}()

	// on linux Go automatically uses copy_file_range, which internally will
	// use reflink if the file system supports it
	_, copyErr := io.Copy(toFile, from)
	doDeferClose = false
	closeErr := toFile.Close()

	return wasReflinked, errors.Join(copyErr, closeErr)
}

func ReflinkOrCopyAfero(fs afero.Fs, from, to string) (wasReflinked bool, joinErr error) {
	var fromFileCloseErr error
	var toFileCloseErr error
	var toDirCloseErr error
	var runningErr error

	defer func() {
		joinErr = errors.Join(fromFileCloseErr, toFileCloseErr, toDirCloseErr, runningErr)
	}()

	fromFile, runningErr := fs.Open(from)
	if runningErr != nil {
		return
	}
	defer func() {
		fromFileCloseErr = fromFile.Close()
	}()

	toDir := filepath.Dir(to)
	toDirFile, runningErr := fs.Open(toDir)
	if runningErr != nil {
		return
	}
	defer func() {
		toDirCloseErr = toDirFile.Close()
	}()

	fromOSFile, fromIsOs := fromFile.(*os.File)
	toDirOSFile, toDirIsOs := toDirFile.(*os.File)

	if fromIsOs && toDirIsOs {
		wasReflinked, runningErr = ReflinkOrCopy(fromOSFile, toDirOSFile, to)
		return
	}

	fullToName := filepath.Join(toDirFile.Name(), to)
	toFile, runningErr := fs.OpenFile(fullToName, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if runningErr != nil {
		return
	}
	defer func() {
		toFileCloseErr = toFile.Close()
	}()

	_, runningErr = io.Copy(toFile, fromFile)

	return
}
