package storage

import (
	"bytes"
	"testing"
)

func TestEntryEncode(t *testing.T) {
	original := Entry{
		key: []byte("a"),
		val: []byte("bb"),
	}

	encoded := original.Encode()

	expected := []byte{
		1, 0, 0, 0,
		2, 0, 0, 0,
		'a',
		'b',
		'b',
	}

	if !bytes.Equal(expected, encoded) {
		t.Fatalf("expected %v, got %v", expected, encoded)
	}
}

func TestEntryDecode(t *testing.T) {
	encoded := []byte{
		1, 0, 0, 0,
		2, 0, 0, 0,
		'a',
		'b',
		'b',
	}

	var decoded Entry

	err := decoded.Decode(bytes.NewReader(encoded))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(decoded.key, []byte("a")) {
		t.Fatalf("expected key %q, got %q", []byte("a"), decoded.key)
	}

	if !bytes.Equal(decoded.val, []byte("bb")) {
		t.Fatalf("expected value %q, got %q", []byte("bb"), decoded.val)
	}
}

func TestEntryEncodeDecode(t *testing.T) {
	original := Entry{
		key: []byte("hello"),
		val: []byte("world"),
	}

	encoded := original.Encode()

	var decoded Entry

	if err := decoded.Decode(bytes.NewReader(encoded)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(original.key, decoded.key) {
		t.Fatalf("expected key %q, got %q", original.key, decoded.key)
	}

	if !bytes.Equal(original.val, decoded.val) {
		t.Fatalf("expected value %q, got %q", original.val, decoded.val)
	}
}
