package app

import (
	"sync"
)

const blockSize int = 128

type GameCrypt struct {
	inKey     []byte
	outKey    []byte
	xorEncKey []byte
	xorDecKey []byte
	isEnabled bool
	mu        sync.Mutex
}

func (gc *GameCrypt) KeyHack2(inKey, outKey []byte) {
	gc.inKey = make([]byte, len(inKey))
	copy(gc.inKey, inKey)

	gc.outKey = make([]byte, len(outKey))
	copy(gc.outKey, outKey)
}

func (ge *GameCrypt) InitialKey(key []byte) {
	ge.mu.Lock()
	ge.inKey = key
	ge.outKey = key
	ge.mu.Unlock()
}

func (ge *GameCrypt) Encrypt2(data []byte) error {
	if ge.outKey == nil {
		return nil
	}

	var value byte
	for index := 0; index < len(data); index++ {
		value = data[index] ^ ge.outKey[index&7] ^ value
		data[index] = value
	}

	ge.mu.Lock()
	ge.outKey = ge.updateXorKey(ge.outKey, len(data))
	ge.mu.Unlock()

	return nil
}

func (ge *GameCrypt) Decrypt2(data []byte) error {
	if ge.inKey == nil {
		return nil
	}

	var value byte
	for index := 0; index < len(data); index++ {
		dataValue := data[index]
		data[index] = dataValue ^ ge.inKey[index&7] ^ value
		value = dataValue
	}

	// Increment inKey by the length of data
	ge.mu.Lock()
	ge.inKey = ge.updateXorKey(ge.inKey, len(data))
	ge.mu.Unlock()

	return nil
}

func (gc *GameCrypt) SetEnabled() {
	gc.isEnabled = true
}

func (gc *GameCrypt) KeyHack(xorEnc, xorDec []byte) {
	gc.xorEncKey = make([]byte, len(xorEnc))
	copy(gc.xorEncKey, xorEnc)

	gc.xorDecKey = make([]byte, len(xorDec))
	copy(gc.xorDecKey, xorDec)
}

func (gc *GameCrypt) SetKey(key uint32) {
	xorKey := []byte{
		byte((key >> 0) & 0xff), byte((key >> 8) & 0xff), byte((key >> 16) & 0xff), byte((key >> 24) & 0xff),
		0xa1, 0x6c, 0x54, 0x87,
	}

	gc.mu.Lock()
	gc.xorEncKey = append([]byte{}, xorKey...)
	gc.xorDecKey = append([]byte{}, xorKey...)
	gc.inKey = append([]byte{}, xorKey...)
	gc.outKey = append([]byte{}, xorKey...)
	gc.mu.Unlock()
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

func (gc *GameCrypt) DecryptFromClient(raw []byte, offset, size int) error {
	if gc.xorDecKey == nil {
		return nil
	}

	return gc.decrypt(raw, offset, size)
}
func (gc *GameCrypt) Decrypt(raw []byte, offset int, size int) error {
	if gc.xorDecKey == nil {
		return nil
	}

	return gc.decrypt(raw, offset, size)
}

func (gc *GameCrypt) decrypt(raw []byte, offset int, size int) error {
	if gc.xorDecKey == nil {
		return nil
	}

	var prevByte byte = 0
	for i := 0; i < size; i++ {
		tmp := raw[offset+i] & 0xff
		raw[offset+i] ^= gc.xorDecKey[i&7] ^ prevByte
		prevByte = tmp
	}

	gc.mu.Lock()
	gc.xorDecKey = gc.updateXorKey(gc.xorDecKey, size)
	gc.mu.Unlock()

	return nil
}

func (gc *GameCrypt) Encrypt(raw []byte, offset int, size int) {
	if gc.xorEncKey == nil {
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

		gc.mu.Lock()
		gc.xorEncKey = gc.updateXorKey(gc.xorEncKey, size)
		gc.mu.Unlock()
	}
}

func NewGameCrypt() *GameCrypt {
	return &GameCrypt{
		isEnabled: false,
		mu:        sync.Mutex{},
	}
}
