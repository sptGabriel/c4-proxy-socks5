package app

const blockSize int = 128

type GameCrypt struct {
	xorEncKey []byte
	xorDecKey []byte
	isEnabled bool
}

func (gc *GameCrypt) SetEnabled() {
	gc.isEnabled = true
}

func (gc *GameCrypt) KeyHack(xorEnc, xorDec []byte) {
	gc.xorEncKey = xorEnc
	gc.xorDecKey = xorDec
}

func (gc *GameCrypt) SetKey(key uint32) {
	xorKey := []byte{
		byte((key >> 0) & 0xff), byte((key >> 8) & 0xff), byte((key >> 16) & 0xff), byte((key >> 24) & 0xff),
		0xa1, 0x6c, 0x54, 0x87,
	}

	gc.xorEncKey = append([]byte{}, xorKey...)
	gc.xorDecKey = append([]byte{}, xorKey...)
}

func (gc *GameCrypt) updateXorKey(key []byte, length int) []byte {
	keyPart := uint32(0)
	keyPart |= uint32(key[0]) << 0
	keyPart |= uint32(key[1]) << 8
	keyPart |= uint32(key[2]) << 16
	keyPart |= uint32(key[3]) << 24
	keyPart += uint32(length)

	key[0] = byte((keyPart >> 0) & 0xff)
	key[1] = byte((keyPart >> 8) & 0xff)
	key[2] = byte((keyPart >> 16) & 0xff)
	key[3] = byte((keyPart >> 24) & 0xff)

	return key
}

func (gc *GameCrypt) Decrypt(raw []byte, offset int, size int) error {
	if !gc.isEnabled {
		gc.isEnabled = true

		return nil
	}

	if gc.xorDecKey != nil {
		var prevByte byte = 0
		for i := 0; i < size; i++ {
			tmp := raw[offset+i] & 0xff
			raw[offset+i] ^= gc.xorDecKey[i&7] ^ prevByte
			prevByte = tmp
		}
		gc.xorDecKey = gc.updateXorKey(gc.xorDecKey, size)
	}

	return nil
}

func (gc *GameCrypt) Encrypt(raw []byte, offset int, size int) {
	if !gc.isEnabled {
		return
	}

	if gc.xorEncKey != nil {
		blockCount := size / blockSize
		if size%blockSize != 0 {
			blockCount++
		}

		var prevByte byte = 0
		for i := 0; i < blockCount; i++ {
			bSize := blockSize
			if i == blockCount-1 {
				bSize = size % blockSize
			}

			for b := 0; b < bSize; b++ {
				raw[offset+i*bSize+b] ^= gc.xorEncKey[b&7] ^ prevByte
				prevByte = raw[offset+i*bSize+b] & 0xff
			}
		}
		gc.xorEncKey = gc.updateXorKey(gc.xorEncKey, size)
	}
}

func NewGameCrypt() *GameCrypt {
	return &GameCrypt{
		isEnabled: false,
	}
}
