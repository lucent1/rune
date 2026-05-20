package store

func (r *Rune) Set(key string, val []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var v []byte
	copy(v, val)
	r.data[key] = v
}

func (r *Rune) Get(key string) []byte {
	return r.data[key]
}

func (r *Rune) Delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, key)
}
