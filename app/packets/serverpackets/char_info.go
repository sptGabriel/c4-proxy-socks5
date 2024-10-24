package serverpackets

import (
	"fmt"

	"github.com/sptGabriel/socks5/app/packets"
)

type CharInfo struct {
	X             int32
	Y             int32
	Z             int32
	Heading       int32
	ObjectID      int32
	Name          string
	UnderItemID   int32
	RearItemID    int32
	LearItemID    int32
	NeckItemID    int32
	RFingerItemID int32
	LFingerItemID int32
	HeadItemID    int32
	RHandItemID   int32
	LHandItemID   int32
	GlovesItemID  int32
	ChestItemID   int32
	LegsItemID    int32
	FeetItemID    int32
	BackItemID    int32
	HairItemID    int32
}

func ReadCharInfo(reader *packets.Reader) (*CharInfo, error) {
	var charInfo CharInfo

	charInfo.X = int32(reader.ReadUInt32())
	charInfo.Y = int32(reader.ReadUInt32())
	charInfo.Z = int32(reader.ReadUInt32())
	charInfo.Heading = int32(reader.ReadUInt32())
	charInfo.ObjectID = int32(reader.ReadUInt32())
	charInfo.Name = reader.ReadString()
	fmt.Println("name:", charInfo.Name)
	_ = int32(reader.ReadUInt32()) // 204
	_ = int32(reader.ReadUInt32()) // 205
	_ = int32(reader.ReadUInt32()) // 206
	fmt.Println("test4", reader.Position())
	charInfo.UnderItemID = int32(reader.ReadUInt32())
	fmt.Println("under item:", charInfo.UnderItemID)
	charInfo.HeadItemID = int32(reader.ReadUInt32()) // 209
	fmt.Println("HeadItemID", reader.Position())
	charInfo.RHandItemID = int32(reader.ReadUInt32())  // 210
	charInfo.LHandItemID = int32(reader.ReadUInt32())  // 211
	charInfo.GlovesItemID = int32(reader.ReadUInt32()) // 212
	charInfo.ChestItemID = int32(reader.ReadUInt32())  // 213
	charInfo.LegsItemID = int32(reader.ReadUInt32())   // 214
	fmt.Println("LegsItemID", reader.Position())
	charInfo.FeetItemID = int32(reader.ReadUInt32()) // 215
	fmt.Println("FeetItemID", reader.Position())
	charInfo.BackItemID = int32(reader.ReadUInt32()) // 216
	fmt.Println("BackItemID", reader.Position())
	charInfo.RHandItemID = int32(reader.ReadUInt32()) // 217
	fmt.Println("RHandItemID", reader.Position())
	charInfo.HairItemID = int32(reader.ReadUInt32()) // 218
	fmt.Println("HairItemID", reader.Position())

	return &charInfo, nil
}
