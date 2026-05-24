package table

type Celltype uint8

const (
	TypeI64 Celltype = 1
	TypeStr Celltype = 2
)

type Cell struct {
	Type Celltype
	I64  int64
	Str  []byte
}

func (cell *Cell) Encode(toAppend []byte) []byte {
    switch cell.Type {
    case TypeI64:
        // TODO
    case TypeStr:
        // TODO
    }
}

func (cell *Cell) Decode(data []byte) (rest []byte, err error)
