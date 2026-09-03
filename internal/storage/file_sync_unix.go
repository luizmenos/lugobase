//go:build unix

package storage

import "os"

func syncDir(dir string) error {
	fp, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer fp.Close()

	return fp.Sync()
}
