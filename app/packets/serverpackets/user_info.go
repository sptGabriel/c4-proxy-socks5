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
	Appearence              int32
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
	FlRunSpd                int32
	FlWalkSpd               int32
	FlyRunSpd               int32
	FlyWalkSpd              int32
	UnderObjectID           int32
	RearObjectID            int32
	LearObjectID            int32
	NeckObjectID            int32
	RFingerObjectID         int32
	LFingerObjectID         int32
	HeadObjectID            int32
	RHandObjectID           int32
	LHandObjectID           int32
	GlovesObjectID          int32
	ChestObjectID           int32
	LegsObjectID            int32
	BackObjectID            int32
	FeetObjectID            int32
	HairObjectID            int32
	UnderItemID             int32
	RearItemID              int32
	LearItemID              int32
	NeckItemID              int32
	RFingerItemID           int32
	LFingerItemID           int32
	HeadItemID              int32
	RHandItemID             int32
	LHandItemID             int32
	GlovesItemID            int32
	ChestItemID             int32
	LegsItemID              int32
	FeetItemID              int32
	BackItemID              int32
	HairItemID              int32
	MovementSpeedMultiplier float64
	AttackSpeedMultiplier   float64
	CollisionRadius         float64
	CollisionHeight         float64
	HairStyle               int32
	HairColor               int32
	Face                    int32
	AccessLevel             int32
	Title                   string
	ActiveWeapon            int32
	ClanID                  int32
	ClanCrestID             int32
	AllyID                  int32
	AllyCrestID             int32
	Relation                int32
	MountType               int32
	PrivateStoreType        int32
	HasDwarvenCraft         int32
	PkKills                 int32
	PvpKills                int32
	NumCubics               uint16
	Cubics                  []int16
	LookingForParty         bool
	AbnormalEffect          int32
	ClanPrivileges          int32
	Recommendations         struct {
		Left int32
		Have int32
	}
	MountNpcID       int32
	InventoryLimit   int32
	MaxCp            int32
	CurrentCp        int32
	EnchantEffect    int32
	EventTeam        int32
	AuraColor        int32
	ClanCrestLargeID int32
	IsMounted        int32
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
	userInfo.Race = int32(reader.ReadUInt32())
	userInfo.Sex = int32(reader.ReadUInt32())
	userInfo.ClassID = int32(reader.ReadUInt32())
	userInfo.Level = int32(reader.ReadUInt32())
	userInfo.Experience = int32(reader.ReadUInt32())
	userInfo.STR = int32(reader.ReadUInt32())
	userInfo.DEX = int32(reader.ReadUInt32())
	userInfo.CON = int32(reader.ReadUInt32())
	userInfo.INT = int32(reader.ReadUInt32())
	userInfo.WIT = int32(reader.ReadUInt32())
	userInfo.MEN = int32(reader.ReadUInt32())
	userInfo.MaxHp = int32(reader.ReadUInt32())
	userInfo.CurrentHp = int32(reader.ReadUInt32())
	userInfo.MaxMp = int32(reader.ReadUInt32())
	userInfo.CurrentMp = int32(reader.ReadUInt32())
	userInfo.SP = int32(reader.ReadUInt32())
	// Carga e armas equipadas
	userInfo.CurrentLoad = int32(reader.ReadUInt32())
	userInfo.MaxLoad = int32(reader.ReadUInt32())
	userInfo.ActiveWeapon = int32(reader.ReadUInt32()) // 20 ou 40 dependendo da arma equipada

	// IDs dos objetos equipados
	userInfo.UnderObjectID = int32(reader.ReadUInt32())
	userInfo.RearObjectID = int32(reader.ReadUInt32())
	userInfo.LearObjectID = int32(reader.ReadUInt32())
	userInfo.NeckObjectID = int32(reader.ReadUInt32())
	userInfo.RFingerObjectID = int32(reader.ReadUInt32())
	userInfo.LFingerObjectID = int32(reader.ReadUInt32())
	userInfo.HeadObjectID = int32(reader.ReadUInt32())
	userInfo.RHandObjectID = int32(reader.ReadUInt32())
	userInfo.LHandObjectID = int32(reader.ReadUInt32())
	userInfo.GlovesObjectID = int32(reader.ReadUInt32())
	userInfo.ChestObjectID = int32(reader.ReadUInt32())
	userInfo.LegsObjectID = int32(reader.ReadUInt32())
	userInfo.FeetObjectID = int32(reader.ReadUInt32())
	userInfo.BackObjectID = int32(reader.ReadUInt32())
	userInfo.RHandObjectID = int32(reader.ReadUInt32())
	userInfo.HairObjectID = int32(reader.ReadUInt32())

	// IDs dos itens equipados
	userInfo.UnderItemID = int32(reader.ReadUInt32())
	userInfo.RearItemID = int32(reader.ReadUInt32())
	userInfo.LearItemID = int32(reader.ReadUInt32())
	userInfo.NeckItemID = int32(reader.ReadUInt32())
	userInfo.RFingerItemID = int32(reader.ReadUInt32())
	userInfo.LFingerItemID = int32(reader.ReadUInt32())
	userInfo.HeadItemID = int32(reader.ReadUInt32())
	userInfo.RHandItemID = int32(reader.ReadUInt32())
	userInfo.LHandItemID = int32(reader.ReadUInt32())
	userInfo.GlovesItemID = int32(reader.ReadUInt32())
	userInfo.ChestItemID = int32(reader.ReadUInt32())
	userInfo.LegsItemID = int32(reader.ReadUInt32())
	userInfo.FeetItemID = int32(reader.ReadUInt32())
	userInfo.BackItemID = int32(reader.ReadUInt32())
	userInfo.RHandItemID = int32(reader.ReadUInt32())
	userInfo.HairItemID = int32(reader.ReadUInt32())

	// Estatísticas de ataque e defesa
	userInfo.PAtk = int32(reader.ReadUInt32())
	userInfo.PAtkSpd = int32(reader.ReadUInt32())
	userInfo.PDef = int32(reader.ReadUInt32())
	userInfo.EvasionRate = int32(reader.ReadUInt32())
	userInfo.Accuracy = int32(reader.ReadUInt32())
	userInfo.CriticalHit = int32(reader.ReadUInt32())
	userInfo.MAtk = int32(reader.ReadUInt32())
	userInfo.MAtkSpd = int32(reader.ReadUInt32())
	userInfo.PAtkSpd = int32(reader.ReadUInt32())
	userInfo.MDef = int32(reader.ReadUInt32())

	// PvP e Karma
	userInfo.PvpFlag = int32(reader.ReadUInt32())
	userInfo.Karma = int32(reader.ReadUInt32())
	// Velocidade de movimento
	userInfo.RunSpd = int32(reader.ReadUInt32())
	userInfo.WalkSpd = int32(reader.ReadUInt32())
	userInfo.SwimRunSpd = int32(reader.ReadUInt32())
	userInfo.SwimWalkSpd = int32(reader.ReadUInt32())
	userInfo.FlRunSpd = int32(reader.ReadUInt32())
	userInfo.FlWalkSpd = int32(reader.ReadUInt32())
	userInfo.FlyRunSpd = int32(reader.ReadUInt32())
	userInfo.FlyWalkSpd = int32(reader.ReadUInt32())
	userInfo.MovementSpeedMultiplier = reader.ReadFloat64()
	userInfo.AttackSpeedMultiplier = reader.ReadFloat64()

	userInfo.CollisionRadius = reader.ReadFloat64()
	userInfo.CollisionHeight = reader.ReadFloat64()

	userInfo.HairStyle = int32(reader.ReadUInt32())
	userInfo.HairColor = int32(reader.ReadUInt32())
	userInfo.Face = int32(reader.ReadUInt32())
	userInfo.AccessLevel = int32(reader.ReadUInt32())
	userInfo.Title = reader.ReadString()

	userInfo.ClanID = int32(reader.ReadUInt32())
	userInfo.ClanCrestID = int32(reader.ReadUInt32())
	userInfo.AllyID = int32(reader.ReadUInt32())
	userInfo.AllyCrestID = int32(reader.ReadUInt32())

	userInfo.Relation = int32(reader.ReadUInt32())
	userInfo.MountType = int32(reader.ReadUInt8())
	userInfo.PrivateStoreType = int32(reader.ReadUInt8())
	userInfo.HasDwarvenCraft = int32(reader.ReadUInt8())
	userInfo.PkKills = int32(reader.ReadUInt32())
	userInfo.PvpKills = int32(reader.ReadUInt32())

	numCubics := reader.ReadUInt16()
	cubicIDs := make([]int16, numCubics)
	for i := 0; i < int(numCubics); i++ {
		cubicIDs[i] = int16(reader.ReadUInt16())
	}

	userInfo.NumCubics = numCubics
	userInfo.Cubics = cubicIDs

	lookingForParty := int32(reader.ReadUInt8())
	if lookingForParty == 1 {
		userInfo.LookingForParty = true
	}

	userInfo.AbnormalEffect = int32(reader.ReadUInt32())
	_ = int32(reader.ReadUInt8())
	userInfo.ClanPrivileges = int32(reader.ReadUInt32())
	_ = int32(reader.ReadUInt32())
	_ = int32(reader.ReadUInt32())
	_ = int32(reader.ReadUInt32())
	_ = int32(reader.ReadUInt32())
	_ = int32(reader.ReadUInt32())
	_ = int32(reader.ReadUInt32())
	_ = int32(reader.ReadUInt32())

	// Lê as recomendações restantes e recebidas
	userInfo.Recommendations.Left = int32(reader.ReadUInt16())
	userInfo.Recommendations.Have = int32(reader.ReadUInt16())

	// Lê o ID do NPC da montaria (somado a 1000000 no Java)
	userInfo.MountNpcID = int32(reader.ReadUInt32())

	// Lê o limite de inventário
	userInfo.InventoryLimit = int32(reader.ReadUInt16())

	// Lê o ID da classe
	userInfo.ClassID = int32(reader.ReadUInt32())

	_ = int32(reader.ReadUInt32())

	// Lê o CP máximo e o CP atual
	userInfo.MaxCp = int32(reader.ReadUInt32())
	userInfo.CurrentCp = int32(reader.ReadUInt32())

	// Lê o efeito de encantamento ou se o personagem está montado
	userInfo.EnchantEffect = int32(reader.ReadUInt8())

	// Lê o time ou a cor da aura
	userInfo.EventTeam = int32(reader.ReadUInt8())

	// Lê o ID do brasão do clã
	userInfo.ClanCrestLargeID = int32(reader.ReadUInt32())

	// Lê se é nobre e se tem a aura de herói
	isNoble := reader.ReadUInt8()
	if isNoble == 1 {
		userInfo.IsNoble = true
	}

	isHero := reader.ReadUInt8()
	if isHero == 1 {
		userInfo.IsHero = true
	}

	// Lê se está pescando
	isFishing := reader.ReadUInt8()
	if isFishing == 1 {
		userInfo.IsFishing = true
	}

	// Lê as coordenadas da pesca (X, Y, Z)
	userInfo.FishingLocation.X = int32(reader.ReadUInt32())
	userInfo.FishingLocation.Y = int32(reader.ReadUInt32())
	userInfo.FishingLocation.Z = int32(reader.ReadUInt32())

	// Lê a cor do nome
	userInfo.NameColor = int32(reader.ReadUInt32())

	return &userInfo, nil
}
