// datatype_test.go checks binary encoding for primitive table cell values.
package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatatype(t *testing.T) {
	cell := Cell{
		Type: TypeI64,
		I64:  -2,
	}

	data := []byte{0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	assert.Equal(t, data, cell.Encode(nil))
}
