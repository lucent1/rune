package store

import (
	"testing"
)

func TestStore(t *testing.T) {
	rune := NewRune()

	key := "aa"
	val := []byte("value")
	rune.Set(key, val)
	got := rune.Get("aa")
	if string(got) != "value" {
		t.Errorf("invalid value!")
	}
	rune.Delete("aa")
	if len(rune.data) > 0 {
		t.Errorf("Invalid length")
	}
}
