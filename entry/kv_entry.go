package entry

import (
	"encoding/binary"
	"io"
)

type Entry struct {
	Key     []byte
	Value   []byte
	Deleted bool
}

func (ent *Entry) Encode() []byte {
	data := make([]byte, 4+4+1+len(ent.Key)+len(ent.Value))
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(ent.Key)))
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(ent.Value)))
	if ent.Deleted {
		data[8] = 1
	} else {
		data[8] = 0
	}

	copy(data[9:], ent.Key)
	copy(data[9+len(ent.Key):], ent.Value)
	return data
}

func (ent *Entry) Decode(r io.Reader) error {
	meta := make([]byte, 4+4+1)
	_, err := r.Read(meta)
	if err != nil {
		// todo handle case when read is not 8
		return err
	}

	kLen := binary.LittleEndian.Uint32(meta[0:4])
	vLen := binary.LittleEndian.Uint32(meta[4:8])
	kv := make([]byte, kLen+vLen)
	r.Read(kv[0:kLen])
	ent.Key = kv[0:kLen]

	if meta[8] == 1 {
		ent.Deleted = true
	} else {
		ent.Deleted = false
		r.Read(kv[kLen : kLen+vLen])
		ent.Value = kv[kLen : kLen+vLen]
	}

	return nil
}
