package loot

import (
	"math/rand"

	"github.com/Pesekjak/173go/pkg/world/inventory"
	"github.com/Pesekjak/173go/pkg/world/material"
)

// GetDrops calculates the item drops for a given block.
// It accepts the block's material, its metadata, and a random source.
func GetDrops(block *material.Block, metadata byte, r *rand.Rand) []inventory.ItemStack {
	tries := getBlockLootTries(block, r)
	var drops = make([]inventory.ItemStack, tries)

	for i := 0; i < tries; i++ {
		// check if this specific drop attempt is successful
		if r.Float32() <= getBlockLootChance(block, metadata, i) {
			stack := getBlockLootStack(block, metadata, i, r)
			if !stack.IsEmpty() {
				drops = append(drops, stack)
			}
		}
	}

	return drops
}

func getBlockLootTries(b *material.Block, r *rand.Rand) int {
	switch b {
	case material.Air, material.Bookshelf, material.CakeBlock, material.Fire,
		material.WaterFlowing, material.WaterStill, material.LavaFlowing, material.LavaStill,
		material.Glass, material.Ice, material.MobSpawner, material.PistonHead,
		material.Portal, material.SnowLayer, material.TNT:
		return 0
	case material.ClayBlock:
		return 4
	case material.Seeds:
		return 4 // 1 for wheat item + 3 for seeds
	case material.GlowstoneBlock:
		return 2 + r.Intn(3)
	case material.Leaves:
		if r.Intn(20) == 0 {
			return 1 // 1 in 20 chance to attempt a sapling drop
		}
		return 0
	case material.LapisLazuliOre:
		return 4 + r.Intn(5)
	case material.RedstoneOre, material.RedstoneOreGlowing:
		return 4 + r.Intn(2)
	case material.SnowBlock:
		return 4
	case material.DoubleSlab:
		return 2
	default:
		return 1
	}
}

func getBlockLootChance(b *material.Block, metadata byte, tryNum int) float32 {
	switch b {
	case material.Seeds:
		if tryNum != 0 {
			// chance to drop extra seeds increases as wheat grows
			// fully grown (meta 7) has 50% per try
			return float32(metadata) / 14.0
		}
	}
	return 1
}

func getBlockLootStack(b *material.Block, metadata byte, tryNum int, r *rand.Rand) inventory.ItemStack {
	switch b {
	// blocks that don't drop
	case material.CakeBlock, material.DeadBush, material.PistonHead, material.MobSpawner:
		return inventory.EmptyStack()

	// special metadata-dependent drops
	case material.BedBlock:
		if (metadata & 0x8) != 0 { // head piece
			return inventory.EmptyStack()
		}
		return inventory.NewItemStack(material.BedItem, 1, 0)
	case material.WoodenDoorBlock:
		if (metadata & 0x8) != 0 { // upper part
			return inventory.EmptyStack()
		}
		return inventory.NewItemStack(material.WoodenDoorItem, 1, 0)
	case material.IronDoorBlock:
		if (metadata & 0x8) != 0 { // upper part
			return inventory.EmptyStack()
		}
		return inventory.NewItemStack(material.IronDoorItem, 1, 0)
	case material.Seeds:
		if tryNum == 0 {
			if metadata == 7 {
				return inventory.NewItemStack(material.Wheat, 1, 0)
			}
			return inventory.EmptyStack()
		}
		return inventory.NewItemStack(material.WheatSeeds, 1, 0)

	// blocks that drop different blocks
	case material.Farmland, material.GrassBlock:
		return inventory.NewItemStack(material.Dirt, 1, 0)
	case material.Stone:
		return inventory.NewItemStack(material.Cobblestone, 1, 0)
	case material.Furnace, material.FurnaceLit:
		return inventory.NewItemStack(material.Furnace, 1, 0)

	// ores and random drops
	case material.GlowstoneBlock:
		return inventory.NewItemStack(material.GlowstoneDust, 1, 0)
	case material.Gravel:
		if r.Intn(10) == 0 {
			return inventory.NewItemStack(material.Flint, 1, 0)
		}
		return inventory.NewItemStack(material.Gravel, 1, 0)
	case material.CoalOre:
		return inventory.NewItemStack(material.Coal, 1, 0)
	case material.DiamondOre:
		return inventory.NewItemStack(material.Diamond, 1, 0)
	case material.RedstoneOre, material.RedstoneOreGlowing:
		return inventory.NewItemStack(material.RedstoneDust, 1, 0)
	case material.LapisLazuliOre:
		return inventory.NewItemStack(material.Dye, 1, 4) // 4 for lapis
	case material.TallGrass:
		if r.Intn(8) == 0 {
			return inventory.NewItemStack(material.WheatSeeds, 1, 0)
		}
		return inventory.EmptyStack()

	// blocks that drop specific items
	case material.ClayBlock:
		return inventory.NewItemStack(material.ClayBalls, 1, 0)
	case material.Leaves:
		return inventory.NewItemStack(material.Sapling, 1, uint16(metadata&3))
	case material.RedstoneWire:
		return inventory.NewItemStack(material.RedstoneDust, 1, 0)
	case material.RedstoneRepeaterOff, material.RedstoneRepeaterOn:
		return inventory.NewItemStack(material.RedstoneRepeaterItem, 1, 0)
	case material.RedstoneTorchOff, material.RedstoneTorchOn:
		return inventory.NewItemStack(material.RedstoneTorchOn, 1, 0)
	case material.SugarCaneBlock:
		return inventory.NewItemStack(material.SugarCaneItem, 1, 0)
	case material.SignBlock, material.SignWall:
		return inventory.NewItemStack(material.SignItem, 1, 0)
	case material.SnowBlock, material.SnowLayer:
		return inventory.NewItemStack(material.Snowball, 1, 0)
	case material.Cobweb:
		return inventory.NewItemStack(material.StringItem, 1, 0)

	// blocks that drop themselves with specific metadata
	case material.Slab, material.DoubleSlab:
		return inventory.NewItemStack(material.Slab, 1, uint16(metadata))
	case material.Wood:
		return inventory.NewItemStack(material.Wood, 1, uint16(metadata))
	case material.Wool:
		return inventory.NewItemStack(material.Wool, 1, uint16(metadata))
	case material.Sapling:
		return inventory.NewItemStack(material.Sapling, 1, uint16(metadata&3))

	// drop the block itself with no metadata
	default:
		return inventory.NewItemStack(b, 1, 0)
	}
}
