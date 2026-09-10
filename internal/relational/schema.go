package relational

type Schema struct {
	Table string
	Cols  []Column
	PKey  []int
}

type Column struct {
	Name string
	Type CellType
}
