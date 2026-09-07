package relational

import (
	"encoding/binary"
	"errors"
	"io"
)

type CellType uint8

const (
	TypeI64 CellType = 1
	TypeStr CellType = 2
)

type Cell struct {
	Type CellType
	I64  int64
	Str  []byte
}

func (cell *Cell) Encode(toAppend []byte) []byte {
	switch cell.Type {
	case TypeI64:
		return binary.LittleEndian.AppendUint64(toAppend, uint64(cell.I64))
	case TypeStr:
		toAppend = binary.LittleEndian.AppendUint32(toAppend, uint32(len(cell.Str)))
		return append(toAppend, cell.Str...)
	}
	panic("invalid cell type")
}

func (cell *Cell) Decode(data []byte) (rest []byte, err error) {
	switch cell.Type {
	case TypeI64:
		if len(data) < 8 {
			return nil, io.ErrUnexpectedEOF
		}

		cell.I64 = int64(binary.LittleEndian.Uint64(data[:8]))
		return data[8:], nil
	case TypeStr:
		if len(data) < 4 {
			return nil, io.ErrUnexpectedEOF
		}

		size := binary.LittleEndian.Uint32(data[:4])
		data = data[4:]

		if uint64(size) > uint64(len(data)) {
			return nil, io.ErrUnexpectedEOF
		}

		cell.Str = data[:size]
		return data[size:], nil
	}

	return nil, errors.New("invalid cell type")
}
