package storage

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"testing"
)

func TestEntryEncode(t *testing.T) {
	original := Entry{
		key:     []byte("a"),
		val:     []byte("bb"),
		deleted: false,
	}

	encoded := original.Encode()

	data := []byte{
		1, 0, 0, 0,
		2, 0, 0, 0,
		0,
		'a',
		'b',
		'b',
	}
	expected := appendChecksum(data)

	if !bytes.Equal(expected, encoded) {
		t.Fatalf("expected %v, got %v", expected, encoded)
	}
}

func TestEntryDecode(t *testing.T) {
	data := []byte{
		1, 0, 0, 0,
		2, 0, 0, 0,
		0,
		'a',
		'b',
		'b',
	}
	encoded := appendChecksum(data)

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

func TestEntryDecodeEOF(t *testing.T) {
	reader := bytes.NewReader((&Entry{
		key: []byte("key"),
		val: []byte("value"),
	}).Encode())

	var decoded Entry
	if err := decoded.Decode(reader); err != nil {
		t.Fatalf("unexpected error decoding valid entry: %v", err)
	}

	if err := decoded.Decode(reader); !errors.Is(err, io.EOF) {
		t.Fatalf("expected %v, got %v", io.EOF, err)
	}
}

func TestEntryDecodeUnexpectedEOF(t *testing.T) {
	valid := (&Entry{
		key: []byte("key"),
		val: []byte("value"),
	}).Encode()

	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "partial header",
			data: valid[:8],
		},
		{
			name: "partial key",
			data: valid[:10],
		},
		{
			name: "partial value",
			data: valid[:14],
		},
		{
			name: "missing checksum",
			data: valid[:len(valid)-4],
		},
		{
			name: "partial checksum",
			data: valid[:len(valid)-2],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoded Entry
			err := decoded.Decode(bytes.NewReader(tt.data))

			if !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Fatalf("expected %v, got %v", io.ErrUnexpectedEOF, err)
			}
		})
	}
}

func TestEntryDecodeBadSum(t *testing.T) {
	encoded := (&Entry{
		key: []byte("key"),
		val: []byte("value"),
	}).Encode()

	encoded[9] ^= 0xFF

	var decoded Entry
	err := decoded.Decode(bytes.NewReader(encoded))

	if !errors.Is(err, ErrBadSum) {
		t.Fatalf("expected %v, got %v", ErrBadSum, err)
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

func appendChecksum(data []byte) []byte {
	encoded := make([]byte, len(data)+4)
	copy(encoded, data)
	binary.LittleEndian.PutUint32(encoded[len(data):], crc32.ChecksumIEEE(data))

	return encoded
}
