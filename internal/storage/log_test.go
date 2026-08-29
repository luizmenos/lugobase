package storage

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestLogWriteRead(t *testing.T) {
	dir := t.TempDir()
	fileName := filepath.Join(dir, "test.log")

	log := Log{
		FileName: fileName,
	}

	if err := log.Open(); err != nil {
		t.Fatalf("failed to open log: %v", err)
	}
	defer log.Close()

	original := Entry{
		key: []byte("language"),
		val: []byte("go"),
	}

	if err := log.Write(&original); err != nil {
		t.Fatalf("failed to write log: %v", err)
	}

	if _, err := log.SeekToStart(); err != nil {
		t.Fatalf("failed to seek: %v", err)
	}

	var decoded Entry

	eof, err := log.Read(&decoded)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if eof {
		t.Fatal("unexpected EOF")
	}

	if !bytes.Equal(decoded.key, original.key) {
		t.Errorf("expected key %q, got %q", original.key, decoded.key)
	}

	if !bytes.Equal(decoded.val, original.val) {
		t.Errorf("expected value %q, got %q", original.val, decoded.val)
	}

	eof, err = log.Read(&decoded)

	if err != nil {
		t.Fatalf("failed to read EOF: %v", err)
	}
	if !eof {
		t.Fatal("expected EOF")
	}
}
