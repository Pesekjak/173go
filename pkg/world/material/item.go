package material

// EquipSlot defines the slot where an item can be equipped.
type EquipSlot int

const (
	SlotHand EquipSlot = iota
	SlotHead
	SlotChest
	SlotLegs
	SlotFeet
)

func (i *Item) withMaxStackSize(size uint16) *Item {
	i.maxStackSize = size
	return i
}

func (i *Item) withMaxDamage(damage uint16) *Item {
	if damage > 0 {
		i.maxDamage = damage
		i.hasDurability = true
	}
	return i
}

func (i *Item) makeTool(durability uint16) *Item {
	i.maxStackSize = 1
	i.isTool = true
	return i.withMaxDamage(durability)
}

func (i *Item) makeFood() *Item {
	i.maxStackSize = 1
	i.isFood = true
	return i
}

func (i *Item) withEquipSlot(slot EquipSlot) *Item {
	i.equipSlot = slot
	return i
}

// MaxStackSize returns the maximum number of this item that can be in a single stack.
func (i *Item) MaxStackSize() uint16 {
	return i.maxStackSize
}

// MaxDamage returns the durability of the item. Returns 0 if the item doesn't have durability.
func (i *Item) MaxDamage() uint16 {
	return i.maxDamage
}

// HasDurability returns true if the item can be damaged.
func (i *Item) HasDurability() bool {
	return i.hasDurability
}

// IsTool returns true if the item is considered a tool.
func (i *Item) IsTool() bool {
	return i.isTool
}

// IsFood returns true if the item is a food source.
func (i *Item) IsFood() bool {
	return i.isFood
}

// EquipSlot returns the slot where the item can be equipped.
func (i *Item) EquipSlot() EquipSlot {
	return i.equipSlot
}
