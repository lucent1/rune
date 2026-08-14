package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEncodeEntry(t *testing.T) {

	data := encodeEntry(0x01, "tushig", []byte("1234"))

	if data[0] != 0x01 {
		t.Errorf("expected 0x01, got 0x%02x", data[0])
	}
}

func TestEncodeEntryDELETE(t *testing.T) {

	data := encodeEntry(0x02, "tushig", []byte(""))

	if data[0] != 0x02 {
		t.Errorf("expected 0x02, got 0x%02x", data[0])
	}

	keyLen := len("tushig")
	offset := 1 + 2 + keyLen

	for i := offset; i < offset+4; i++ {
		if data[i] != 0x00 {
			t.Errorf("expected value 0, got %d", data[i])
		}
	}
}

func TestDecodeEntry(t *testing.T) {
	encoded := encodeEntry(0x01, "tushig", []byte("1234"))
	op, key, val, err := decodeEntry(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if op != 0x01 || key != "tushig" || string(val) != "1234" {
		t.Errorf("round trip fail")
	}
}

func TestDecodeEntryDELETE(t *testing.T) {
	encoded := encodeEntry(0x02, "tushig", []byte(""))
	op, key, val, err := decodeEntry(encoded)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if op != 0x02 || key != "tushig" || len(val) != 0 {
		t.Error("round trip fail")
	}
}

func TestDecodeEntryCorruption(t *testing.T) {
	malformed := []byte{0xFF}
	_, _, _, err := decodeEntry(malformed)

	if err == nil {
		t.Error("expected error for malformed data, got nil")
	}
}

func TestEncodeEntryEmptyValues(t *testing.T) {
	data := encodeEntry(0x01, "", []byte("1234"))

	if data[0] != 0x01 {
		t.Errorf("expected operation 0x01, got 0x%02x", data[0])
	}

	if data[1] != 0x00 || data[2] != 0x00 {
		t.Error("expected key_len to be 0x0000 for empty key")
	}

	data = encodeEntry(0x01, "tushig", []byte(""))

	if data[0] != 0x01 {
		t.Errorf("expected operation 0x01, got 0x%02x", data[0])
	}

	offset := 1 + 2 + len("tushig")

	for i := offset; i < offset+4; i++ {
		if data[i] != 0x00 {
			t.Errorf("expected value_len byte to be 0x00, got 0x%02x", data[i])
		}
	}
}

func TestDecodeEntryInvalidOp(t *testing.T) {
	encoded := encodeEntry(0xFF, "tushig", []byte("1234"))

	_, _, _, err := decodeEntry(encoded)
	if err == nil {
		t.Error("decodeEntry should return error")
	}
}

func TestDecodeEntryTruncated(t *testing.T) {
	data := encodeEntry(0x01, "tushig", []byte("1234"))
	data = data[:7]
	_, _, _, err := decodeEntry(data)
	if err == nil {
		t.Error("truncated data, should return error")
	}
}

func TestDecodeEntryKeyLengthTooLarge(t *testing.T) {
	data := encodeEntry(0x01, "tushig", []byte("1234"))
	data[1] = byte(9)

	_, _, _, err := decodeEntry(data)
	if err == nil {
		t.Error("key length larger than actual key, should return error")
	}
}

func TestDecodeEntryValueLengthTooLarge(t *testing.T) {
	data := encodeEntry(0x01, "tushig", []byte("1234"))
	offset := 1 + 2 + len("tushig")
	data[offset] = byte(9)

	_, _, _, err := decodeEntry(data)
	if err == nil {
		t.Error("value length larger than actual key, should return error")
	}
}

func TestAOFWrite(t *testing.T) {
	tmpDir := t.TempDir()
	aofPath := filepath.Join(tmpDir, "test.aof")

	aof, err := NewAOF(aofPath)
	if err != nil {
		t.Fatalf("unable to create new aof: %v", err)
	}
	defer aof.Close()

	err = aof.WriteEntry(0x01, "tushig", []byte("1234"))
	if err != nil {
		t.Fatalf("unable to write entry: %v", err)
	}

	err = aof.Sync()
	if err != nil {
		t.Fatalf("unable to sync data: %v:", err)
	}

	//close now then read from disk
	aof.Close()

	data, err := os.ReadFile(aofPath)
	if err != nil {
		t.Fatalf("unable to read from disk: %v", err)
	}

	if len(data) == 0 || data[0] != 0x01 {
		t.Fatalf("data not matching: expected first byte: 0x01, got : 0x%02x", data[0])
	}
}

func TestAOFRecovery(t *testing.T) {
	tmpDir := t.TempDir()
	aofPath := filepath.Join(tmpDir, "test.aof")

	aof, err := NewAOF(aofPath)
	if err != nil {
		t.Fatalf("unable to create new aof: %v", err)
	}
	defer aof.Close()

	entries := []struct {
		op    byte
		key   string
		value []byte
	}{
		{0x01, "w1", []byte("1234")},
		{0x01, "w2", []byte("1234")},
		{0x01, "w3", []byte("1234")},
		{0x01, "w4", []byte("1234")},
		{0x01, "w5", []byte("1234")},
		{0x02, "d1", []byte("")}, // DELETE
	}

	for _, e := range entries {
		err = aof.WriteEntry(e.op, e.key, e.value)
		if err != nil {
			t.Fatalf("unable to write entry: %v", err)
		}
	}

	err = aof.Sync()
	if err != nil {
		t.Fatalf("unable to sync data: %v:", err)
	}

	aof.Close()

	data, err := os.ReadFile(aofPath)
	if err != nil {
		t.Fatalf("unable to read file: %v", err)
	}

	offset := 0
	for i, expected := range entries {
		if offset >= len(data) {
			t.Fatalf("ran out of data at entry: %d", i)
		}

		op, key, val, err := decodeEntry(data[offset:])
		if err != nil {
			t.Fatalf("failed to decode entry: %v", err)
		}

		if op != expected.op || key != expected.key {
			t.Errorf("entry %d missmatch: got op=0x%02x key=%s, expected op=0x%02x key=%s", i, key, val, expected.key, expected.value)
		}

		entryBytes := encodeEntry(op, key, val)
		offset += len(entryBytes)
	}

}
