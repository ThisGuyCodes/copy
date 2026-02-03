//go:build !(darwin || linux)

package reflink

import (
	"os"
	"path/filepath"
)

func createFile(dir *os.File, name string) (*os.File, error) {
	fileName := filepath.Join(dir.Name(), name)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}
