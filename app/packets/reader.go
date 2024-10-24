package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Reader struct {
	*bytes.Reader
}

func NewReader(buffer []byte) *Reader {
	return &Reader{bytes.NewReader(buffer)}
}

func (r *Reader) ReadBytes(number int) []byte {
	buffer := make([]byte, number)
	n, _ := r.Read(buffer)
	if n < number {
		return []byte{}
	}

	return buffer
}

func (r *Reader) ReadUInt64() uint64 {
	var result uint64

	buffer := make([]byte, 8)
	n, _ := r.Read(buffer)
	if n < 8 {
		return 0
	}

	buf := bytes.NewBuffer(buffer)

	binary.Read(buf, binary.LittleEndian, &result)

	return result
}

func (r *Reader) ReadUInt32() uint32 {
	var result uint32

	buffer := make([]byte, 4)
	n, _ := r.Read(buffer)
	if n < 4 {
		return 0
	}

	buf := bytes.NewBuffer(buffer)

	binary.Read(buf, binary.LittleEndian, &result)

	return result
}

func (r *Reader) ReadUInt16() uint16 {
	var result uint16

	buffer := make([]byte, 2)
	n, _ := r.Read(buffer)
	if n < 2 {
		return 0
	}

	buf := bytes.NewBuffer(buffer)

	binary.Read(buf, binary.LittleEndian, &result)

	return result
}

func (r *Reader) ReadUInt8() uint8 {
	var result uint8

	buffer := make([]byte, 1)
	n, _ := r.Read(buffer)
	if n < 1 {
		return 0
	}

	buf := bytes.NewBuffer(buffer)

	binary.Read(buf, binary.LittleEndian, &result)

	return result
}

func (r *Reader) Position() int64 {
	position, _ := r.Seek(0, io.SeekCurrent)

	return position
}

func (r *Reader) ReadString() string {
	position, _ := r.Seek(0, io.SeekCurrent)
	fmt.Print(position)

	var result []byte
	var first_byte, second_byte byte

	for {
		first_byte, _ = r.ReadByte()
		second_byte, _ = r.ReadByte()
		if first_byte == 0x00 && second_byte == 0x00 {
			break
		} else {
			result = append(result, first_byte, second_byte)
		}
	}

	positionend, _ := r.Seek(0, io.SeekCurrent)
	fmt.Print(position, positionend)

	return string(result)
}

func (r *Reader) ReadFloat64() float64 {
	var result float64

	buffer := make([]byte, 8)
	n, _ := r.Read(buffer)
	if n < 8 {
		return 0
	}

	buf := bytes.NewBuffer(buffer)
	binary.Read(buf, binary.LittleEndian, &result)

	return result
}
