package prot

import (
	"github.com/Pesekjak/173go/pkg/buff"
	"github.com/Pesekjak/173go/pkg/world"
)

func init() {
	RegisterOut(0x00, &PacketOutKeepAlive{})
	RegisterOut(0x01, &PacketOutLogin{})
	RegisterOut(0x02, &PacketOutHandShake{})
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
	Dimension  world.Dimension
}

func (p *PacketOutLogin) Push(buf *buff.MCWriter) error {
	pusher := newPusher(buf)
	pusher.push(func() error { return buf.WriteInt(p.EntityId) })
	pusher.push(func() error { return buf.WriteString16(p.ServerName) })
	pusher.push(func() error { return buf.WriteLong(p.MapSeed) })
	pusher.push(func() error { return buf.WriteByte(byte(p.Dimension)) })
	return pusher.err
}

type PacketOutHandShake struct {
	Hash string
}

func (p *PacketOutHandShake) Push(buf *buff.MCWriter) error {
	pusher := newPusher(buf)
	pusher.push(func() error { return buf.WriteString16(p.Hash) })
	return pusher.err
}

type PacketOutKick struct {
	Reason string
}

func (p *PacketOutKick) Push(buf *buff.MCWriter) error {
	pusher := newPusher(buf)
	pusher.push(func() error { return buf.WriteString16(p.Reason) })
	return pusher.err
}
