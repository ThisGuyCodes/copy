package createmount

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// MountDiskImageMacOS creates a disk image with the specified filesystem type,
// mounts it, and registers cleanup on the provided testing.TB.
func MountDiskImageMacOS(t testing.TB, mountpoint, fsType string) {
	t.Helper()

	if runtime.GOOS != "darwin" {
		t.Fatalf("this only works on macOS")
	}

	// Create the disk image file path in the temp directory managed by TB
	imagePath := filepath.Join(t.TempDir(), "test_image.sparseimage")

	// Calculate size for the image (e.g., 500MB)
	const size = "512m"

	// Create the disk image file
	cmd := exec.Command("hdiutil", "create", "-size", size, "-fs", fsType, "-volname", "TestVolume", "-type", "SPARSE", "-quiet", imagePath)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to create disk image: %v, output: %s", err, string(out))
	}

	// Mount the disk image to the provided location
	mountCmd := exec.Command("hdiutil", "attach", "-mountPoint", mountpoint, imagePath)
	if err := mountCmd.Run(); err != nil {
		t.Fatalf("failed to mount disk image: %v", err)
	}

	// Register cleanup to unmount the disk after the test
	t.Cleanup(func() {
		unmountCmd := exec.Command("hdiutil", "detach", mountpoint)
		if out, err := unmountCmd.CombinedOutput(); err != nil {
			t.Errorf("failed to unmount disk image: %v, output: %s", err, string(out))
		}
	})
}
