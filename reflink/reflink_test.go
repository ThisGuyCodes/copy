package reflink_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/thisguycodes/copy/reflink"
	"github.com/thisguycodes/copy/reflink/testutils/createmount"
	"github.com/thisguycodes/copy/reflink/testutils/ts"
)

func TestReflinkOnDarwinWithinAPFS(t *testing.T) {
	ts.OnlyOn(t, "darwin_")
	t.Parallel()

	apfsMount := t.TempDir()
	createmount.MountDiskImageMacOS(t, apfsMount, "APFS")

	fileName := filepath.Join(apfsMount, "test.txt")

	ts.NoErr(0, os.WriteFile(fileName, []byte("Hello, World!"), 0o644))(t)

	fromFD := ts.NoErr(os.Open(fileName))(t)
	toDirFD := ts.NoErr(os.Open(apfsMount))(t)
	defer fromFD.Close()
	defer toDirFD.Close()

	toName := "test-reflink.txt"

	ts.NoErr(0, reflink.Reflink(fromFD, toDirFD, toName))
}

func TestReflinkOnDarwinAcrossAPFS(t *testing.T) {
	ts.OnlyOn(t, "darwin_")
	t.Parallel()

	apfsMount1 := t.TempDir()
	apfsMount2 := t.TempDir()

	createmount.MountDiskImageMacOS(t, apfsMount1, "APFS")
	createmount.MountDiskImageMacOS(t, apfsMount2, "APFS")

	fileName := filepath.Join(apfsMount1, "test.txt")

	ts.NoErr(0, os.WriteFile(fileName, []byte("Hello, World!"), 0o644))(t)

	fromFD := ts.NoErr(os.Open(fileName))(t)
	toDirFD := ts.NoErr(os.Open(apfsMount2))(t)
	defer fromFD.Close()
	defer toDirFD.Close()

	toName := "test-reflink.txt"

	err := reflink.Reflink(fromFD, toDirFD, toName)
	ts.True(errors.Is(err, reflink.ErrCanNotReflink{}))(t)
}

func TestReflinkOnDarwinWithinExFAT(t *testing.T) {
	ts.OnlyOn(t, "darwin_")
	t.Parallel()

	exfatMount := t.TempDir()
	createmount.MountDiskImageMacOS(t, exfatMount, "ExFAT")

	fileName := filepath.Join(exfatMount, "test.txt")

	ts.NoErr(0, os.WriteFile(fileName, []byte("Hello, World!"), 0o644))(t)

	fromFD := ts.NoErr(os.Open(fileName))(t)
	toDirFD := ts.NoErr(os.Open(exfatMount))(t)
	defer fromFD.Close()
	defer toDirFD.Close()

	toName := "test-reflink.txt"

	err := reflink.Reflink(fromFD, toDirFD, toName)
	ts.True(errors.Is(err, reflink.ErrCanNotReflink{}))(t)
}
