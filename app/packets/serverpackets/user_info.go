package serverpackets

import (
	"github.com/sptGabriel/socks5/app/packets"
)

type UserInfo struct {
	X                       int32
	Y                       int32
	Z                       int32
	Heading                 int32
	ObjectID                int32
	Name                    string
	Race                    int32
	Sex                     int32
	ClassID                 int32
	Level                   int32
	Experience              int32
	STR                     int32
	DEX                     int32
	CON                     int32
	INT                     int32
	WIT                     int32
	MEN                     int32
	MaxHp                   int32
	CurrentHp               int32
	MaxMp                   int32
	CurrentMp               int32
	SP                      int32
	CurrentLoad             int32
	MaxLoad                 int32
	WeaponEquipped          int32
	PaperdollObjectID       [12]int32
	PaperdollItemID         [12]int32
	PAtk                    int32
	PAtkSpd                 int32
	PDef                    int32
	EvasionRate             int32
	Accuracy                int32
	CriticalHit             int32
	MAtk                    int32
	MAtkSpd                 int32
	MDef                    int32
	PvpFlag                 int32
	Karma                   int32
	RunSpd                  int32
	WalkSpd                 int32
	SwimRunSpd              int32
	SwimWalkSpd             int32
	FlyRunSpd               int32
	FlyWalkSpd              int32
	MovementSpeedMultiplier float32
	AttackSpeedMultiplier   float32
	CollisionRadius         float32
	CollisionHeight         float32
	HairStyle               int32
	HairColor               int32
	Face                    int32
	Title                   string
	ClanID                  int32
	ClanCrestID             int32
	AllyID                  int32
	AllyCrestID             int32
	Relation                int32
	MountType               int32
	PrivateStoreType        int32
	HasDwarvenCraft         bool
	PkKills                 int32
	PvpKills                int32
	Cubics                  []int16
	LookingForParty         bool
	AbnormalEffect          int32
	ClanPrivileges          int32
	Recommendations         struct {
		Left int16
		Have int16
	}
	MountNpcID       int32
	InventoryLimit   int16
	MaxCp            int32
	CurrentCp        int32
	EnchantEffect    int8
	EventTeam        int32
	AuraColor        int32
	ClanCrestLargeID int32
	IsNoble          bool
	IsHero           bool
	IsFishing        bool
	FishingLocation  struct {
		X int32
		Y int32
		Z int32
	}
	NameColor int32
}

func ReadUserInfo(reader *packets.Reader) (*UserInfo, error) {
	var userInfo UserInfo

	userInfo.X = int32(reader.ReadUInt32())
	userInfo.Y = int32(reader.ReadUInt32())
	userInfo.Z = int32(reader.ReadUInt32())
	userInfo.Heading = int32(reader.ReadUInt32())
	userInfo.ObjectID = int32(reader.ReadUInt32())
	userInfo.Name = reader.ReadString()

	return &userInfo, nil
}
