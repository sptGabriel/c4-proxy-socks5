package app

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

func relay(src, dst net.Conn, modifyFunc func([]byte) []byte) error {
	buffer := make([]byte, 64*1024) // Tamanho do buffer para ler os dados
	for {
		// Ler os dados da conexão fonte
		n, err := src.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Erro ao ler da conexão: %v", err)
			}
			return err
		}

		data := buffer[:n]

		if modifyFunc != nil {
			data = modifyFunc(data)
		}

		_, err = dst.Write(data)
		if err != nil {
			log.Printf("Erro ao escrever na conexão de destino: %v", err)
			return err
		}
	}
}

func relay2(src, dst net.Conn, analyzeFunc func([]byte) []byte) error {
	buffer := make([]byte, 0) // Buffer inicial para armazenar pacotes incompletos

	for {
		readBuf := make([]byte, 4096) // Buffer temporário para ler dados da conexão
		n, err := src.Read(readBuf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Erro ao ler da conexão: %v", err)
			}
			return err
		}

		// Adicionar os novos dados lidos ao buffer de dados acumulado
		buffer = append(buffer, readBuf[:n]...)

		// Processar pacotes completos no buffer
		for len(buffer) > 2 {
			// O primeiro campo de 2 bytes é o tamanho do pacote
			packetSize := int(binary.LittleEndian.Uint16(buffer[:2]))

			// Verificar se já temos o pacote completo no buffer
			if len(buffer) < packetSize {
				break // Esperar mais dados se o pacote estiver incompleto
			}

			// Extrair o pacote completo do buffer
			packetData := buffer[:packetSize]

			// Analisar o pacote, sem modificá-lo
			if analyzeFunc != nil {
				analyzeFunc(packetData) // Processar (ex. descriptografar, logar) o conteúdo sem alterar
			}

			// Enviar o pacote exatamente como foi recebido para a conexão de destino
			_, err = dst.Write(packetData)
			if err != nil {
				log.Printf("Erro ao escrever na conexão de destino: %v", err)
				return err
			}

			// Remover o pacote processado do buffer
			buffer = buffer[packetSize:]
		}
	}
}
