package storage

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

type Entry struct {
	key     []byte
	val     []byte
	deleted bool
}

var ErrBadSum = errors.New("bad checksum")

func (ent *Entry) Encode() []byte {
	dataLen := 4 + 4 + 1 + len(ent.key) + len(ent.val)

	buf := make([]byte, dataLen+4)

	binary.LittleEndian.PutUint32(buf[0:4], uint32(len(ent.key)))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(ent.val)))

	if ent.deleted {
		buf[8] = 1
	}

	copy(buf[9:], ent.key)
	copy(buf[9+len(ent.key):], ent.val)

	checksum := crc32.ChecksumIEEE(buf[:dataLen])
	binary.LittleEndian.PutUint32(buf[dataLen:], checksum)

	return buf
}

func (ent *Entry) Decode(r io.Reader) error {
	header := make([]byte, 9)

	if _, err := io.ReadFull(r, header); err != nil {
		return err
	}

	keyLen := binary.LittleEndian.Uint32(header[0:4])
	valLen := binary.LittleEndian.Uint32(header[4:8])

	ent.deleted = header[8] == 1

	ent.key = make([]byte, keyLen)
	ent.val = make([]byte, valLen)

	if err := readEntryPart(r, ent.key); err != nil {
		return err
	}

	if err := readEntryPart(r, ent.val); err != nil {
		return err
	}

	checksumBytes := make([]byte, 4)
	if err := readEntryPart(r, checksumBytes); err != nil {
		return err
	}

	expected := binary.LittleEndian.Uint32(checksumBytes)

	hash := crc32.NewIEEE()
	_, _ = hash.Write(header)
	_, _ = hash.Write(ent.key)
	_, _ = hash.Write(ent.val)

	if hash.Sum32() != expected {
		return ErrBadSum
	}

	return nil
}

func readEntryPart(r io.Reader, buf []byte) error {
	_, err := io.ReadFull(r, buf)

	if errors.Is(err, io.EOF) {
		return io.ErrUnexpectedEOF
	}

	return err
}
