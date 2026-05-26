// kv_ser_de.go serializes and deserializes log records.
//
// Record format:
// | crc32 | key size | val size | deleted | key data | val data |
// |  4B   |    4B    |    4B    |   1B    |    ...   |    ...   |
package kvstore

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

type Store struct {
	key     []byte
	val     []byte
	deleted bool
}

var ErrBadSum = errors.New("bad checksum")

func (st *Store) Encode() []byte {
	data := make([]byte, 4+4+4+1+len(st.key)+len(st.val))
	crc := crc32.ChecksumIEEE(data[4:])
	binary.LittleEndian.PutUint32(data[0:4], crc)

	binary.LittleEndian.PutUint32(data[4:8], uint32(len(st.key)))
	binary.LittleEndian.PutUint32(data[8:12], uint32(len(st.val)))
	if st.deleted {
		data[12] = 1
	}
	copy(data[13:len(st.key)+13], st.key)
	copy(data[13+len(st.key):], st.val)
	return data
}

func (st *Store) Decode(r io.Reader) error {
	var head [13]byte
	if _, err := io.ReadFull(r, head[:]); err != nil {
		return err
	}

	keylen := int(binary.LittleEndian.Uint32(head[4:8]))
	vallen := int(binary.LittleEndian.Uint32(head[8:12]))

	deleted := head[12]
	data := make([]byte, keylen+vallen)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil
	}

	h := crc32.NewIEEE()
	h.Write(head[:4])
	h.Write(data)

	if h.Sum32() != binary.LittleEndian.Uint32(head[0:4]) {
		return ErrBadSum
	}

	st.key = data[0:keylen]
	if deleted != 0 {
		st.deleted = true
	} else {
		st.val = data[keylen:]
		st.deleted = false
	}
	return nil
}
