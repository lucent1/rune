package store

import "sync"

type Rune struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewRune() *Rune {
	return &Rune{
		data: make(map[string][]byte),
	}
}
