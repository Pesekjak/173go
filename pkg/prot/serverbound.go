package prot

import (
	"github.com/Pesekjak/173go/pkg/buff"
)

func init() {
	RegisterIn(0x00, func() PacketIn { return &PacketInKeepAlive{} })
	RegisterIn(0x01, func() PacketIn { return &PacketInLogin{} })
	RegisterIn(0x02, func() PacketIn { return &PacketInHandShake{} })
	RegisterIn(0x0B, func() PacketIn { return &PacketInPlayerPosition{} })
	RegisterIn(0x0D, func() PacketIn { return &PacketInPlayerPositionAndLook{} })
}

type PacketInKeepAlive struct {
}

func (p *PacketInKeepAlive) Pull(*buff.MCReader) error {
	return nil
}

func (p *PacketInKeepAlive) Handle(handler PacketHandler) error {
	return handler.OnKeepAlive(p)
}

type PacketInLogin struct {
	Protocol  int32
	Username  string
	MapSeed   int64
	Dimension byte
}

func (p *PacketInLogin) Pull(buf *buff.MCReader) error {
	puller := buff.NewPuller(buf)
	puller.Pull(func() { p.Protocol, puller.Err = buf.ReadInt() })
	puller.Pull(func() { p.Username, puller.Err = buf.ReadString16() })
	puller.Pull(func() { p.MapSeed, puller.Err = buf.ReadLong() })
	puller.Pull(func() { p.Dimension, puller.Err = buf.ReadByte() })
	return puller.Err
}

func (p *PacketInLogin) Handle(handler PacketHandler) error {
	return handler.OnLogin(p)
}

type PacketInHandShake struct {
	Username string
}

func (p *PacketInHandShake) Pull(buf *buff.MCReader) error {
	puller := buff.NewPuller(buf)
	puller.Pull(func() { p.Username, puller.Err = buf.ReadString16() })
	return puller.Err
}

func (p *PacketInHandShake) Handle(handler PacketHandler) error {
	return handler.OnHandShake(p)
}

type PacketInPlayerPosition struct {
	X        float64
	Y        float64
	Stance   float64
	Z        float64
	OnGround bool
}

func (p *PacketInPlayerPosition) Pull(buf *buff.MCReader) error {
	puller := buff.NewPuller(buf)
	puller.Pull(func() { p.X, puller.Err = buf.ReadDouble() })
	puller.Pull(func() { p.Y, puller.Err = buf.ReadDouble() })
	puller.Pull(func() { p.Stance, puller.Err = buf.ReadDouble() })
	puller.Pull(func() { p.Z, puller.Err = buf.ReadDouble() })
	puller.Pull(func() { p.OnGround, puller.Err = buf.ReadBool() })
	return puller.Err
}

func (p *PacketInPlayerPosition) Handle(handler PacketHandler) error {
	return handler.OnPlayerPosition(p)
}

type PacketInPlayerPositionAndLook struct {
	X        float64
	Y        float64
	Stance   float64
	Z        float64
	Yaw      float32
	Pitch    float32
	OnGround bool
}

func (p *PacketInPlayerPositionAndLook) Pull(buf *buff.MCReader) error {
	puller := buff.NewPuller(buf)
	puller.Pull(func() { p.X, puller.Err = buf.ReadDouble() })
	puller.Pull(func() { p.Y, puller.Err = buf.ReadDouble() })
	puller.Pull(func() { p.Stance, puller.Err = buf.ReadDouble() })
	puller.Pull(func() { p.Z, puller.Err = buf.ReadDouble() })
	puller.Pull(func() { p.Yaw, puller.Err = buf.ReadFloat() })
	puller.Pull(func() { p.Pitch, puller.Err = buf.ReadFloat() })
	puller.Pull(func() { p.OnGround, puller.Err = buf.ReadBool() })
	return puller.Err
}

func (p *PacketInPlayerPositionAndLook) Handle(handler PacketHandler) error {
	return handler.OnPlayerPositionAndLook(p)
}
