package inventory

type Type byte

const (
	Chest Type = iota
	Crafting
	Furnace
	Dispenser
)

type Window interface {
	Type() Type
	Name() string
	Size() byte
}
