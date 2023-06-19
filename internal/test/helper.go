package test

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Fatalf("%+v is not equal to %+v", actual, expected)
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
