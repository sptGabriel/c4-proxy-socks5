package app

import (
	"testing"
)

func TestXxx(t *testing.T) {
	// const (
	// 	xorKeySize   = 4
	// 	checkSumSize = 4

	// 	xorPadding      = 4
	// 	blowFishPadding = 8
	// )

	// encrypted := []byte{34, 0, 155, 92, 123, 104, 119, 95, 68, 78, 50, 112, 33, 57, 75, 62, 170, 239, 157, 65, 190, 111, 32, 196, 50, 31, 17, 46, 105, 175, 152, 183, 218, 106}

	// cryptKey := []byte{
	// 	0x5f, 0x3b, 0x35, 0x2e,
	// 	0x5d, 0x39, 0x34, 0x2d,
	// 	0x33, 0x31, 0x3d, 0x3d,
	// 	0x2d, 0x25, 0x78, 0x54,
	// 	0x21, 0x5e, 0x5b, 0x24, 0x00,
	// }

	// cryptSvc, err := NewCrypt()
	// require.NoError(t, err)

	// toDecrypt := encrypted[2:]
	// err = cryptSvc.Decrypt(toDecrypt, cryptKey)
	// fmt.Print(err, toDecrypt)
}

// func TestEncrypt(t *testing.T) {
// 	want := []byte{42, 0, 71, 229, 156, 239, 211, 38, 107, 132, 238, 118, 217, 88, 202, 233, 202, 100, 86, 165, 77, 95, 89, 234, 223, 79, 65, 97, 111, 92, 75, 189, 111, 165, 164, 56, 26, 226, 111, 153, 82, 85}
// 	nonEncryptedData := []byte{7, 170, 7, 98, 67, 35, 1, 0, 0, 103, 69, 0, 0, 171, 137, 0, 0, 239, 205, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// 	cryptSvc, err := NewLoginCrypt()
// 	require.NoError(t, err)

// 	encryptedData := cryptSvc.Encrypt(nonEncryptedData, 0, 40)

// 	packetSize := len(nonEncryptedData)
// 	sendable := make([]byte, packetSize+2)
// 	sendable[0] = byte((packetSize + 2) & 0xff)
// 	sendable[1] = byte((packetSize + 2) >> 8)

// 	copy(sendable[2:], encryptedData[:packetSize])

// 	require.Equal(t, want, sendable)
// }

// func TestDecrypt(t *testing.T) {
// 	want := []byte{11, 8, 217, 120, 187, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 176, 1, 217, 120, 0, 0, 0, 0}
// 	encryptedData := []byte{137, 65, 69, 189, 6, 93, 196, 235, 65, 97, 111, 92, 75, 189, 111, 165, 141, 53, 105, 20, 90, 215, 183, 193}

// 	cryptSvc, err := NewLoginCrypt()
// 	require.NoError(t, err)

// 	err = cryptSvc.Decrypt(encryptedData, 0, 24)
// 	require.NoError(t, err)

// 	require.Equal(t, want, encryptedData)
// }
