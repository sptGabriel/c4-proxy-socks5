package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// before xor enc [153, 26, 0, 0, 161, 108, 84, 135]
// after xor enc [190, 26, 0, 0, 161, 108, 84, 135]

func TestGameCryptEncrypt(t *testing.T) {
	key := 6809
	data := []byte{8, 103, 0, 98, 0, 99, 0, 120, 0, 98, 0, 97, 0, 117, 0, 0, 0, 52, 229, 20, 0, 168, 30, 1, 0, 52, 229, 20, 0, 62, 89, 138, 85, 1, 0, 0, 0}
	want := []byte{145, 236, 236, 142, 47, 32, 116, 139, 18, 106, 106, 11, 170, 179, 231, 96, 249, 215, 50, 38, 135, 67, 9, 143, 22, 56, 221, 201, 104, 58, 55, 58, 246, 237, 237, 237, 76}

	gameCrypt := NewGameCrypt()
	gameCrypt.SetKey(uint32(key))
	gameCrypt.SetEnabled()

	gameCrypt.Encrypt(data, 0, 37)
	require.Equal(t, want, data)
}

// dec after := 155, 32, 0, 0, 161, 108, 84, 135
func TestGameCryptDecrypt(t *testing.T) {
	xorDec := []byte{114, 32, 0, 0, 161, 108, 84, 135}

	data := []byte{56, 24, 24, 24, 185, 208, 132, 3, 113, 2, 2, 91, 250, 197, 145, 44, 94, 126, 126, 43, 138, 149, 193, 35, 81, 81, 81, 126, 223, 156, 200, 39, 85, 26, 26, 119, 214, 223, 139, 12, 126}
	want := []byte{74, 0, 0, 0, 0, 5, 0, 0, 0, 83, 0, 89, 0, 83, 0, 58, 0, 0, 0, 85, 0, 115, 0, 101, 0, 32, 0, 47, 0, 47, 0, 104, 0, 111, 0, 109, 0, 101, 0, 0, 0}

	gameCrypt := NewGameCrypt()
	gameCrypt.SetEnabled()
	gameCrypt.KeyHack(nil, xorDec)

	gameCrypt.Decrypt(data, 0, 41)
	require.Equal(t, want, data)
}
