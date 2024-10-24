package app

import "fmt"

type LoginClient struct {
	encryption *LoginCrypt
}

func (c LoginClient) Encrypt(data []byte) ([]byte, error) {
	c.encryption.Encrypt(data[2:], 0, len(data))

	return data, nil
}

func (c LoginClient) HandleDataFromClient(data []byte) {
	return

	c.encryption.Decrypt(data, 0, len(data))
}

func (c LoginClient) HandleData(data []byte) {
	return

	c.encryption.Decrypt(data, 0, len(data))

	opcode := data[0] & 0xff
	switch opcode {
	case 0x00:
		fmt.Println("RECEIVED PACKET Init")
	case 0x01:
		fmt.Println("RECEIVED PACKET LoginFail")
	case 0x02:
		fmt.Println("RECEIVED PACKET AccountKicked")
	case 0x03:
		fmt.Println("RECEIVED PACKET LoginOk")
	case 0x04:
		fmt.Println("RECEIVED PACKET ServerList")
	case 0x06:
		fmt.Println("RECEIVED PACKET PlayFail")
	case 0x07:
		fmt.Println("RECEIVED PACKET PlayOk")
	case 0x0b:
		fmt.Println("RECEIVED PACKET GGAuth")
	default:
		if len(data) > 2 {
			fmt.Printf("Unknown game packet received. [0x%x 0x%x] len=%d\n", opcode, data[1], len(data))
		}
	}
}

func NewLoginClient() (LoginClient, error) {
	loginCrypt, err := NewLoginCrypt()
	if err != nil {
		return LoginClient{}, nil
	}

	return LoginClient{
		encryption: loginCrypt,
	}, nil
}
