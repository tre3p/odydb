package entry

import (
	"encoding/binary"
	"io"
)

type Entry struct {
	key []byte
	value []byte
}

func (ent *Entry) Encode() []byte {
	data := make([]byte, 4 + 4 + len(ent.key) + len (ent.value))
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(ent.key)))
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(ent.value)))
	copy(data[8:], ent.key)
	copy(data[8+len(ent.key):], ent.value)
	return data
}

func (ent *Entry) Decode(r io.Reader) error {
	lengths := make([]byte, 4 + 4)
	_, err := r.Read(lengths)
	if err != nil {
		// todo handle case when read is not 8
		return err
	}

	kLen := binary.LittleEndian.Uint32(lengths[0:4])
	vLen := binary.LittleEndian.Uint32(lengths[4:8])
	kv := make([]byte, kLen + vLen)
	r.Read(kv[0:kLen])
	r.Read(kv[kLen:kLen+vLen])

	ent.key = kv[0:kLen]
	ent.value = kv[kLen:kLen+vLen]
	return nil
}