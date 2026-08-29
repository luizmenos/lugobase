package storage

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
)

func logTest(t *testing.T) *KV {
	t.Helper()

	dir := t.TempDir()

	kv := &KV{
		log: Log{
			FileName: filepath.Join(dir, "test.log"),
		},
	}

	if err := kv.Open(); err != nil {
		t.Fatalf("failed to open KV: %v", err)
	}

	t.Cleanup(func() {
		kv.Close()
	})

	return kv
}

func TestLifecycle(t *testing.T) {
	kv := logTest(t)

	key := []byte("language")
	val := []byte("go")

	t.Run("Set", func(t *testing.T) {
		updated, err := kv.Set(key, val)
		if err != nil {
			t.Fatalf("failed to set: %v", err)
		}
		if updated {
			t.Error("expected updated=false for new key")
		}
	})

	t.Run("Get", func(t *testing.T) {
		got, ok, _ := kv.Get(key)
		if !ok || !bytes.Equal(got, val) {
			t.Errorf("expected %s, got %s", val, got)
		}
	})

	t.Run("Update", func(t *testing.T) {
		newVal := []byte("golang")
		updated, _ := kv.Set(key, newVal)
		if !updated {
			t.Error("expected updated=true for existing key")
		}

		got, _, _ := kv.Get(key)
		if !bytes.Equal(got, newVal) {
			t.Errorf("expected updated value %s, got %s", newVal, got)
		}
	})
	t.Run("Delete", func(t *testing.T) {
		deleted, _ := kv.Del(key)
		if !deleted {
			t.Error("expected deleted=true")
		}

		_, ok, _ := kv.Get(key)
		if ok {
			t.Error("key should be gone")
		}
	})
}

func TestTableDriven(t *testing.T) {
	kv := logTest(t)

	cases := []struct {
		name  string
		key   []byte
		value []byte
	}{
		{"Standard", []byte("name"), []byte("john")},
		{"EmptyValue", []byte("key"), []byte("")},
		{"BinaryData", []byte("img"), []byte{0xFF, 0x00, 0xAB, 0x12}},
		{"EmptyKey", []byte(""), []byte("empty")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kv.Set(tc.key, tc.value)
			got, ok, _ := kv.Get(tc.key)
			if !ok || !bytes.Equal(got, tc.value) {
				t.Errorf("failed for %s: expected %v, got %s", tc.name, tc.value, got)
			}
			kv.Del(tc.key)
		})
	}
}

func TestConcurrency(t *testing.T) {
	kv := logTest(t)

	var wg sync.WaitGroup
	workers := 50
	ops := 100

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ops; j++ {
				key := []byte(fmt.Sprintf("key-%d-%d", id, j))
				val := []byte("data")

				kv.Set(key, val)
				kv.Get(key)
				if j%10 == 0 {
					kv.Del(key)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestPersistence(t *testing.T) {
	dir := t.TempDir()
	fileName := filepath.Join(dir, "test.log")

	kv := &KV{
		log: Log{
			FileName: fileName,
		},
	}

	if err := kv.Open(); err != nil {
		t.Fatal(err)
	}

	_, err := kv.Set([]byte("language"), []byte("go"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = kv.Set([]byte("language"), []byte("php"))

	if err != nil {
		t.Fatal(err)
	}

	kv.Close()

	kv2 := &KV{
		log: Log{
			FileName: fileName,
		},
	}

	if err := kv2.Open(); err != nil {
		t.Fatal(err)
	}
	defer kv2.Close()

	got, ok, err := kv2.Get([]byte("language"))
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("expected key to exist after reopening")
	}

	if !bytes.Equal(got, []byte("php")) {
		t.Errorf("expected php, got %s", got)
	}
}

func TestDeletePersistence(t *testing.T) {
	dir := t.TempDir()
	fileName := filepath.Join(dir, "test.log")

	kv := &KV{
		log: Log{
			FileName: fileName,
		},
	}

	if err := kv.Open(); err != nil {
		t.Fatal(err)
	}

	_, err := kv.Set([]byte("language"), []byte("go"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = kv.Del([]byte("language"))
	if err != nil {
		t.Fatal(err)
	}

	kv.Close()

	kv2 := &KV{
		log: Log{
			FileName: fileName,
		},
	}

	if err := kv2.Open(); err != nil {
		t.Fatal(err)
	}
	defer kv2.Close()

	_, ok, _ := kv2.Get([]byte("language"))

	if ok {
		t.Fatal("expected key to remain deleted after reopening")
	}
}
