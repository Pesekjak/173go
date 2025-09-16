package prot

import (
	"github.com/Pesekjak/173go/pkg/buff"
)

func init() {
	RegisterOut(0x00, &PacketOutKeepAlive{})
	RegisterOut(0x01, &PacketOutLogin{})
	RegisterOut(0x02, &PacketOutHandShake{})
	RegisterOut(0x04, &PacketOutTimeUpdate{})
	RegisterOut(0x06, &PacketOutSpawnPosition{})
	RegisterOut(0x0D, &PacketOutPlayerPositionAndLook{})
	RegisterOut(0x32, &PacketOutPreChunk{})
	RegisterOut(0x33, &PacketOutMapChunk{})
	RegisterOut(0xFF, &PacketOutKick{})
}

type PacketOutKeepAlive struct {
}

func (p *PacketOutKeepAlive) Push(*buff.MCWriter) error {
	return nil
}

type PacketOutLogin struct {
	EntityId   int32
	ServerName string
	MapSeed    int64
	Dimension  byte
}

func (p *PacketOutLogin) Push(buf *buff.MCWriter) error {
	pusher := buff.NewPusher(buf)
	pusher.Push(func() error { return buf.WriteInt(p.EntityId) })
	pusher.Push(func() error { return buf.WriteString16(p.ServerName) })
	pusher.Push(func() error { return buf.WriteLong(p.MapSeed) })
	pusher.Push(func() error { return buf.WriteByte(p.Dimension) })
	return pusher.Err
}

type PacketOutHandShake struct {
	Hash string
}

func (p *PacketOutHandShake) Push(buf *buff.MCWriter) error {
	pusher := buff.NewPusher(buf)
	pusher.Push(func() error { return buf.WriteString16(p.Hash) })
	return pusher.Err
}

type PacketOutTimeUpdate struct {
	Time int64
}

func (p *PacketOutTimeUpdate) Push(buf *buff.MCWriter) error {
	pusher := buff.NewPusher(buf)
	pusher.Push(func() error { return buf.WriteLong(p.Time) })
	return pusher.Err
}

type PacketOutSpawnPosition struct {
	X int32
	Y int32
	Z int32
}

func (p *PacketOutSpawnPosition) Push(buf *buff.MCWriter) error {
	pusher := buff.NewPusher(buf)
	pusher.Push(func() error { return buf.WriteInt(p.X) })
	pusher.Push(func() error { return buf.WriteInt(p.Y) })
	pusher.Push(func() error { return buf.WriteInt(p.Z) })
	return pusher.Err
}

type PacketOutPlayerPositionAndLook struct {
	X        float64
	Stance   float64
	Y        float64
	Z        float64
	Yaw      float32
	Pitch    float32
	OnGround bool
}

func (p *PacketOutPlayerPositionAndLook) Push(buf *buff.MCWriter) error {
	pusher := buff.NewPusher(buf)
	pusher.Push(func() error { return buf.WriteDouble(p.X) })
	pusher.Push(func() error { return buf.WriteDouble(p.Stance) })
	pusher.Push(func() error { return buf.WriteDouble(p.Y) })
	pusher.Push(func() error { return buf.WriteDouble(p.Z) })
	pusher.Push(func() error { return buf.WriteFloat(p.Yaw) })
	pusher.Push(func() error { return buf.WriteFloat(p.Pitch) })
	pusher.Push(func() error { return buf.WriteBool(p.OnGround) })
	return pusher.Err
}

type PacketOutPreChunk struct {
	X    int32
	Z    int32
	Load bool
}

func (p *PacketOutPreChunk) Push(buf *buff.MCWriter) error {
	pusher := buff.NewPusher(buf)
	pusher.Push(func() error { return buf.WriteInt(p.X) })
	pusher.Push(func() error { return buf.WriteInt(p.Z) })
	pusher.Push(func() error { return buf.WriteBool(p.Load) })
	return pusher.Err
}

type PacketOutMapChunk struct {
	X     int32
	Y     int16
	Z     int32
	SizeX byte
	SizeY byte
	SizeZ byte
	Data  []byte
}

func (p *PacketOutMapChunk) Push(buf *buff.MCWriter) error {
	pusher := buff.NewPusher(buf)
	pusher.Push(func() error { return buf.WriteInt(p.X) })
	pusher.Push(func() error { return buf.WriteShort(p.Y) })
	pusher.Push(func() error { return buf.WriteInt(p.Z) })
	pusher.Push(func() error { return buf.WriteByte(p.SizeX) })
	pusher.Push(func() error { return buf.WriteByte(p.SizeY) })
	pusher.Push(func() error { return buf.WriteByte(p.SizeZ) })
	pusher.Push(func() error { return buf.WriteInt(int32(len(p.Data))) })
	pusher.Push(func() error { return buf.WriteBytes(p.Data) })
	return pusher.Err
}

type PacketOutKick struct {
	Reason string
}

func (p *PacketOutKick) Push(buf *buff.MCWriter) error {
	pusher := buff.NewPusher(buf)
	pusher.Push(func() error { return buf.WriteString16(p.Reason) })
	return pusher.Err
}
