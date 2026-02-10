package ts

import "testing"

func NoErr[T any](ret T, err error) func(t testing.TB) T {
	return func(t testing.TB) T {
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		return ret
	}
}
