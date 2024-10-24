package app

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/sptGabriel/socks5/app/packets"
	"github.com/sptGabriel/socks5/app/packets/clientpackets"
	"github.com/sptGabriel/socks5/app/packets/serverpackets"
)

// var clientKey = []byte{
// 	0x94, 0x35, 0x00, 0x00,
// 	0xa1, 0x6c, 0x54, 0x87,
// }

var clientKey = []byte{
	0x12, 0x34, 0x00, 0x00,
	0x56, 0x78, 0x21, 0x43,
}

type GameClient struct {
	userInfo           *serverpackets.UserInfo
	serverConn         net.Conn
	encryption         *GameCrypt
	clientEncryption   *GameCrypt // Criptografia para o cliente
	firstPacketHandled bool
}

func (c GameClient) Encrypt(data []byte) ([]byte, error) {
	c.encryption.Encrypt(data[2:], 0, len(data))

	return data, nil
}

func (c *GameClient) HandleDataFromClient(data []byte) {
	fmt.Println("before manipulation client data", data)

	if c.firstPacketHandled {
		c.firstPacketHandled = false
	}

	c.clientEncryption.Decrypt2(data)
	fmt.Println("after decrypt client data", data)

	c.encryption.Encrypt2(data)
	fmt.Println("after encrypt client data", data)
}

func (c *GameClient) HandleData(data []byte) {
	fmt.Println("before manipulation server data", data)
	c.encryption.Decrypt2(data)
	fmt.Println("before decrypt server data", data)
	c.processOpcode(data)

	if c.firstPacketHandled {
		return
	}

	c.clientEncryption.Encrypt2(data)
	fmt.Println("after encrypt server data", data)
}

func (c *GameClient) processOpcode(data []byte) {
	opcode := data[0] & 0xff

	switch opcode {
	case 0x00:
		fmt.Println("RECEIVED PACKET INIT LOGIN")

		key := binary.LittleEndian.Uint32(data[2:6])
		c.encryption.SetKey(key)
		copy(data[2:], clientKey)

		c.clientEncryption.InitialKey(clientKey)

		c.firstPacketHandled = true
	case 0x01:
		fmt.Println("RECEIVED PACKET MoveToLocation")
	case 0x02:
		fmt.Println("RECEIVED PACKET NpcSay")
	case 0x03:
		fmt.Println("RECEIVED PACKET CharInfo")
		fmt.Println("CharInfo:", data)
		reader := packets.NewReader(data[1:])
		charInfo, err := serverpackets.ReadCharInfo(reader)
		fmt.Println(err)

		tatto, ok := tattoosMap[int(charInfo.UnderItemID)]
		if !ok {
		}

		fmt.Printf("O user: %s, está com a tatto: %s", charInfo.Name, tatto)

		// Criar e enviar o pacote "Say" com a mensagem
		sayPacket := clientpackets.Say{
			Text:   "ok",
			Type:   0,  // Definido como 2, conforme o formato
			Target: "", // Nenhum destinatário específico, pois Type == 2
		}

		// Convertendo o pacote em bytes
		dataToSend, err := sayPacket.ToBytes()
		if err != nil {
			fmt.Println("Erro ao converter o pacote Say:", err)
			break
		}

		fmt.Println(dataToSend)

		reader2 := packets.NewReader(dataToSend)
		clientpackets.ReadSay(reader2)

		c.encryption.Encrypt(dataToSend, 0, len(dataToSend))

		dataSize := len(dataToSend)

		sendable := make([]byte, dataSize+2)
		sendable[0] = byte((dataSize + 2) & 0xff)
		sendable[1] = byte((dataSize + 2) >> 8)

		copy(sendable[2:], dataToSend)

		fmt.Println(sendable)
		// Enviar o pacote para o servidor
		_, err = c.serverConn.Write(sendable) // Supondo que você tenha o serverConn disponível
		if err != nil {
			fmt.Println("Erro ao enviar o pacote Say para o servidor:", err)
		}

	case 0x04:
		fmt.Println("RECEIVED PACKET UserInfo")
		reader := packets.NewReader(data[1:])
		userInfo, err := serverpackets.ReadUserInfo(reader)
		fmt.Println(err)

		c.userInfo = userInfo
		tatto, ok := tattoosMap[int(c.userInfo.UnderItemID)]
		if !ok {
			break
		}

		fmt.Printf("O user: %s, está com a tatto: %s\n", c.userInfo.Name, tatto)

		// Criar e enviar o pacote "Say" com a mensagem
		sayPacket := clientpackets.Say{
			Text:   "ok",
			Type:   0,  // Definido como 2, conforme o formato
			Target: "", // Nenhum destinatário específico, pois Type == 2
		}

		// Convertendo o pacote em bytes
		dataToSend, err := sayPacket.ToBytes()
		if err != nil {
			fmt.Println("Erro ao converter o pacote Say:", err)
			break
		}

		c.encryption.Encrypt(dataToSend, 0, len(dataToSend))

		dataSize := len(dataToSend)

		sendable := make([]byte, dataSize+2)
		sendable[0] = byte((dataSize + 2) & 0xff)
		sendable[1] = byte((dataSize + 2) >> 8)

		copy(sendable[2:], dataToSend[:dataSize])

		fmt.Println(sendable)
		// Enviar o pacote para o servidor
		_, err = c.serverConn.Write(sendable) // Supondo que você tenha o serverConn disponível
		if err != nil {
			fmt.Println("Erro ao enviar o pacote Say para o servidor:", err)
		}
	case 0x05:
		fmt.Println("RECEIVED PACKET Attack")
	case 0x06:
		fmt.Println("RECEIVED PACKET Die")
	case 0x07:
		fmt.Println("RECEIVED PACKET Revive")
	case 0x0B:
		fmt.Println("RECEIVED PACKET SpawnItem")
	case 0x0C:
		fmt.Println("RECEIVED PACKET DropItem")
	case 0x0D:
		fmt.Println("RECEIVED PACKET GetItem")
	case 0x0E:
		fmt.Println("RECEIVED PACKET StatusUpdate")
	case 0x0F:
		fmt.Println("RECEIVED PACKET NpcHtmlMessage")
	case 0x12:
		fmt.Println("RECEIVED PACKET DeleteObject")
	case 0x13:
		fmt.Println("RECEIVED PACKET CharSelectionInfo")
	case 0x15:
		fmt.Println("RECEIVED PACKET CharSelected")
	case 0x16:
		fmt.Println("RECEIVED PACKET NpcInfo")
	case 0x7F:
		fmt.Println("RECEIVED PACKET TutorialShowHtml")
	case 0xA1:
		fmt.Println("RECEIVED PACKET TutorialShowQuestionMark")
	case 0xA2:
		fmt.Println("RECEIVED PACKET TutorialEnableClientEvent")
	case 0xA3:
		fmt.Println("RECEIVED PACKET TutorialCloseHtml")
	case 0xA6:
		fmt.Println("RECEIVED PACKET MyTargetSelected")
	case 0xA7:
		fmt.Println("RECEIVED PACKET PartyMemberPosition")
	case 0xB6:
		fmt.Println("RECEIVED PACKET PetDelete")
	case 0xBA:
		fmt.Println("RECEIVED PACKET VehicleStarted")
	case 0xF8:
		fmt.Println("RECEIVED PACKET SSQInfo")
	case 0xE4:
		fmt.Println("RECEIVED PACKET HennaInfo")
	case 0x1B:
		fmt.Println("RECEIVED PACKET ItemList")
	case 0x1A:
		subOpcode := data[1] & 0xff
		switch subOpcode {
		case 0x87:
			fmt.Println("RECEIVED PACKET UserInfo")
			reader := packets.NewReader(data[2:])
			userInfo, err := serverpackets.ReadUserInfo(reader)
			fmt.Print(err, userInfo)
			fmt.Println("Char Nome: ", userInfo.Name)
			fmt.Println("Char Title: ", userInfo.Title)
		default:
			fmt.Printf("Unknown extended game packet received. [0xfe 0x%x] len=%d\n", subOpcode, len(data))
		}
	case 0x4A:
		fmt.Println("RECEIVED PACKET CreatureSay")
	case 0x39:
		fmt.Println("RECEIVED PACKET AskJoinParty")
	case 0x3A:
		fmt.Println("RECEIVED PACKET JoinParty")
	case 0x64:
		fmt.Println("RECEIVED PACKET SystemMessage")
	case 0x76:
		fmt.Println("RECEIVED PACKET MagicSkillLaunched")
	case 0x48:
		fmt.Println("RECEIVED PACKET MagicSkillUse")
	case 0x60:
		fmt.Println("RECEIVED PACKET MoveToPawn")
	case 0xCE:
		fmt.Println("RECEIVED PACKET RelationChanged")
	case 0x2D:
		fmt.Println("RECEIVED PACKET SocialAction")
	case 0x29:
		fmt.Println("RECEIVED PACKET TargetSelected")
	case 0x2A:
		fmt.Println("RECEIVED PACKET TargetUnselected")
	case 0x2B:
		fmt.Println("RECEIVED PACKET AutoAttackStart")
	case 0x2C:
		fmt.Println("RECEIVED PACKET AutoAttackStop")
	case 0xf3:
		fmt.Println("RECEIVED PACKET EtcStatusUpdate")
	case 0x61:
		fmt.Println("RECEIVED PACKET ValidateLocation")
	case 0x80:
		fmt.Println("RECEIVED PACKET QuestList")
	case 0xE7:
		fmt.Println("RECEIVED PACKET SendMacroList")
	case 0xfe: // Verifica pacotes com opcode extendido
		subOpcode := data[1] & 0xff
		switch subOpcode {
		case 0x1B:
			fmt.Println("RECEIVED PACKET ExSendManorList")
		default:
			fmt.Printf("Unknown extended game packet received. [0xfe 0x%x] len=%d\n", subOpcode, len(data))
		}
	default:
		if len(data) > 2 {
			fmt.Printf("Unknown game packet received. [0x%x 0x%x] len=%d\n", opcode, data[1], len(data))
		}
	}
}

func NewGameClient(serverConn net.Conn) *GameClient {
	return &GameClient{
		clientEncryption: NewGameCrypt(),
		encryption:       NewGameCrypt(),
		serverConn:       serverConn,
	}
}
