package kvstore

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVSeriDer(t *testing.T) {
	st := Store{key: []byte("key"), val: []byte("val1")}

	data := []byte{3, 0, 0, 0, 4, 0, 0, 0, 'k', 'e', 'y', 'v', 'a', 'l', '1'}

	assert.Equal(t, data, st.Encode())

	decode := Store{}
	err := decode.Decode(bytes.NewBuffer(data))
	assert.Nil(t, err)
}
