package app

import (
	"crypto/cipher"
	"errors"

	"golang.org/x/crypto/blowfish"
)

type Crypt struct {
	cipher cipher.Block
}

func (c Crypt) VerifyChecksum(raw []byte, offset int, size int) error {
	// Se o tamanho não for múltiplo de 4 ou for menor ou igual a 4, retorna falso
	if size%4 != 0 || size <= 4 {
		return errors.New("the data is not a multiple of 4")
	}

	// Inicializa a variável chksum e outras
	var chksum uint32 = 0
	count := size - 4
	var check uint32
	var i int

	// Loop para calcular o checksum
	for i = offset; i < count; i += 4 {
		check = uint32(raw[i]) |
			uint32(raw[i+1])<<8 |
			uint32(raw[i+2])<<16 |
			uint32(raw[i+3])<<24
		chksum ^= check
	}

	// Calcula o último bloco de 4 bytes
	check = uint32(raw[i]) |
		uint32(raw[i+1])<<8 |
		uint32(raw[i+2])<<16 |
		uint32(raw[i+3])<<24

	// Verifica se o checksum calculado é igual ao valor esperado
	if check != chksum {
		return errors.New("could not validate checksum")
	}

	return nil
}

func (c Crypt) Decrypt(data []byte) error {
	// Check if the decrypted data is a multiple of our block size
	if len(data)%8 != 0 {
		return errors.New("the decrypted data is not a multiple of the block size")
	}

	blockSize := 8
	count := len(data) / 8

	convertToLittleEndian(data)

	for i := 0; i < count; i++ {
		block := data[i*blockSize:]
		c.cipher.Decrypt(block, block)
	}

	convertToBigEndian(data)

	return nil
}

func (c Crypt) Encrypt(data []byte) error {
	// Check if the decrypted data is a multiple of our block size
	if len(data)%8 != 0 {
		return errors.New("the decrypted data is not a multiple of the block size")
	}

	blockSize := 8
	count := len(data) / 8

	convertToLittleEndian(data)

	for i := 0; i < count; i++ {
		block := data[i*blockSize:]
		c.cipher.Encrypt(block, block)
	}

	convertToBigEndian(data)

	return nil
}

func (c Crypt) AppendChecksum(data []byte, offset, size int) {
	var checksum uint32
	for i := offset; i < size-4; i += 4 {
		val := uint32(data[i]) |
			(uint32(data[i+1]) << 8) |
			(uint32(data[i+2]) << 16) |
			(uint32(data[i+3]) << 24)
		checksum ^= val
	}
	checksum = uint32(checksum)

	// Adiciona o checksum no final
	data[size-4] = byte(checksum & 0xff)
	data[size-3] = byte((checksum >> 8) & 0xff)
	data[size-2] = byte((checksum >> 16) & 0xff)
	data[size-1] = byte((checksum >> 24) & 0xff)
}

func NewCrypt(key []byte) (Crypt, error) {
	cipher, err := blowfish.NewCipher(blowFishStatic)
	if err != nil {
		return Crypt{}, err
	}

	return Crypt{
		cipher: cipher,
	}, nil
}
