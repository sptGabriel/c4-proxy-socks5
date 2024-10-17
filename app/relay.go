package app

import (
	"io"
	"log"
	"net"
)

func relay(src, dst net.Conn, modifyFunc func([]byte) []byte) error {
	buffer := make([]byte, 4096) // Tamanho do buffer para ler os dados
	for {
		// Ler os dados da conexão fonte
		n, err := src.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Erro ao ler da conexão: %v", err)
			}
			return err
		}

		// Manipular os dados, se necessário
		data := buffer[:n]

		if data[0] == 159 {
			// sessionID := data[3:7]
			// publicKey := data[11:139]
			// log.Printf("Pacote Init: SessionID: %d, PublicKey: %d, BlowfishKey: %x", sessionID, publicKey)
		}

		if modifyFunc != nil {
			data = modifyFunc(data)
		}

		// Escrever os dados na conexão de destino
		_, err = dst.Write(data)
		if err != nil {
			log.Printf("Erro ao escrever na conexão de destino: %v", err)
			return err
		}
	}
}
