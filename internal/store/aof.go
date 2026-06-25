package store

import "os"

type AOF struct {
	file *os.File
}

func NewAOF(path string) (*AOF, error) {}

func (a *AOF) WriteEntry(op byte, key string, value []byte) error {}

func (a *AOF) Sync() error {}

func (a *AOF) Close() error {}
