package inventory

import (
	"strconv"

	"github.com/Pesekjak/173go/pkg/world/material"
)

type ItemStack struct {
	Material material.Material
	Count    byte
	Data     uint16
}

func NewItemStack(material material.Material, count byte, data uint16) ItemStack {
	return ItemStack{
		Material: material,
		Count:    count,
		Data:     data,
	}
}

func EmptyStack() ItemStack {
	return NewItemStack(material.Air, 0, 0)
}

func (s ItemStack) IsEmpty() bool {
	return s.Count == 0 || s.Material == material.Air
}

func (s ItemStack) String() string {
	return strconv.Itoa(int(s.Count)) + "x " + s.Material.String()
}
