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
