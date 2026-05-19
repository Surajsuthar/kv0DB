package kvstore

import (
	"encoding/binary"
	"io"
)

type Store struct {
	key 	[]byte
	val 	[]byte
	deleted bool
}

func (st *Store) Encode() []byte {
	data := make([]byte, 4+4+1+len(st.key)+len(st.val))
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(st.key)))
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(st.val)))
	if st.deleted {
		data[9] = 1
	}
	copy(data[9:len(st.key)], st.key)
	copy(data[9+len(st.key):], st.val)
	return data
}

func (st *Store) Decode(r io.Reader) error {
	var head [9]byte
	if _, err := io.ReadFull(r, head[:]); err != nil {
		return err
	}
	keylen := int(binary.LittleEndian.Uint32(head[0:4]))
	vallen := int(binary.LittleEndian.Uint32(head[4:8]))

	deleted := head[8];
	data := make([]byte, keylen+vallen)
	if _, err := io.ReadFull(r, data); err != nil {
		return err
	}

	st.key = data[0:keylen]
	if deleted != 0 {
		st.deleted = true
	}else {
		st.val = data[keylen:]
		st.deleted = false
	}
	return nil
}
