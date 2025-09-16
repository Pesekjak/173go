package material

import "math"

// Group defines the general category of a block.
type Group int

const (
	GroupAir Group = iota
	GroupGrass
	GroupGround
	GroupWood
	GroupRock
	GroupIron
	GroupWater
	GroupLava
	GroupLeaves
	GroupPlant
	GroupSponge
	GroupCloth
	GroupFire
	GroupSand
	GroupCircuit
	GroupGlass
	GroupTNT
	GroupUnused // unused but kept for consistency
	GroupIce
	GroupSnow
	GroupSnowBlock
	GroupCactus
	GroupClay
	GroupPumpkin
	GroupPortal
	GroupCake
	GroupCobweb
	GroupPiston
)

// IsSolid returns true if the group represents a solid block.
func (g Group) IsSolid() bool {
	switch g {
	case GroupAir, GroupWater, GroupLava, GroupPlant, GroupSnow, GroupCircuit, GroupPortal, GroupFire:
		return false
	default:
		return true
	}
}

// IsFluid returns true if the group is a fluid (Water or Lava).
func (g Group) IsFluid() bool {
	return g == GroupWater || g == GroupLava
}

// IsTranslucent returns true for groups that are solid but not fully opaque.
func (g Group) IsTranslucent() bool {
	switch g {
	case GroupLeaves, GroupGlass, GroupTNT, GroupIce, GroupSnow, GroupCactus:
		return true
	default:
		return false
	}
}

// IsOpaque returns true if the group is solid and not translucent.
func (g Group) IsOpaque() bool {
	return !g.IsTranslucent() && g.IsSolid()
}

// IsReplaceable returns true if another block can easily replace this one.
func (g Group) IsReplaceable() bool {
	switch g {
	case GroupAir, GroupWater, GroupLava, GroupSnow, GroupFire:
		return true
	default:
		return false
	}
}

// IsBreakableByDefault returns false for groups that require specific tools.
func (g Group) IsBreakableByDefault() bool {
	switch g {
	case GroupRock, GroupIron, GroupSnow, GroupSnowBlock, GroupCobweb:
		return false
	default:
		return true
	}
}

// PistonPolicy defines how a block interacts with pistons.
type PistonPolicy int

const (
	// PistonPolicyBreak indicates the block will break when pushed.
	PistonPolicyBreak PistonPolicy = iota
	// PistonPolicyPushPull indicates the block can be pushed and pulled.
	PistonPolicyPushPull
	// PistonPolicyStop indicates the block cannot be moved by a piston.
	PistonPolicyStop
)

func (b *Block) withGroup(g Group) *Block {
	b.Group = g
	b.isFluidProof = g.IsSolid() // default fluid-proofing to solidity

	// update light opacity and piston policy based on group, can be overridden later
	if !g.IsOpaque() {
		b.lightOpacity = 0
	} else {
		b.lightOpacity = 255
	}
	//goland:noinspection GoDfaConstantCondition keep GroupAir for consistency
	switch g {
	case GroupAir, GroupWater, GroupLava, GroupLeaves, GroupPlant, GroupFire, GroupCircuit, GroupSnow, GroupCactus,
		GroupPumpkin, GroupCake, GroupCobweb:
		b.pistonPolicy = PistonPolicyBreak
	default:
		b.pistonPolicy = PistonPolicyPushPull
	}
	return b
}

func (b *Block) nonCube() *Block {
	b.isCube = false
	b.lightOpacity = 0 // most non-cubes don't fully block light
	return b
}

func (b *Block) makeFluidProof() *Block {
	b.isFluidProof = true
	return b
}

func (b *Block) withLightOpacity(opacity uint8) *Block {
	b.lightOpacity = opacity
	return b
}

func (b *Block) withLightEmission(emission uint8) *Block {
	b.lightEmission = emission
	return b
}

func (b *Block) withSlipperiness(slipperiness float32) *Block {
	b.slipperiness = slipperiness
	return b
}

func (b *Block) withHardness(hardness float32) *Block {
	b.breakHardness = hardness
	b.explosionResistance = hardness
	return b
}

func (b *Block) withExplosionResistance(resistance float32) *Block {
	b.explosionResistance = resistance
	return b
}

func (b *Block) unbreakable() *Block {
	b.breakHardness = float32(math.Inf(1))
	b.explosionResistance = 18000000.0 / 5.0 // bedrock resistance
	return b
}

func (b *Block) withPistonPolicy(policy PistonPolicy) *Block {
	b.pistonPolicy = policy
	return b
}

func (b *Block) withFireProperties(flammability, burnRate uint16) *Block {
	b.fireFlammability = flammability
	b.fireBurnRate = burnRate
	return b
}

// IsCube returns true if a block is a full 1x1x1 cube.
func (b *Block) IsCube() bool {
	return b.isCube
}

// IsOpaqueCube returns true if a block is a full, opaque cube that blocks vision.
func (b *Block) IsOpaqueCube() bool {
	return b.isCube && b.Group.IsOpaque()
}

// IsNormalCube returns true if the block is a full, opaque cube.
func (b *Block) IsNormalCube() bool {
	return b.Group.IsOpaque() && b.isCube
}

// IsFluidProof returns true if the block can block fluid flow.
func (b *Block) IsFluidProof() bool {
	return b.isFluidProof
}

// LightOpacity returns the amount of light a block absorbs (0-255).
func (b *Block) LightOpacity() uint8 {
	return b.lightOpacity
}

// LightEmission returns the amount of light a block emits (0-15).
func (b *Block) LightEmission() uint8 {
	return b.lightEmission
}

// Slipperiness returns the slipperiness value for entities.
func (b *Block) Slipperiness() float32 {
	return b.slipperiness
}

// BreakHardness returns the time it takes to break a block. Returns +inf for unbreakable blocks.
func (b *Block) BreakHardness() float32 {
	return b.breakHardness
}

// ExplosionResistance returns the block's resistance to explosions.
func (b *Block) ExplosionResistance() float32 {
	return b.explosionResistance
}

// FireFlammability returns the chance for fire to spread to this block.
func (b *Block) FireFlammability() uint16 {
	return b.fireFlammability
}

// FireBurnRate returns the chance for this block to be destroyed by fire.
func (b *Block) FireBurnRate() uint16 {
	return b.fireBurnRate
}

// PistonPolicy returns how a piston interacts with the block.
func (b *Block) PistonPolicy(metadata byte) PistonPolicy {
	switch b {
	// special case for pistons
	case Piston, StickyPiston:
		// a piston is extended if the 4th bit (0x8) of its metadata is set
		if (metadata & 0x8) != 0 {
			return PistonPolicyStop
		}
		return PistonPolicyPushPull
	default:
		return b.pistonPolicy
	}
}
