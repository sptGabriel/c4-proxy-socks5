package app

import (
	"net"
	"sync"
)

var (
	loginServerPORT = "2106"
	gameServerPORT  = "7777"
)

type crypt interface {
	Decrypt(raw []byte, offset, size int) error
	Encrypt(data []byte, offset, size int) []byte
}

type BaseClient struct {
	cryptService crypt
}

// func (bc *BaseClient) InitializeCrypt(port string) error {
// 	switch port {
// 	case "2106":
// 		loginCrypt, err := NewLoginCrypt()
// 		if err != nil {
// 			return err
// 		}
// 		bc.loginCrypt = &loginCrypt
// 		bc.cryptEnabled = true

// 	case "7777":
// 		gameCrypt := NewGameCrypt()
// 		bc.gameCrypt = &gameCrypt
// 		bc.cryptEnabled = true
// 	}
// 	return nil
// }

// Função para lidar com o tráfego de pacotes
func (bc *BaseClient) HandleConnection(clientConn, serverConn net.Conn, port string) error {
	var wg sync.WaitGroup
	wg.Add(2)

	// Cliente -> Servidor
	go func() {
		defer wg.Done()
		relay(clientConn, serverConn, func(data []byte) []byte {
			return data
		})
	}()

	// Servidor -> Cliente
	go func() {
		defer wg.Done()

	}()

	wg.Wait()
	return nil
}

func NewBaseClient(port string) (BaseClient, error) {
	if port == loginServerPORT {
		loginCrypt, err := NewLoginCrypt()
		if err != nil {
			return BaseClient{}, nil
		}

		return BaseClient{
			cryptService: loginCrypt,
		}, nil
	}

	gameCrypt := NewGameCrypt()
	return BaseClient{
		cryptService: gameCrypt,
	}, nil
}
