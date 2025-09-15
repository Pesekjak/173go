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
	puller := newPuller(buf)
	puller.pull(func() { p.Protocol, puller.err = buf.ReadInt() })
	puller.pull(func() { p.Username, puller.err = buf.ReadString16() })
	puller.pull(func() { p.MapSeed, puller.err = buf.ReadLong() })
	puller.pull(func() { p.Dimension, puller.err = buf.ReadByte() })
	return puller.err
}

func (p *PacketInLogin) Handle(handler PacketHandler) error {
	return handler.OnLogin(p)
}

type PacketInHandShake struct {
	Username string
}

func (p *PacketInHandShake) Pull(buf *buff.MCReader) error {
	puller := newPuller(buf)
	puller.pull(func() { p.Username, puller.err = buf.ReadString16() })
	return puller.err
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
	puller := newPuller(buf)
	puller.pull(func() { p.X, puller.err = buf.ReadDouble() })
	puller.pull(func() { p.Y, puller.err = buf.ReadDouble() })
	puller.pull(func() { p.Stance, puller.err = buf.ReadDouble() })
	puller.pull(func() { p.Z, puller.err = buf.ReadDouble() })
	puller.pull(func() { p.OnGround, puller.err = buf.ReadBool() })
	return puller.err
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
	puller := newPuller(buf)
	puller.pull(func() { p.X, puller.err = buf.ReadDouble() })
	puller.pull(func() { p.Y, puller.err = buf.ReadDouble() })
	puller.pull(func() { p.Stance, puller.err = buf.ReadDouble() })
	puller.pull(func() { p.Z, puller.err = buf.ReadDouble() })
	puller.pull(func() { p.Yaw, puller.err = buf.ReadFloat() })
	puller.pull(func() { p.Pitch, puller.err = buf.ReadFloat() })
	puller.pull(func() { p.OnGround, puller.err = buf.ReadBool() })
	return puller.err
}

func (p *PacketInPlayerPositionAndLook) Handle(handler PacketHandler) error {
	return handler.OnPlayerPositionAndLook(p)
}
