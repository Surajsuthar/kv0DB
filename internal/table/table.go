// table.go defines schemas and rows, plus how rows are split into key bytes and
// value bytes for storage in the lower-level key/value engine.
//
// Row storage shape:
// key   = table name + 0x00 + primary-key columns
// value = non-primary-key columns
package table

import (
	"errors"
	"slices"
)

type Schema struct {
	Table string
	Cols  []Column
	PKey  []int // primary key column
}

type Column struct {
	Name string
	Type Celltype
}

type Row []Cell

func (schema *Schema) NewRow() Row {
	return make(Row, len(schema.Cols))
}

func (row Row) EncodeKey(schema *Schema) (key []byte) {
	key = append([]byte(schema.Table), 0x00)
	if len(row) != len(schema.Cols) {
		panic("Assert Fail")
	}

	for idx, val := range row {
		if val.Type != schema.Cols[idx].Type {
			panic("Assert Fail")
		}
		if slices.Contains(schema.PKey, idx) {
			key = row[idx].Encode(key)
		}
	}

	return key
}

func (row Row) EncodeVal(schema *Schema) (val []byte) {
	if len(row) != len(schema.Cols) {
		panic("Assert Fail")
	}
	for idx, val := range row {
		if val.Type != schema.Cols[idx].Type {
			panic("Assert Fail")
		}
		if !slices.Contains(schema.PKey, idx) {
			val = row[idx].Encode(val)
		}
	}
	return val
}

func (row Row) DecodeKey(schema *Schema, key []byte) (err error) {
	if len(row) != len(schema.Cols) {
		panic("Assert Fail")
	}

	if len(key) < len(schema.Table)+1 {
		return errors.New("bad key")
	}

	if string(key[:len(schema.Table)+1]) != schema.Table+'\x00' {
		return errors.New("bad key")
	}

	for idx, col := range schema.Cols {
		if !slices.Contains(schema.PKey, idx) {
			continue
		}
		row[idx] = Cell{Type: col.Type}
		if key, err = row[idx].Decode(key); err != nil {
			return err
		}
	}

	if len(key) != 0 {
		return errors.New("trailing garbage")
	}
	return nil
}

func (row Row) DecodeVal(schema *Schema, val []byte) (err error) {
	if len(row) != len(schema.Cols) {
		panic("Assert Fail")
	}

	for idx, col := range schema.Cols {
		if slices.Contains(schema.PKey, idx) {
			continue
		}

		row[idx] = Cell{Type: col.Type}
		if val, err = row[idx].Decode(val); err != nil {
			return err
		}
	}

	if len(val) != 0 {
		return errors.New("trailing garbage")
	}
	return nil
}
