package storage

import (
	"bytes"
	"testing"
)

func TestEntryEncode(t *testing.T) {
	original := Entry{
		key:     []byte("a"),
		val:     []byte("bb"),
		deleted: false,
	}

	encoded := original.Encode()

	expected := []byte{
		1, 0, 0, 0,
		2, 0, 0, 0,
		0,
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
		0,
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

	if decoded.deleted {
		t.Fatalf("expected deleted:%v, got %v", false, decoded.deleted)
	}
}

func TestEntryEncodeDecode(t *testing.T) {
	cases := []Entry{
		{
			key:     []byte("hello"),
			val:     []byte("world"),
			deleted: false,
		},
		{
			key:     []byte("foo"),
			val:     nil,
			deleted: true,
		},
		{
			key:     []byte("binary"),
			val:     []byte{0xFF, 0x00, 0xAB},
			deleted: false,
		},
	}

	for _, original := range cases {
		t.Run(string(original.key), func(t *testing.T) {
			encoded := original.Encode()

			var decoded Entry

			if err := decoded.Decode(bytes.NewReader(encoded)); err != nil {
				t.Fatalf("decode failed: %v", err)
			}

			if !bytes.Equal(original.key, decoded.key) {
				t.Errorf("expected key %v, got %v", original.key, decoded.key)
			}

			if !bytes.Equal(original.val, decoded.val) {
				t.Errorf("expected value %v, got %v", original.val, decoded.val)
			}

			if original.deleted != decoded.deleted {
				t.Errorf("expected deleted=%v, got %v", original.deleted, decoded.deleted)
			}
		})
	}
}
