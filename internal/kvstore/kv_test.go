// kv_test.go covers the public KVStore behavior and the Store record encoding
// used by the append-only log.
package kvstore

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKVStoreStartAndClose(t *testing.T) {
	kv := KVStore{}

	require.NoError(t, kv.Start())
	require.NotNil(t, kv.mp)

	assert.NoError(t, kv.Close())
}

func TestKVStoreSetAndGet(t *testing.T) {
	kv := newStartedStore(t)

	updated, err := kv.Set([]byte("key1"), []byte("val1"))
	require.NoError(t, err)
	assert.True(t, updated)

	val, ok, err := kv.Get([]byte("key1"))
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []byte("val1"), val)
}

func TestKVStoreGetMissingKey(t *testing.T) {
	kv := newStartedStore(t)

	val, ok, err := kv.Get([]byte("missing"))
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestKVStoreSetReportsWhetherValueChanged(t *testing.T) {
	kv := newStartedStore(t)

	updated, err := kv.Set([]byte("key1"), []byte("val1"))
	require.NoError(t, err)
	assert.True(t, updated)

	updated, err = kv.Set([]byte("key1"), []byte("val1"))
	require.NoError(t, err)
	assert.False(t, updated)

	updated, err = kv.Set([]byte("key1"), []byte("val2"))
	require.NoError(t, err)
	assert.True(t, updated)

	val, ok, err := kv.Get([]byte("key1"))
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []byte("val2"), val)
}

func TestKVStoreDel(t *testing.T) {
	kv := newStartedStore(t)
	_, err := kv.Set([]byte("key1"), []byte("val1"))
	require.NoError(t, err)

	deleted, err := kv.Del([]byte("key1"))
	require.NoError(t, err)
	assert.True(t, deleted)

	val, ok, err := kv.Get([]byte("key1"))
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestKVStoreDelMissingKey(t *testing.T) {
	kv := newStartedStore(t)

	deleted, err := kv.Del([]byte("missing"))
	require.NoError(t, err)
	assert.False(t, deleted)
}

func newStartedStore(t *testing.T) *KVStore {
	t.Helper()

	kv := &KVStore{}
	require.NoError(t, kv.Start())
	t.Cleanup(func() {
		require.NoError(t, kv.Close())
	})

	return kv
}

func TestKVSeriDer(t *testing.T) {
	st := Store{key: []byte("key"), val: []byte("val1")}

	data := []byte{3, 0, 0, 0, 4, 0, 0, 0, 0, 'k', 'e', 'y', 'v', 'a', 'l', '1'}

	assert.Equal(t, data, st.Encode())

	decode := Store{}
	err := decode.Decode(bytes.NewBuffer(data))
	assert.Nil(t, err)
	assert.Equal(t, data, decode)

	st = Store{key: []byte("key"), deleted: true}
	data = []byte{3, 0, 0, 0, 0, 0, 0, 0, 1, 'k', 'e', 'y'}

	assert.Equal(t, data, st.Encode())
	decode = Store{}
	err = decode.Decode(bytes.NewBuffer(data))
	assert.Nil(t, err)
	assert.Equal(t, data, decode)
}
