// log.go wraps the append-only Store record file used by KVStore.
// Writes are synced immediately, and reads stream records back in insertion order.
package kvstore

import (
	"io"
	"os"
)

type Log struct {
	Filename string
	fp       *os.File
}

func (log *Log) Open() (err error) {
	log.fp, err = createFileSync(log.Filename)
	return err
}

func (Log *Log) Close() error {
	return Log.fp.Close()
}

func (Log *Log) Write(st *Store) error {
	if _, err := Log.fp.Write(st.Encode()); err != nil {
		return err
	}
	return Log.fp.Sync()
}

func (Log *Log) Read(st *Store) (eof bool, err error) {
	err = st.Decode(Log.fp)

	if err == io.EOF || err == ErrBadSum || err == io.ErrUnexpectedEOF {
		return true, err
	} else if err != nil {
		return false, nil
	} else {
		return false, nil
	}
}
