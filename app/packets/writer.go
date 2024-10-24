package packets

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

// Funções de escrita individuais
func Write(buf *bytes.Buffer, p []byte) (int, error) {
	return buf.Write(p)
}

func WriteUInt64(buf *bytes.Buffer, value uint64) {
	binary.Write(buf, binary.LittleEndian, value)
}

func WriteUInt32(buf *bytes.Buffer, value uint32) {
	binary.Write(buf, binary.LittleEndian, value)
}

func WriteUInt16(buf *bytes.Buffer, value uint16) {
	binary.Write(buf, binary.LittleEndian, value)
}

func WriteUInt8(buf *bytes.Buffer, value uint8) {
	binary.Write(buf, binary.LittleEndian, value)
}

func WriteFloat64(buf *bytes.Buffer, value float64) {
	binary.Write(buf, binary.LittleEndian, value)
}

func WriteFloat32(buf *bytes.Buffer, value float32) {
	binary.Write(buf, binary.LittleEndian, value)
}

func WriteS(buf *bytes.Buffer, value string) {
	// Converter a string para UCS-2
	ucs2 := utf16.Encode([]rune(value))

	// Escrever os bytes da string em UCS-2 (little-endian)
	for _, char := range ucs2 {
		buf.WriteByte(byte(char & 0xff))        // Byte menos significativo
		buf.WriteByte(byte((char >> 8) & 0xff)) // Byte mais significativo
	}

	// Escrever o terminador null (2 bytes de 0x00)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)
}
