package copy

import (
	"os"
	"path/filepath"

	"github.com/thisguycodes/copy/reflink"
)

func Copy(from, to string) error {
	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toDir := filepath.Dir(to)
	toDirFile, err := os.Open(toDir)
	if err != nil {
		return err
	}
	defer toDirFile.Close()

	toFile := filepath.Base(to)

	_, err = reflink.ReflinkOrCopy(fromFile, toDirFile, toFile)
	return err
}
