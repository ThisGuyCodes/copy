package ts

import (
	"errors"
	"syscall"
	"testing"
)

func NoErr[T any](ret T, err error) func(t testing.TB) T {
	return func(t testing.TB) T {
		t.Helper()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return ret
	}
}

func True(ret bool) func(t testing.TB) bool {
	return func(t testing.TB) bool {
		t.Helper()
		if !ret {
			t.Fatalf("expected true, got false")
		}
		return ret
	}
}

func IsErr(t testing.TB, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		if err, ok := err.(syscall.Errno); ok {
			t.Logf("hint: received syscall error %v", uintptr(err))
		}
		t.Fatalf("expected error %q, got %q", target, err)
	}
}
