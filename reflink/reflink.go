package reflink

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func Reflink(from *os.File, toDir *os.File, toName string) error {
	err := clonefile(from, toDir, toName)
	if err == nil || err != ErrNotOnPlatform {
		return err
	}
	err = ioctlFileClone(from, toDir, toName)
	if err == nil || err != ErrNotOnPlatform {
		return err
	}
	return ErrNotOnPlatform
}

func ReflinkOrCopy(from *os.File, toDir *os.File, toName string) error {
	err := Reflink(from, toDir, toName)
	if err == nil || err != ErrNotOnPlatform {
		return err
	}

	toFile, err := createFile(toDir, toName)
	doDeferClose := true
	defer func() {
		if doDeferClose {
			toFile.Close()
		}
	}()

	_, copyErr := io.Copy(toFile, from)
	doDeferClose = false
	closeErr := toFile.Close()

	return errors.Join(copyErr, closeErr)
}

func ReflinkOrCopyAfero(fs afero.Fs, from, toName string) (joinErr error) {
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

	toDir := filepath.Dir(toName)
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
		runningErr = ReflinkOrCopy(fromOSFile, toDirOSFile, toName)
		return
	}

	fullToName := filepath.Join(toDirFile.Name(), toName)
	toFile, runningErr := fs.OpenFile(fullToName, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if runningErr != nil {
		return
	}
	defer func() {
		toFileCloseErr = toFile.Close()
	}()

	_, runningErr = io.Copy(toFile, fromFile)

	return
}
