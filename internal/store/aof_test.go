package store

import "testing"

func TestEncodeEntry(t *testing.T) {

	data := encodeEntry(0x01, "tushig", []byte("1234"))

	if data[0] != 0x01 {
		t.Error("expected 0x01, got 0x02", data[0])
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

func TestDecodeEntryCorruption(t *testing.T) {
	malformed := []byte{0xFF}
	_, _, _, err := decodeEntry(malformed)

	if err == nil {
		t.Error("expected error for malformed data, got nil")
	}
}
