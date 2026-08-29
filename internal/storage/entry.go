package storage

import (
	"encoding/binary"
	"io"
)

type Entry struct {
	key     []byte
	val     []byte
	deleted bool
}

func (ent *Entry) Encode() []byte {
	totalLen := 4 + 4 + len(ent.key) + len(ent.val)

	buf := make([]byte, totalLen)

	binary.LittleEndian.PutUint32(buf[0:4], uint32(len(ent.key)))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(ent.val)))

	copy(buf[8:], ent.key)
	copy(buf[8+len(ent.key):], ent.val)
	return buf
}

func (ent *Entry) Decode(r io.Reader) error {
	header := make([]byte, 8)

	if _, err := io.ReadFull(r, header); err != nil {
		return err
	}

	keyLen := binary.LittleEndian.Uint32(header[0:4])
	valLen := binary.LittleEndian.Uint32(header[4:8])

	ent.key = make([]byte, keyLen)
	ent.val = make([]byte, valLen)

	if _, err := io.ReadFull(r, ent.key); err != nil {
		return err
	}

	if _, err := io.ReadFull(r, ent.val); err != nil {
		return err
	}

	return nil
}
