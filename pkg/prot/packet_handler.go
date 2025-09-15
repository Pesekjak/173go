package prot

type PacketHandler interface {
	OnKeepAlive(packet *PacketInKeepAlive) error
	OnLogin(packet *PacketInLogin) error
	OnHandShake(packet *PacketInHandShake) error
	OnPlayerPosition(packet *PacketInPlayerPosition) error
	OnPlayerPositionAndLook(packet *PacketInPlayerPositionAndLook) error
}
