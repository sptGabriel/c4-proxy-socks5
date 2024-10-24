package app

import (
	"errors"
)

var blowFishStatic = []byte{
	0x5f, 0x3b, 0x35, 0x2e,
	0x5d, 0x39, 0x34, 0x2d,
	0x33, 0x31, 0x3d, 0x3d,
	0x2d, 0x25, 0x78, 0x54,
	0x21, 0x5e, 0x5b, 0x24, 0x00,
}

type LoginCrypt struct {
	isInit bool
	crypt  Crypt
}

func (lc *LoginCrypt) DecryptFromClient(raw []byte, offset, size int) error {
	return lc.decrypt(raw, offset, size)
}

func (lc *LoginCrypt) Decrypt(raw []byte, offset, size int) error {
	if lc.isInit {
		lc.isInit = false

		return nil
	}

	return lc.decrypt(raw, offset, size)
}

func (lc *LoginCrypt) decrypt(raw []byte, offset, size int) error {
	if offset+size > len(raw) {
		return errors.New("raw array too short for size starting from offset")
	}

	lc.crypt.Decrypt(raw)

	return lc.crypt.VerifyChecksum(raw, offset, size)
}

func (lc LoginCrypt) Encrypt(data []byte, offset, size int) []byte {
	padSize := 8 - (size % 8)
	size += padSize

	remaing := size - len(data)
	if remaing > 0 {
		data = append(data, make([]byte, remaing)...)
	}

	lc.crypt.AppendChecksum(data, offset, size-12)
	lc.crypt.Encrypt(data)

	return data
}

func convertToLittleEndian(data []byte) {
	for i := 0; i < len(data); i += 4 {
		data[i], data[i+3] = data[i+3], data[i]
		data[i+1], data[i+2] = data[i+2], data[i+1]
	}
}

func convertToBigEndian(data []byte) {
	for i := 0; i < len(data); i += 4 {
		data[i], data[i+3] = data[i+3], data[i]
		data[i+1], data[i+2] = data[i+2], data[i+1]
	}
}

func NewLoginCrypt() (*LoginCrypt, error) {
	crypt, err := NewCrypt(blowFishStatic)
	if err != nil {
		return &LoginCrypt{}, err
	}

	return &LoginCrypt{
		crypt:  crypt,
		isInit: true,
	}, nil
}
