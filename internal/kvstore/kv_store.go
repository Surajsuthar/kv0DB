// kv_store.go owns the in-memory key/value map and replays the append-only log
// on startup so persisted Set and Del operations are reflected in memory.
package kvstore

import (
	"bytes"
)

type KVStore struct {
	log Log
	mp  map[string][]byte
}

type UpdateMode int

const (
	ModeUpsert UpdateMode = 0 // insert or update
	ModeInsert UpdateMode = 1 // insert new
	ModeUpdate UpdateMode = 2 // update existing
)

func (Kv *KVStore) Start() error {
	if err := Kv.log.Open(); err != nil {
		return err
	}
	Kv.mp = map[string][]byte{}
	for {
		st := &Store{}
		eof, err := Kv.log.Read(st)
		if err != nil {
			return err
		} else if eof {
			break
		}

		if st.deleted {
			delete(Kv.mp, string(st.key))
		} else {
			Kv.mp[string(st.key)] = st.val
		}
	}
	return nil
}

func (Kv *KVStore) Close() error {
	return nil
}

func (Kv *KVStore) Get(key []byte) ([]byte, bool, error) {
	val, ok := Kv.mp[string(key)]
	return val, ok, nil
}

func (Kv *KVStore) Set(key, value []byte) (bool, error) {
	return Kv.SetEx(key, value, ModeUpsert)
}

func (Kv *KVStore) Del(key []byte) (bool, error) {
	_, ok := Kv.mp[string(key)]
	if ok {
		if err := Kv.log.Write(&Store{key: key, deleted: true}); err != nil {
			return false, err
		}
		delete(Kv.mp, string(key))
	}
	return ok, nil
}

func (kv *KVStore) SetEx(key []byte, val []byte, mode UpdateMode) (bool, error) {
	preVal, ok := kv.mp[string(key)]
	var update bool
	switch mode {
	case ModeUpsert:
		update = !ok || bytes.Equal(preVal, val)
	case ModeInsert:
		update = !ok
	case ModeUpdate:
		update = ok || !bytes.Equal(preVal, val)
	default:
		panic("unreachable")
	}

	if update {
		if err := kv.log.Write(&Store{key: key, val: val}); err != nil {
			return false, err
		}
		kv.mp[string(key)] = val
	}

	return update, nil
}
