package entry

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

const headerSize = 4 + 4 + 4 + 1

type entryHeader struct {
	Checksum uint32
	KLen uint32
	VLen uint32
	Deleted bool
}

var ErrBadSum = errors.New("bad checksum")

type Entry struct {
	Key     []byte
	Value   []byte
	Deleted bool
}

func (ent *Entry) Encode() []byte {
	h := entryHeader{
		KLen: uint32(len(ent.Key)),
		VLen: uint32(len(ent.Value)),
		Deleted: ent.Deleted,
	}

	buf := make([]byte, 0, headerSize+len(ent.Key)+len(ent.Value))
	buf, _ = binary.Append(buf, binary.LittleEndian, &h)
	buf = append(buf, ent.Key...)
	buf = append(buf, ent.Value...)

	binary.LittleEndian.PutUint32(buf[0:4], crc32.ChecksumIEEE(buf[4:]))
	return buf
}

func (ent *Entry) Decode(r io.Reader) error {
	var raw [headerSize]byte
	if _, err := io.ReadFull(r, raw[:]); err != nil {
		return err
	}

	var h entryHeader
	if _, err := binary.Decode(raw[:], binary.LittleEndian, &h); err != nil {
		return err
	}

	body := make([]byte, int(h.KLen) + int(h.VLen))
	if _, err := io.ReadFull(r, body); err != nil {
		if err == io.EOF { // If header is read = there's must be body
			err = io.ErrUnexpectedEOF
		}

		return err
	}

	checksum := crc32.Update(crc32.ChecksumIEEE(raw[4:]), crc32.IEEETable, body)
	if checksum != h.Checksum {
		return ErrBadSum
	}

	ent.Key = body[:h.KLen:h.KLen]
	ent.Deleted = h.Deleted
	if h.Deleted {
		ent.Value = nil
	} else {
		ent.Value = body[h.KLen:]
	}

	return nil
}
