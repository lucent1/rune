package store

func (r *Rune) Set(key string, val []byte) {
	cp := make([]byte, len(val))
	copy(cp, val)

	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = cp
}

func (r *Rune) Get(key string) []byte {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.data[key]
	if !ok {
		return nil
	}

	cp := make([]byte, len(v))
	copy(cp, v)
	return cp
}

func (r *Rune) Delete(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.data[key]
	if ok {
		delete(r.data, key)
	}

	return ok
}
