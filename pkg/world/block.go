package world

import "github.com/Pesekjak/173go/pkg/world/material"

type Block interface {
	Position() BlockPos
	Material() *material.Block
	Data() byte

	Set(material *material.Block, data byte) error
}
