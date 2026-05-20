package rune

import "sync"

type Rune struct {
	data map[string][]byte
	mu   sync.RWMutex
}
