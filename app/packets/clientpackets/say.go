package clientpackets

import (
	"bytes"
	"fmt"

	"github.com/sptGabriel/socks5/app/packets"
)

// Say representa a estrutura do pacote
type Say struct {
	Text   string
	Type   int
	Target string
}

func (s *Say) ToBytes() ([]byte, error) {
	targetSize := len(s.Target)

	size := 5 + len(s.Text) + targetSize

	// Cria um slice de bytes com o tamanho exato
	buf := bytes.NewBuffer(make([]byte, 0, size))

	buf.WriteByte(0x38) // Escreve o byte 0x49

	// Escreve o campo Text (String)
	packets.WriteS(buf, s.Text)

	// Escreve o campo Type (int32)
	packets.WriteUInt32(buf, uint32(s.Type))

	if s.Type == 1 {
		packets.WriteS(buf, s.Target)
	}

	return buf.Bytes(), nil
}

func ReadSay(reader *packets.Reader) {
	packetType := int32(reader.ReadUInt8())
	text := reader.ReadString()
	textType := int32(reader.ReadUInt32())

	fmt.Println(packetType, text, textType)
}
