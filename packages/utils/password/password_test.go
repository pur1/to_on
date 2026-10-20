package password

import (
	"errors"
	"strings"
	"testing"
)

func TestHashCompare(t *testing.T) {
	h, err := Hash("correct horse")
	if err != nil {
		t.Fatal(err)
	}
	if !Compare(h, "correct horse") {
		t.Error("expected match")
	}
	if Compare(h, "wrong horse") {
		t.Error("expected mismatch")
	}
}

func TestHashTooLong(t *testing.T) {
	_, err := Hash(strings.Repeat("a", MaxLength+1))
	if !errors.Is(err, ErrTooLong) {
		t.Errorf("got %v, want ErrTooLong", err)
	}
}
