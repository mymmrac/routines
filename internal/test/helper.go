package test

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Fatalf("%+v is not equal to %+v", actual, expected)
	}
}

func EqualEl[T comparable, S []T](t *testing.T, actual, expected S) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("%+v is not equal to %+v", actual, expected)
	}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("%+v is not equal to %+v", actual, expected)
		}
	}
}

func True(t *testing.T, actual bool) {
	t.Helper()

	Equal(t, actual, true)
}

func False(t *testing.T, actual bool) {
	t.Helper()

	Equal(t, actual, false)
}
