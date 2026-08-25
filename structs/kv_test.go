package structs

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

// TestLifecycle handles the full CRUD cycle using subtests
func TestLifecycle(t *testing.T) {
	kv := &KV{}
	kv.Open()
	defer kv.Close()

	key := []byte("language")
	val := []byte("go")

	// 1. Test Set
	t.Run("Set", func(t *testing.T) {
		updated, err := kv.Set(key, val)
		if err != nil {
			t.Fatalf("failed to set: %v", err)
		}
		if updated {
			t.Error("expected updated=false for new key")
		}
	})

	// 2. Test Get
	t.Run("Get", func(t *testing.T) {
		got, ok, _ := kv.Get(key)
		if !ok || !bytes.Equal(got, val) {
			t.Errorf("expected %s, got %s", val, got)
		}
	})

	// 3. Test Update
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

	// 4. Test Delete
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

// TestTableDriven tests multiple cases using a table-driven approach
func TestTableDriven(t *testing.T) {
	kv := &KV{}
	kv.Open()

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

// TestConcurrency runs multiple goroutines to check for race conditions
func TestConcurrency(t *testing.T) {
	kv := &KV{}
	kv.Open()

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
