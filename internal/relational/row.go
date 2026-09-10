package relational

type Row []Cell

func (schema *Schema) NewRow() Row {
	return make(Row, len(schema.Cols))
}

func (row Row) EncodeKey(schema *Schema) (key []byte) {
	if len(row) != len(schema.Cols) {
		panic("row length does not match schema")
	}

	for i, cell := range row {
		if cell.Type != schema.Cols[i].Type {
			panic("cell type does not match schema")
		}
	}

	key = append([]byte(schema.Table), 0x00)

	for _, index := range schema.PKey {
		cell := row[index]

		key = cell.Encode(key)
	}

	return key
}

func (row Row) EncodeVal(schema *Schema) (val []byte)
func (row Row) DecodeKey(schema *Schema, key []byte) (err error)
func (row Row) DecodeVal(schema *Schema, val []byte) (err error)
