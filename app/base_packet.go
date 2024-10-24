package app

import (
	"net"
	"sync"
)

var (
	loginServerPORT = "2106"
	// gameServerPORT  = "7777"
)

type mmoClient interface {
	Encrypt(data []byte) ([]byte, error)
	HandleData(data []byte)
	HandleDataFromClient(data []byte)
}

type BaseClient struct {
	mmoClient  mmoClient
	clientConn net.Conn
	serverConn net.Conn
}

func (bc *BaseClient) HandleConnection() error {
	var wg sync.WaitGroup
	wg.Add(2)

	// Cliente -> Servidor
	go func() {
		defer wg.Done()
		relay2(bc.clientConn, bc.serverConn, func(data []byte) []byte {
			// fmt.Println("SEND TO SERVER -> ", data)
			bc.mmoClient.HandleDataFromClient(data[2:])

			return data
		})
	}()

	// Servidor -> Cliente
	go func() {
		defer wg.Done()

		relay2(bc.serverConn, bc.clientConn, func(data []byte) []byte {
			//fmt.Println("RECEIVED FROM SERVER -> ", data)
			bc.mmoClient.HandleData(data[2:])

			return data
		})
	}()

	wg.Wait()
	return nil
}

func NewBaseClient(server, client net.Conn, port string) (BaseClient, error) {
	if port == loginServerPORT {
		loginClient, err := NewLoginClient()
		if err != nil {
			return BaseClient{}, err
		}

		return BaseClient{
			mmoClient:  loginClient,
			clientConn: client,
			serverConn: server,
		}, nil
	}

	return BaseClient{
		mmoClient:  NewGameClient(server),
		clientConn: client,
		serverConn: server,
	}, nil
}
