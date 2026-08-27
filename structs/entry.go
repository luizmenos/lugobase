package structs

import "encoding/binary"

type Entry struct {
	key []byte
	val []byte
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
