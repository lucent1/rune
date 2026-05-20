package store

func (r *Rune) Set(key string, val []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var v []byte
	copy(v, val)
	r.data[key] = v
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

func (r *Rune) Delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, key)
}
