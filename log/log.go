package log

import (
	"io"
	"os"
	"odydb/entry"
)

type Log struct {
	FileName string
	fp *os.File
}

func (log *Log) Open() (err error) {
	log.fp, err = os.OpenFile(log.FileName, os.O_RDWR|os.O_CREATE, 0o644)
	return err
}

func (log *Log) Close() error {
	return log.fp.Close()
}

func (log *Log) Write(ent *entry.Entry) error {
	_, err := log.fp.Write(ent.Encode())
	return err
}

func (log *Log) Read(ent *entry.Entry) (eof bool, err error) {
	err = ent.Decode(log.fp)

	if err == io.EOF {
		return true, nil
	} else if err != nil {
		return false, err
	} else {
		return false, nil
	}
}