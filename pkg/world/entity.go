package world

import (
	"fmt"

	"github.com/Pesekjak/173go/pkg/net"
	"github.com/Pesekjak/173go/pkg/world/entity_data"
)

type EntityType byte

const (
	Creeper = iota
	Skeleton
	Spider
	GiantZombie
	Zombie
	Slime
	Ghast
	ZombiePigman
	Pig
	Sheep
	Cow
	Hen
	Squid
	Wolf

	Boat
	Minecart
	StorageCart
	PoweredCart
	ActivatedTNT
	Arrow
	ThrownSnowball
	ThrownEgg
	FallingSand
	FallingGravel
	FishingFloat

	Player

	Item

	Painting

	Lightning
)

func IsMobType(entityType EntityType) bool {
	return entityType <= 13
}

func IsObjectType(entityType EntityType) bool {
	return entityType >= 14 && entityType <= 24
}

var entityIdMap = map[EntityType]int{
	Creeper:      50,
	Skeleton:     51,
	Spider:       52,
	GiantZombie:  53,
	Zombie:       54,
	Slime:        55,
	Ghast:        56,
	ZombiePigman: 57,
	Pig:          90,
	Sheep:        91,
	Cow:          92,
	Hen:          93,
	Squid:        94,
	Wolf:         95,

	Boat:           1,
	Minecart:       10,
	StorageCart:    11,
	PoweredCart:    12,
	ActivatedTNT:   20,
	Arrow:          60,
	ThrownSnowball: 61,
	ThrownEgg:      62,
	FallingSand:    70,
	FallingGravel:  71,
	FishingFloat:   90,
}

func GetEntityTypeId(entityType EntityType) (int, error) {
	id, ok := entityIdMap[entityType]
	if !ok {
		return -1, fmt.Errorf("entity type without id: %v", entityType)
	}
	return id, nil
}

type Entity interface {
	Id() int32
	Location() Location
	World() *World
	EntityType() EntityType
}

type MobEntity interface {
	Entity
	IsAlive() bool
	Health() uint32
	Metadata() entity_data.Metadata
}

type ObjectEntity interface {
	Entity
}

type PaintingEntity interface {
	Entity
	ArtType() entity_data.ArtType
}

type PlayerEntity interface {
	MobEntity
	Username() string
	IsOnline() bool
	Connection() *net.Connection
	Disconnect(error)
	Kick(string)
}
