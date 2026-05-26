package store

import (
	"fmt"
	"sync"
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

func TestConcurrentSetGet(t *testing.T) {
	rune := NewRune()

	const workers = 100

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("key-%d", i)
			value := []byte(fmt.Sprintf("value-%d", i))

			rune.Set(key, value)

			got := rune.Get(key)

			if string(got) != string(value) {
				t.Errorf("wrong value for %s", key)
			}
		}(i)
	}

	wg.Wait()
}
