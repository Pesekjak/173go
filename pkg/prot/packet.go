package prot

import (
	"fmt"
	"io"
	"reflect"

	"github.com/Pesekjak/173go/pkg/buff"
)

type PacketIn interface {
	Pull(buf *buff.MCReader) error
	Handle(handler PacketHandler) error
}

type PacketOut interface {
	Push(buf *buff.MCWriter) error
}

var inRegistry = make(map[byte]func() PacketIn)
var outRegistry = make(map[reflect.Type]byte)

func RegisterIn(packetID byte, factory func() PacketIn) {
	if _, exists := inRegistry[packetID]; exists {
		panic(fmt.Sprintf("incoming packet with ID 0x%02X is already registered", packetID))
	}
	inRegistry[packetID] = factory
}

func RegisterOut(packetID byte, packet PacketOut) {
	packetType := reflect.TypeOf(packet)
	if _, exists := outRegistry[packetType]; exists {
		panic(fmt.Sprintf("outgoing packet type %v is already registered", packetType))
	}
	outRegistry[packetType] = packetID
}

func GetPacket(packetID byte) (PacketIn, error) {
	factory, ok := inRegistry[packetID]
	if !ok {
		return nil, fmt.Errorf("unknown packet ID: 0x%02X", packetID)
	}
	return factory(), nil
}

func GetPacketID(packet PacketOut) (byte, error) {
	packetType := reflect.TypeOf(packet)
	id, ok := outRegistry[packetType]
	if !ok {
		return 0, fmt.Errorf("unknown outgoing packet type: %T", packet)
	}
	return id, nil
}

func ReadPacket(packetID byte, reader io.Reader) (PacketIn, error) {
	packet, err := GetPacket(packetID)
	if err != nil {
		return nil, fmt.Errorf("received packet with unknown ID: 0x%02X", packetID)
	}
	if err = packet.Pull(buff.NewReader(reader)); err != nil {
		return nil, err
	}
	return packet, nil
}

func WritePacket(packet PacketOut, writer io.Writer) error {
	buf := buff.NewWriter(writer)
	id, err := GetPacketID(packet)
	if err != nil {
		return err
	}
	if err = buf.WriteByte(id); err != nil {
		return err
	}
	if err = packet.Push(buf); err != nil {
		return err
	}
	return nil
}
