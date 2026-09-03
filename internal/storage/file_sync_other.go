//go:build !unix

package storage

func syncDir(string) error {
	return nil
}
