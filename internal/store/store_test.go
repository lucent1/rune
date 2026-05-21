package store

import (
	"testing"
)

func TestSetAndGetandDelete(t *testing.T) {
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

func TestGetDoesNotAllowMutation(t *testing.T) {
	rune := NewRune()

	key := "aa"
	val := []byte("value")
	rune.Set(key, val)
	got := rune.Get("aa")
	got[0] = 't'
	got2 := rune.Get("aa")
	if string(got2) != "value" {
		t.Errorf("invalid value!")
	}
}

func TestConcurrentSetGet(t *testing.T) {}
