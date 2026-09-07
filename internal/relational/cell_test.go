package relational

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestCellEncoding(t *testing.T) {
	cases := []struct {
		name string
		cell Cell
		data []byte
	}{
		{"zero", Cell{Type: TypeI64}, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{"little endian", Cell{Type: TypeI64, I64: 0x0102030405060708}, []byte{8, 7, 6, 5, 4, 3, 2, 1}},
		{"negative", Cell{Type: TypeI64, I64: -2}, []byte{0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{"minimum int64", Cell{Type: TypeI64, I64: -1 << 63}, []byte{0, 0, 0, 0, 0, 0, 0, 0x80}},
		{"maximum int64", Cell{Type: TypeI64, I64: 1<<63 - 1}, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}},
		{"nil string", Cell{Type: TypeStr}, []byte{0, 0, 0, 0}},
		{"empty string", Cell{Type: TypeStr, Str: []byte{}}, []byte{0, 0, 0, 0}},
		{"string", Cell{Type: TypeStr, Str: []byte("abc")}, []byte{3, 0, 0, 0, 'a', 'b', 'c'}},
		{"binary string", Cell{Type: TypeStr, Str: []byte{0, 0xff, 0x80}}, []byte{3, 0, 0, 0, 0, 0xff, 0x80}},
		{"multibyte length", Cell{Type: TypeStr, Str: bytes.Repeat([]byte{'x'}, 256)}, append([]byte{0, 1, 0, 0}, bytes.Repeat([]byte{'x'}, 256)...)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("encode", func(t *testing.T) {
				if got := tc.cell.Encode(nil); !bytes.Equal(got, tc.data) {
					t.Fatalf("expected %x, got %x", tc.data, got)
				}
			})

			t.Run("append to prefix", func(t *testing.T) {
				prefix := []byte{0xaa, 0xbb}
				want := append([]byte{0xaa, 0xbb}, tc.data...)
				if got := tc.cell.Encode(prefix); !bytes.Equal(got, want) {
					t.Fatalf("expected %x, got %x", want, got)
				}
			})

			t.Run("decode", func(t *testing.T) {
				// Decode a known byte sequence independently of Encode.
				data := append(append([]byte(nil), tc.data...), 0xaa, 0xbb)
				decoded := Cell{Type: tc.cell.Type}
				rest, err := decoded.Decode(data)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if decoded.Type != tc.cell.Type || decoded.I64 != tc.cell.I64 || !bytes.Equal(decoded.Str, tc.cell.Str) {
					t.Fatalf("expected %+v, got %+v", tc.cell, decoded)
				}
				if !bytes.Equal(rest, []byte{0xaa, 0xbb}) {
					t.Fatalf("unexpected remaining bytes: %x", rest)
				}
			})

			t.Run("truncated", func(t *testing.T) {
				// Every cut before the end must fail, including empty input.
				for n := 0; n < len(tc.data); n++ {
					decoded := Cell{Type: tc.cell.Type}
					_, err := decoded.Decode(tc.data[:n])
					if !errors.Is(err, io.ErrUnexpectedEOF) {
						t.Fatalf("with %d bytes: expected ErrUnexpectedEOF, got %v", n, err)
					}
				}
			})
		})
	}
}

func TestCellEncodeDecodeSequence(t *testing.T) {
	cells := []Cell{
		{Type: TypeI64, I64: -42},
		{Type: TypeStr, Str: []byte("hello")},
		{Type: TypeStr},
		{Type: TypeI64, I64: 123},
	}
	var data []byte
	for _, cell := range cells {
		data = cell.Encode(data)
	}

	for i, want := range cells {
		decoded := Cell{Type: want.Type}
		rest, err := decoded.Decode(data)
		if err != nil {
			t.Fatalf("cell %d: unexpected error: %v", i, err)
		}
		if decoded.I64 != want.I64 || !bytes.Equal(decoded.Str, want.Str) {
			t.Fatalf("cell %d: expected %+v, got %+v", i, want, decoded)
		}
		data = rest
	}
	if len(data) != 0 {
		t.Fatalf("expected no remaining bytes, got %x", data)
	}
}

func TestCellDecodeOversizedString(t *testing.T) {
	cell := Cell{Type: TypeStr}
	_, err := cell.Decode([]byte{0xff, 0xff, 0xff, 0xff})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected ErrUnexpectedEOF, got %v", err)
	}
}

func TestCellInvalidType(t *testing.T) {
	for _, typ := range []CellType{0, 3, 255} {
		t.Run(fmt.Sprint(typ), func(t *testing.T) {
			cell := Cell{Type: typ}
			t.Run("decode returns error", func(t *testing.T) {
				if _, err := cell.Decode(nil); err == nil {
					t.Fatal("expected an error for invalid cell type")
				}
			})
			t.Run("encode panics", func(t *testing.T) {
				defer func() {
					if recover() == nil {
						t.Error("expected a panic for invalid cell type")
					}
				}()
				cell.Encode(nil)
			})
		})
	}
}
