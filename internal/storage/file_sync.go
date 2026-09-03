package storage

import (
	"os"
	"path/filepath"
)

func createFileSync(name string) (*os.File, error) {
	fp, err := os.OpenFile(
		name,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	if err := syncDir(filepath.Dir(name)); err != nil {
		_ = fp.Close()
		return nil, err
	}

	return fp, nil
}
