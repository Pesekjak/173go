package material

import "fmt"

type Material interface {
	Id() uint16
	String() string
}

type Block struct {
	id   uint16
	name string
	Group
	isCube              bool
	isFluidProof        bool
	lightOpacity        uint8
	lightEmission       uint8
	slipperiness        float32
	breakHardness       float32
	explosionResistance float32
	pistonPolicy        PistonPolicy
	fireFlammability    uint16
	fireBurnRate        uint16
}

func (b *Block) Id() uint16 {
	return b.id
}

func (b *Block) String() string {
	return b.name
}

type Item struct {
	id            uint16
	name          string
	maxStackSize  uint16
	maxDamage     uint16
	hasDurability bool
	isTool        bool
	isFood        bool
	equipSlot     EquipSlot
}

func (i *Item) Id() uint16 {
	return i.id
}

func (i *Item) String() string {
	return i.name
}

var materials = make(map[uint16]Material)

func newBlock(id byte, name string) *Block {
	iid := uint16(id)
	if _, exists := materials[iid]; exists {
		panic(fmt.Sprintf("block with ID %v is already registered", id))
	}
	b := &Block{
		id:                  iid,
		name:                name,
		Group:               GroupRock, // a common default
		isCube:              true,
		isFluidProof:        true,
		lightOpacity:        255, // fully opaque
		lightEmission:       0,
		slipperiness:        0.6,
		breakHardness:       0.0,
		explosionResistance: 0.0,
		pistonPolicy:        PistonPolicyPushPull,
		fireFlammability:    0,
		fireBurnRate:        0,
	}
	materials[iid] = b
	return b
}

func newItem(id uint16, name string) *Item {
	finalID := id + 256 // items offset
	if _, exists := materials[finalID]; exists {
		panic(fmt.Sprintf("item with final ID %v (base ID %v) is already registered", finalID, id))
	}
	i := &Item{
		id:            finalID,
		name:          name,
		maxStackSize:  64,
		maxDamage:     0,
		hasDurability: false,
		isTool:        false,
		isFood:        false,
		equipSlot:     SlotHand,
	}
	materials[finalID] = i
	return i
}

func FromID(id uint16) (Material, error) {
	m, ok := materials[id]
	if !ok {
		return nil, fmt.Errorf("material with id %d not found", id)
	}
	return m, nil
}

var (
	Air                 *Block
	Stone               *Block
	GrassBlock          *Block
	Dirt                *Block
	Cobblestone         *Block
	WoodenPlanks        *Block
	Sapling             *Block
	Bedrock             *Block
	WaterFlowing        *Block
	WaterStill          *Block
	LavaFlowing         *Block
	LavaStill           *Block
	Sand                *Block
	Gravel              *Block
	GoldOre             *Block
	IronOre             *Block
	CoalOre             *Block
	Wood                *Block
	Leaves              *Block
	Sponge              *Block
	Glass               *Block
	LapisLazuliOre      *Block
	LapisLazuliBlock    *Block
	Dispenser           *Block
	Sandstone           *Block
	NoteBlock           *Block
	BedBlock            *Block
	PoweredRail         *Block
	DetectorRail        *Block
	StickyPiston        *Block
	Cobweb              *Block
	TallGrass           *Block
	DeadBush            *Block
	Piston              *Block
	PistonHead          *Block
	Wool                *Block
	Dandelion           *Block
	Rose                *Block
	BrownMushroom       *Block
	RedMushroom         *Block
	GoldBlock           *Block
	IronBlock           *Block
	DoubleSlab          *Block
	Slab                *Block
	Bricks              *Block
	TNT                 *Block
	Bookshelf           *Block
	MossStone           *Block
	Obsidian            *Block
	Torch               *Block
	Fire                *Block
	MobSpawner          *Block
	WoodenStairs        *Block
	Chest               *Block
	RedstoneWire        *Block
	DiamondOre          *Block
	DiamondBlock        *Block
	CraftingTable       *Block
	Seeds               *Block
	Farmland            *Block
	Furnace             *Block
	FurnaceLit          *Block
	SignBlock           *Block
	WoodenDoorBlock     *Block
	Ladder              *Block
	Rails               *Block
	CobblestoneStairs   *Block
	SignWall            *Block
	Lever               *Block
	StonePressurePlate  *Block
	IronDoorBlock       *Block
	WoodenPressurePlate *Block
	RedstoneOre         *Block
	RedstoneOreGlowing  *Block
	RedstoneTorchOff    *Block
	RedstoneTorchOn     *Block
	StoneButton         *Block
	SnowLayer           *Block
	Ice                 *Block
	SnowBlock           *Block
	Cactus              *Block
	ClayBlock           *Block
	SugarCaneBlock      *Block
	Jukebox             *Block
	Fence               *Block
	Pumpkin             *Block
	Netherrack          *Block
	SoulSand            *Block
	GlowstoneBlock      *Block
	Portal              *Block
	JackOLantern        *Block
	CakeBlock           *Block
	RedstoneRepeaterOff *Block
	RedstoneRepeaterOn  *Block
	LockedChest         *Block
	Trapdoor            *Block

	IronShovel           *Item
	IronPickaxe          *Item
	IronAxe              *Item
	FlintAndSteel        *Item
	Apple                *Item
	Bow                  *Item
	Arrow                *Item
	Coal                 *Item
	Diamond              *Item
	IronIngot            *Item
	GoldIngot            *Item
	IronSword            *Item
	WoodenSword          *Item
	WoodenShovel         *Item
	WoodenPickaxe        *Item
	WoodenAxe            *Item
	StoneSword           *Item
	StoneShovel          *Item
	StonePickaxe         *Item
	StoneAxe             *Item
	DiamondSword         *Item
	DiamondShovel        *Item
	DiamondPickaxe       *Item
	DiamondAxe           *Item
	Stick                *Item
	Bowl                 *Item
	MushroomStew         *Item
	GoldSword            *Item
	GoldShovel           *Item
	GoldPickaxe          *Item
	GoldAxe              *Item
	StringItem           *Item
	Feather              *Item
	Gunpowder            *Item
	WoodenHoe            *Item
	StoneHoe             *Item
	IronHoe              *Item
	DiamondHoe           *Item
	GoldHoe              *Item
	WheatSeeds           *Item
	Wheat                *Item
	Bread                *Item
	LeatherHelmet        *Item
	LeatherTunic         *Item
	LeatherPants         *Item
	LeatherBoots         *Item
	ChainmailHelmet      *Item
	ChainmailChestplate  *Item
	ChainmailLeggings    *Item
	ChainmailBoots       *Item
	IronHelmet           *Item
	IronChestplate       *Item
	IronLeggings         *Item
	IronBoots            *Item
	DiamondHelmet        *Item
	DiamondChestplate    *Item
	DiamondLeggings      *Item
	DiamondBoots         *Item
	GoldHelmet           *Item
	GoldChestplate       *Item
	GoldLeggings         *Item
	GoldBoots            *Item
	Flint                *Item
	RawPorkchop          *Item
	CookedPorkchop       *Item
	Painting             *Item
	GoldenApple          *Item
	SignItem             *Item
	WoodenDoorItem       *Item
	Bucket               *Item
	WaterBucket          *Item
	LavaBucket           *Item
	Minecart             *Item
	Saddle               *Item
	IronDoorItem         *Item
	RedstoneDust         *Item
	Snowball             *Item
	Boat                 *Item
	Leather              *Item
	MilkBucket           *Item
	ClayBrick            *Item
	ClayBalls            *Item
	SugarCaneItem        *Item
	Paper                *Item
	Book                 *Item
	Slimeball            *Item
	StorageMinecart      *Item
	PoweredMinecart      *Item
	Egg                  *Item
	Compass              *Item
	FishingRod           *Item
	Clock                *Item
	GlowstoneDust        *Item
	RawFish              *Item
	CookedFish           *Item
	Dye                  *Item
	Bone                 *Item
	Sugar                *Item
	CakeItem             *Item
	BedItem              *Item
	RedstoneRepeaterItem *Item
	Cookie               *Item
	Map                  *Item
	Shears               *Item
	Record13             *Item
	RecordCat            *Item
)

func init() {
	Air = newBlock(0, "Air").withGroup(GroupAir).nonCube().withHardness(0)
	Stone = newBlock(1, "Stone").withGroup(GroupRock).withHardness(1.5).withExplosionResistance(30.0 / 5.0)
	GrassBlock = newBlock(2, "Grass Block").withGroup(GroupGrass).withHardness(0.6)
	Dirt = newBlock(3, "Dirt").withGroup(GroupGround).withHardness(0.5)
	Cobblestone = newBlock(4, "Cobblestone").withGroup(GroupRock).withHardness(2.0).withExplosionResistance(30.0 / 5.0)
	WoodenPlanks = newBlock(5, "Wooden Planks").withGroup(GroupWood).withHardness(2.0).withExplosionResistance(15.0/5.0).withFireProperties(5, 20)
	Sapling = newBlock(6, "Sapling").withGroup(GroupPlant).nonCube().withHardness(0)
	Bedrock = newBlock(7, "Bedrock").withGroup(GroupRock).unbreakable().withPistonPolicy(PistonPolicyStop)
	WaterFlowing = newBlock(8, "Flowing Water").withGroup(GroupWater).nonCube().withLightOpacity(3).unbreakable()
	WaterStill = newBlock(9, "Still Water").withGroup(GroupWater).nonCube().withLightOpacity(3).unbreakable()
	LavaFlowing = newBlock(10, "Flowing Lava").withGroup(GroupLava).nonCube().withLightEmission(15).unbreakable()
	LavaStill = newBlock(11, "Still Lava").withGroup(GroupLava).nonCube().withLightEmission(15).unbreakable()
	Sand = newBlock(12, "Sand").withGroup(GroupSand).withHardness(0.5)
	Gravel = newBlock(13, "Gravel").withGroup(GroupSand).withHardness(0.6)
	GoldOre = newBlock(14, "Gold Ore").withGroup(GroupRock).withHardness(3.0).withExplosionResistance(15.0 / 5.0)
	IronOre = newBlock(15, "Iron Ore").withGroup(GroupRock).withHardness(3.0).withExplosionResistance(15.0 / 5.0)
	CoalOre = newBlock(16, "Coal Ore").withGroup(GroupRock).withHardness(3.0).withExplosionResistance(15.0 / 5.0)
	Wood = newBlock(17, "Wood").withGroup(GroupWood).withHardness(2.0).withExplosionResistance(15.0/5.0).withFireProperties(5, 5)
	Leaves = newBlock(18, "Leaves").withGroup(GroupLeaves).withLightOpacity(1).withHardness(0.2).withFireProperties(30, 60)
	Sponge = newBlock(19, "Sponge").withGroup(GroupSponge).withHardness(0.6)
	Glass = newBlock(20, "Glass").withGroup(GroupGlass).withHardness(0.3)
	LapisLazuliOre = newBlock(21, "Lapis Lazuli Ore").withGroup(GroupRock).withHardness(3.0).withExplosionResistance(15.0 / 5.0)
	LapisLazuliBlock = newBlock(22, "Lapis Lazuli Block").withGroup(GroupRock).withHardness(3.0).withExplosionResistance(15.0 / 5.0)
	Dispenser = newBlock(23, "Dispenser").withGroup(GroupRock).withHardness(3.5)
	Sandstone = newBlock(24, "Sandstone").withGroup(GroupRock).withHardness(0.8)
	NoteBlock = newBlock(25, "Note Block").withGroup(GroupWood).withHardness(0.8)
	BedBlock = newBlock(26, "Bed Block").withGroup(GroupCloth).nonCube().withHardness(0.2).withPistonPolicy(PistonPolicyBreak)
	PoweredRail = newBlock(27, "Powered Rail").withGroup(GroupCircuit).nonCube().withHardness(0.7)
	DetectorRail = newBlock(28, "Detector Rail").withGroup(GroupCircuit).nonCube().withHardness(0.7)
	StickyPiston = newBlock(29, "Sticky Piston").withGroup(GroupPiston).withHardness(0.5)
	Cobweb = newBlock(30, "Cobweb").withGroup(GroupCobweb).nonCube().withLightOpacity(1).withHardness(4.0)
	TallGrass = newBlock(31, "Tall Grass").withGroup(GroupPlant).nonCube().withHardness(0).withFireProperties(60, 100)
	DeadBush = newBlock(32, "Dead Bush").withGroup(GroupPlant).nonCube().withHardness(0)
	Piston = newBlock(33, "Piston").withGroup(GroupPiston).nonCube().withHardness(0.5)
	PistonHead = newBlock(34, "Piston Head").withGroup(GroupPiston).nonCube().withHardness(0.5).withPistonPolicy(PistonPolicyStop)
	Wool = newBlock(35, "Wool").withGroup(GroupCloth).withHardness(0.8).withFireProperties(30, 60)
	Dandelion = newBlock(37, "Dandelion").withGroup(GroupPlant).nonCube().withHardness(0)
	Rose = newBlock(38, "Rose").withGroup(GroupPlant).nonCube().withHardness(0)
	BrownMushroom = newBlock(39, "Brown Mushroom").withGroup(GroupPlant).nonCube().withLightEmission(1).withHardness(0)
	RedMushroom = newBlock(40, "Red Mushroom").withGroup(GroupPlant).nonCube().withHardness(0)
	GoldBlock = newBlock(41, "Gold Block").withGroup(GroupIron).withHardness(3.0).withExplosionResistance(30.0 / 5.0)
	IronBlock = newBlock(42, "Iron Block").withGroup(GroupIron).withHardness(5.0).withExplosionResistance(30.0 / 5.0)
	DoubleSlab = newBlock(43, "Double Slab").withGroup(GroupRock).withHardness(2.0).withExplosionResistance(30.0 / 5.0)
	Slab = newBlock(44, "Slab").withGroup(GroupRock).nonCube().withHardness(2.0).withExplosionResistance(30.0 / 5.0)
	Bricks = newBlock(45, "Bricks").withGroup(GroupRock).withHardness(2.0).withExplosionResistance(30.0 / 5.0)
	TNT = newBlock(46, "TNT").withGroup(GroupTNT).withHardness(0).withFireProperties(15, 100)
	Bookshelf = newBlock(47, "Bookshelf").withGroup(GroupWood).withHardness(1.5).withFireProperties(30, 20)
	MossStone = newBlock(48, "Moss Stone").withGroup(GroupRock).withHardness(2.0).withExplosionResistance(30.0 / 5.0)
	Obsidian = newBlock(49, "Obsidian").withGroup(GroupRock).withHardness(10.0).withExplosionResistance(6000.0 / 5.0).withPistonPolicy(PistonPolicyStop)
	Torch = newBlock(50, "Torch").withGroup(GroupCircuit).nonCube().withLightEmission(14).withHardness(0)
	Fire = newBlock(51, "Fire").withGroup(GroupFire).nonCube().withLightEmission(15).withHardness(0)
	MobSpawner = newBlock(52, "Mob Spawner").withGroup(GroupRock).withHardness(5.0)
	WoodenStairs = newBlock(53, "Wooden Stairs").withGroup(GroupWood).nonCube().withHardness(2.0).withExplosionResistance(15.0/5.0).withFireProperties(5, 20)
	Chest = newBlock(54, "Chest").withGroup(GroupWood).withHardness(2.5)
	RedstoneWire = newBlock(55, "Redstone Wire").withGroup(GroupCircuit).nonCube().withHardness(0)
	DiamondOre = newBlock(56, "Diamond Ore").withGroup(GroupRock).withHardness(3.0).withExplosionResistance(15.0 / 5.0)
	DiamondBlock = newBlock(57, "Diamond Block").withGroup(GroupIron).withHardness(5.0).withExplosionResistance(30.0 / 5.0)
	CraftingTable = newBlock(58, "Crafting Table").withGroup(GroupWood).withHardness(2.5)
	Seeds = newBlock(59, "Seeds").withGroup(GroupPlant).nonCube().withHardness(0)
	Farmland = newBlock(60, "Farmland").withGroup(GroupGround).nonCube().withHardness(0.6)
	Furnace = newBlock(61, "Furnace").withGroup(GroupRock).withHardness(3.5)
	FurnaceLit = newBlock(62, "Lit Furnace").withGroup(GroupRock).withLightEmission(13).withHardness(3.5)
	SignBlock = newBlock(63, "Sign Block").withGroup(GroupWood).nonCube().withHardness(1.0).makeFluidProof()
	WoodenDoorBlock = newBlock(64, "Wooden Door Block").withGroup(GroupWood).nonCube().withHardness(3.0).withPistonPolicy(PistonPolicyBreak).makeFluidProof()
	Ladder = newBlock(65, "Ladder").withGroup(GroupCircuit).nonCube().withHardness(0.4).makeFluidProof()
	Rails = newBlock(66, "Rails").withGroup(GroupCircuit).nonCube().withHardness(0.7)
	CobblestoneStairs = newBlock(67, "Cobblestone Stairs").withGroup(GroupRock).nonCube().withHardness(2.0).withExplosionResistance(30.0 / 5.0)
	SignWall = newBlock(68, "Wall Sign").withGroup(GroupWood).nonCube().withHardness(1.0).makeFluidProof()
	Lever = newBlock(69, "Lever").withGroup(GroupCircuit).nonCube().withHardness(0.5)
	StonePressurePlate = newBlock(70, "Stone Pressure Plate").withGroup(GroupRock).nonCube().withHardness(0.5).withPistonPolicy(PistonPolicyBreak)
	IronDoorBlock = newBlock(71, "Iron Door Block").withGroup(GroupIron).nonCube().withHardness(5.0).withExplosionResistance(15.0 / 5.0).withPistonPolicy(PistonPolicyBreak).makeFluidProof()
	WoodenPressurePlate = newBlock(72, "Wooden Pressure Plate").withGroup(GroupWood).nonCube().withHardness(0.5).withPistonPolicy(PistonPolicyBreak)
	RedstoneOre = newBlock(73, "Redstone Ore").withGroup(GroupRock).withHardness(3.0).withExplosionResistance(15.0 / 5.0)
	RedstoneOreGlowing = newBlock(74, "Glowing Redstone Ore").withGroup(GroupRock).withLightEmission(9).withHardness(3.0).withExplosionResistance(15.0 / 5.0)
	RedstoneTorchOff = newBlock(75, "Redstone Torch (Off)").withGroup(GroupCircuit).nonCube().withHardness(0)
	RedstoneTorchOn = newBlock(76, "Redstone Torch (On)").withGroup(GroupCircuit).nonCube().withLightEmission(7).withHardness(0)
	StoneButton = newBlock(77, "Stone Button").withGroup(GroupRock).nonCube().withHardness(0.5)
	SnowLayer = newBlock(78, "Snow Layer").withGroup(GroupSnow).nonCube().withHardness(0.1)
	Ice = newBlock(79, "Ice").withGroup(GroupIce).withLightOpacity(3).withHardness(0.5).withSlipperiness(0.95)
	SnowBlock = newBlock(80, "Snow Block").withGroup(GroupSnowBlock).withHardness(0.2)
	Cactus = newBlock(81, "Cactus").withGroup(GroupCactus).nonCube().withHardness(0.4)
	ClayBlock = newBlock(82, "Clay Block").withGroup(GroupClay).withHardness(0.6)
	SugarCaneBlock = newBlock(83, "Sugar Cane Block").withGroup(GroupPlant).nonCube().withHardness(0).makeFluidProof()
	Jukebox = newBlock(84, "Jukebox").withGroup(GroupWood).withHardness(2.0).withExplosionResistance(30.0 / 5.0)
	Fence = newBlock(85, "Fence").withGroup(GroupWood).nonCube().withHardness(2.0).withExplosionResistance(15.0/5.0).withFireProperties(5, 20)
	Pumpkin = newBlock(86, "Pumpkin").withGroup(GroupPumpkin).withHardness(1.0)
	Netherrack = newBlock(87, "Netherrack").withGroup(GroupRock).withHardness(0.4)
	SoulSand = newBlock(88, "Soul Sand").withGroup(GroupSand).withHardness(0.5)
	GlowstoneBlock = newBlock(89, "Glowstone Block").withGroup(GroupRock).withLightEmission(15).withHardness(0.3)
	Portal = newBlock(90, "Portal").withGroup(GroupPortal).nonCube().withLightEmission(11).unbreakable().withPistonPolicy(PistonPolicyStop)
	JackOLantern = newBlock(91, "Jack 'o' Lantern").withGroup(GroupPumpkin).withLightEmission(15).withHardness(1.0)
	CakeBlock = newBlock(92, "Cake Block").withGroup(GroupCake).nonCube().withHardness(0.5)
	RedstoneRepeaterOff = newBlock(93, "Redstone Repeater (Off)").withGroup(GroupCircuit).nonCube().withHardness(0)
	RedstoneRepeaterOn = newBlock(94, "Redstone Repeater (On)").withGroup(GroupCircuit).nonCube().withLightEmission(9).withHardness(0)
	LockedChest = newBlock(95, "Locked Chest").withGroup(GroupWood).withHardness(2.5)
	Trapdoor = newBlock(96, "Trapdoor").withGroup(GroupWood).nonCube().withHardness(3.0)

	const (
		woodMaxUses    = 59
		stoneMaxUses   = 131
		ironMaxUses    = 250
		goldMaxUses    = 32
		diamondMaxUses = 1561
	)

	IronShovel = newItem(0, "Iron Shovel").makeTool(ironMaxUses)
	IronPickaxe = newItem(1, "Iron Pickaxe").makeTool(ironMaxUses)
	IronAxe = newItem(2, "Iron Axe").makeTool(ironMaxUses)
	FlintAndSteel = newItem(3, "Flint and Steel").makeTool(64)
	Apple = newItem(4, "Apple").makeFood()
	Bow = newItem(5, "Bow").makeTool(384) // bow durability is 384
	Arrow = newItem(6, "Arrow")
	Coal = newItem(7, "Coal")
	Diamond = newItem(8, "Diamond")
	IronIngot = newItem(9, "Iron Ingot")
	GoldIngot = newItem(10, "Gold Ingot")
	IronSword = newItem(11, "Iron Sword").makeTool(ironMaxUses)
	WoodenSword = newItem(12, "Wooden Sword").makeTool(woodMaxUses)
	WoodenShovel = newItem(13, "Wooden Shovel").makeTool(woodMaxUses)
	WoodenPickaxe = newItem(14, "Wooden Pickaxe").makeTool(woodMaxUses)
	WoodenAxe = newItem(15, "Wooden Axe").makeTool(woodMaxUses)
	StoneSword = newItem(16, "Stone Sword").makeTool(stoneMaxUses)
	StoneShovel = newItem(17, "Stone Shovel").makeTool(stoneMaxUses)
	StonePickaxe = newItem(18, "Stone Pickaxe").makeTool(stoneMaxUses)
	StoneAxe = newItem(19, "Stone Axe").makeTool(stoneMaxUses)
	DiamondSword = newItem(20, "Diamond Sword").makeTool(diamondMaxUses)
	DiamondShovel = newItem(21, "Diamond Shovel").makeTool(diamondMaxUses)
	DiamondPickaxe = newItem(22, "Diamond Pickaxe").makeTool(diamondMaxUses)
	DiamondAxe = newItem(23, "Diamond Axe").makeTool(diamondMaxUses)
	Stick = newItem(24, "Stick")
	Bowl = newItem(25, "Bowl")
	MushroomStew = newItem(26, "Mushroom Stew").makeFood()
	GoldSword = newItem(27, "Gold Sword").makeTool(goldMaxUses)
	GoldShovel = newItem(28, "Gold Shovel").makeTool(goldMaxUses)
	GoldPickaxe = newItem(29, "Gold Pickaxe").makeTool(goldMaxUses)
	GoldAxe = newItem(30, "Gold Axe").makeTool(goldMaxUses)
	StringItem = newItem(31, "String")
	Feather = newItem(32, "Feather")
	Gunpowder = newItem(33, "Gunpowder")
	WoodenHoe = newItem(34, "Wooden Hoe").makeTool(woodMaxUses)
	StoneHoe = newItem(35, "Stone Hoe").makeTool(stoneMaxUses)
	IronHoe = newItem(36, "Iron Hoe").makeTool(ironMaxUses)
	DiamondHoe = newItem(37, "Diamond Hoe").makeTool(diamondMaxUses)
	GoldHoe = newItem(38, "Gold Hoe").makeTool(goldMaxUses)
	WheatSeeds = newItem(39, "Wheat Seeds")
	Wheat = newItem(40, "Wheat")
	Bread = newItem(41, "Bread").makeFood()
	LeatherHelmet = newItem(42, "Leather Helmet").makeTool(11 * 3).withEquipSlot(SlotHead)
	LeatherTunic = newItem(43, "Leather Tunic").makeTool(16 * 3).withEquipSlot(SlotChest)
	LeatherPants = newItem(44, "Leather Pants").makeTool(15 * 3).withEquipSlot(SlotLegs)
	LeatherBoots = newItem(45, "Leather Boots").makeTool(13 * 3).withEquipSlot(SlotFeet)
	ChainmailHelmet = newItem(46, "Chainmail Helmet").makeTool(11 * 6).withEquipSlot(SlotHead)
	ChainmailChestplate = newItem(47, "Chainmail Chestplate").makeTool(16 * 6).withEquipSlot(SlotChest)
	ChainmailLeggings = newItem(48, "Chainmail Leggings").makeTool(15 * 6).withEquipSlot(SlotLegs)
	ChainmailBoots = newItem(49, "Chainmail Boots").makeTool(13 * 6).withEquipSlot(SlotFeet)
	IronHelmet = newItem(50, "Iron Helmet").makeTool(11 * 12).withEquipSlot(SlotHead)
	IronChestplate = newItem(51, "Iron Chestplate").makeTool(16 * 12).withEquipSlot(SlotChest)
	IronLeggings = newItem(52, "Iron Leggings").makeTool(15 * 12).withEquipSlot(SlotLegs)
	IronBoots = newItem(53, "Iron Boots").makeTool(13 * 12).withEquipSlot(SlotFeet)
	DiamondHelmet = newItem(54, "Diamond Helmet").makeTool(11 * 24).withEquipSlot(SlotHead)
	DiamondChestplate = newItem(55, "Diamond Chestplate").makeTool(16 * 24).withEquipSlot(SlotChest)
	DiamondLeggings = newItem(56, "Diamond Leggings").makeTool(15 * 24).withEquipSlot(SlotLegs)
	DiamondBoots = newItem(57, "Diamond Boots").makeTool(13 * 24).withEquipSlot(SlotFeet)
	GoldHelmet = newItem(58, "Gold Helmet").makeTool(11 * 6).withEquipSlot(SlotHead)
	GoldChestplate = newItem(59, "Gold Chestplate").makeTool(16 * 6).withEquipSlot(SlotChest)
	GoldLeggings = newItem(60, "Gold Leggings").makeTool(15 * 6).withEquipSlot(SlotLegs)
	GoldBoots = newItem(61, "Gold Boots").makeTool(13 * 6).withEquipSlot(SlotFeet)
	Flint = newItem(62, "Flint")
	RawPorkchop = newItem(63, "Raw Porkchop").makeFood()
	CookedPorkchop = newItem(64, "Cooked Porkchop").makeFood()
	Painting = newItem(65, "Painting")
	GoldenApple = newItem(66, "Golden Apple").makeFood()
	SignItem = newItem(67, "Sign").withMaxStackSize(1)
	WoodenDoorItem = newItem(68, "Wooden Door").withMaxStackSize(1)
	Bucket = newItem(69, "Bucket").withMaxStackSize(16) // Buckets stack to 16
	WaterBucket = newItem(70, "Water Bucket").withMaxStackSize(1)
	LavaBucket = newItem(71, "Lava Bucket").withMaxStackSize(1)
	Minecart = newItem(72, "Minecart").withMaxStackSize(1)
	Saddle = newItem(73, "Saddle").withMaxStackSize(1)
	IronDoorItem = newItem(74, "Iron Door").withMaxStackSize(1)
	RedstoneDust = newItem(75, "Redstone Dust")
	Snowball = newItem(76, "Snowball").withMaxStackSize(16)
	Boat = newItem(77, "Boat").withMaxStackSize(1)
	Leather = newItem(78, "Leather")
	MilkBucket = newItem(79, "Milk Bucket").makeFood()
	ClayBrick = newItem(80, "Clay Brick")
	ClayBalls = newItem(81, "Clay Balls")
	SugarCaneItem = newItem(82, "Sugar Cane")
	Paper = newItem(83, "Paper")
	Book = newItem(84, "Book")
	Slimeball = newItem(85, "Slimeball")
	StorageMinecart = newItem(86, "Storage Minecart").withMaxStackSize(1)
	PoweredMinecart = newItem(87, "Powered Minecart").withMaxStackSize(1)
	Egg = newItem(88, "Egg").withMaxStackSize(16)
	Compass = newItem(89, "Compass").withMaxStackSize(1)
	FishingRod = newItem(90, "Fishing Rod").makeTool(64)
	Clock = newItem(91, "Clock").withMaxStackSize(1)
	GlowstoneDust = newItem(92, "Glowstone Dust")
	RawFish = newItem(93, "Raw Fish").makeFood()
	CookedFish = newItem(94, "Cooked Fish").makeFood()
	Dye = newItem(95, "Dye")
	Bone = newItem(96, "Bone")
	Sugar = newItem(97, "Sugar")
	CakeItem = newItem(98, "Cake").withMaxStackSize(1)
	BedItem = newItem(99, "Bed").withMaxStackSize(1)
	RedstoneRepeaterItem = newItem(100, "Redstone Repeater")
	Cookie = newItem(101, "Cookie").withMaxStackSize(8).makeFood()
	Map = newItem(102, "Map").withMaxStackSize(1)
	Shears = newItem(103, "Shears").makeTool(238)
	Record13 = newItem(2000, "13 Disc").withMaxStackSize(1)
	RecordCat = newItem(2001, "Cat Disc").withMaxStackSize(1)
}
