package store

import (
	"encoding/binary"
	"fmt"
	"os"
)

type AOF struct {
	file *os.File
}

// func NewAOF(path string) (*AOF, error) {}

// func (a *AOF) WriteEntry(op byte, key string, value []byte) error {}

// func (a *AOF) Sync() error {}

// func (a *AOF) Close() error {}

func encodeEntry(op byte, key string, value []byte) []byte {
	buf := make([]byte, 1+2+len(key)+4+len(value))

	buf[0] = op

	binary.BigEndian.PutUint16(buf[1:3], uint16(len(key)))

	copy(buf[3:3+len(key)], key)

	binary.BigEndian.PutUint32(buf[3+len(key):7+len(key)], uint32(len(value)))

	copy(buf[7+len(key):], value)

	return buf
}

func decodeEntry(data []byte) (op byte, key string, value []byte, err error) {
	if len(data) < 1+2+4 {
		return 0, "", nil, fmt.Errorf("data too short: %d bytes", len(data))
	}

	op = data[0]

	//validate if op is known
	if op != 0x01 && op != 0x02 {
		return 0, "", nil, fmt.Errorf("unknown operation: 0x%02x", op)
	}

	keyLen := binary.BigEndian.Uint16(data[1:3])

	if len(data) < 1+2+int(keyLen)+4 {
		return 0, "", nil, fmt.Errorf("data too short for key: need %d, have %d", 1+2+int(keyLen)+4, len(data))
	}

	key = string(data[3 : 3+keyLen])

	valueLen := binary.BigEndian.Uint32(data[3+keyLen : 7+keyLen])

	if len(data) < 1+2+int(keyLen)+4+int(valueLen) {
		return 0, "", nil, fmt.Errorf("data too short for value: need %d, have %d", 1+2+int(keyLen)+4+int(valueLen), len(data))
	}

	value = data[7+keyLen : 7+keyLen+uint16(valueLen)]

	return op, key, value, nil
}
