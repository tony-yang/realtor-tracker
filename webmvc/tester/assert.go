package tester

import (
	"strings"
	"testing"
)

func AssertStringEqual(t *testing.T, got, want string) {
	t.Helper()
	if strings.TrimSpace(got) != strings.TrimSpace(want) {
		t.Errorf("got string '%s', want '%s'", strings.TrimSpace(got), strings.TrimSpace(want))
	}
}

func AssertIntEqual(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got int '%d', want '%d'", got, want)
	}
}

func AssertObjectEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("got object '%v', want '%v'", got, want)
	}
}

func AssertTrue(t *testing.T, got bool) {
	t.Helper()
	if got != true {
		t.Errorf("got boolean '%t', want 'true'", got)
	}
}

func AssertFalse(t *testing.T, got bool) {
	t.Helper()
	if got != false {
		t.Errorf("got boolean '%t', want 'false'", got)
	}
}
