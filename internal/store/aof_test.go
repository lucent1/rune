package store

import (
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
	// Test SET with empty key (edge case)
	// Test SET with empty value
	// Both should work
}

func TestDecodeEntryInvalidOp(t *testing.T) {
	// Create data with op=0xFF
	// decodeEntry should return error
}

func TestDecodeEntryTruncated(t *testing.T) {
	// Encode a normal entry
	// Truncate it (remove last few bytes)
	// decodeEntry should detect and return error
}

func TestDecodeEntryKeyLengthTooLarge(t *testing.T) {
	// Encode normally
	// Corrupt the key length to be larger than actual data
	// decodeEntry should detect bounds violation and return error
}

func TestDecodeEntryValueLengthTooLarge(t *testing.T) {
	// Encode normally
	// Corrupt the value length to be larger than actual data
	// decodeEntry should detect and return error
}
