package world

import (
	"bytes"
	"compress/zlib"
	"fmt"

	"github.com/Pesekjak/173go/pkg/world/material"
)

const ChunkSize uint32 = 16
const DefaultChunkHeight byte = 128

type Chunk struct {
	pos    ChunkPos
	height byte

	valid bool

	blocks []Block

	blockTypes    []byte
	blockMetadata []byte
	blockLight    *light
	skyLight      *light

	cache     []byte
	updater   func(block Block) error
	generated bool
}

func newChunk(pos ChunkPos, height byte, updater func(block Block) error) (*Chunk, error) {
	bCount := ChunkSize * uint32(height) * ChunkSize
	c := &Chunk{
		pos:    pos,
		height: height,

		valid: true,

		blocks: make([]Block, bCount),

		blockTypes:    make([]byte, bCount),
		blockMetadata: make([]byte, bCount/2),
		blockLight:    newLight(height),
		skyLight:      newLight(height),

		cache:     nil,
		updater:   updater,
		generated: false,
	}

	for i := uint32(0); i < bCount; i++ {
		x, y, z := blockPos(i, height)
		bPos := NewBlockPos(c.pos.X*int32(ChunkSize)+int32(x), int32(y), c.pos.Z*int32(ChunkSize)+int32(z))
		mat, err := material.FromID(uint16(c.blockTypes[i]))
		if err != nil || mat.(*material.Block) == nil {
			return nil, fmt.Errorf("failed to create block from ID %v", c.blockTypes[i])
		}
		data := c.blockMetadata[i/2]
		if i%2 == 0 {
			data = data & 0x0F
		} else {
			data = (data & 0xF0) >> 4
		}
		b := &chunkBlock{
			owner:    c,
			pos:      bPos,
			material: mat.(*material.Block),
			data:     data,
		}
		c.blocks[i] = b
	}

	return c, nil
}

func (c *Chunk) GetBlock(x, y, z uint32) (Block, error) {
	if x >= ChunkSize || z >= ChunkSize || y >= uint32(c.height) {
		return nil, fmt.Errorf("coodinates %v;%v;%v are out of bounds of a chunk", x, y, z)
	}
	i := blockIndex(x, y, z, c.height)
	return c.blocks[i], nil
}

func (c *Chunk) Height() byte {
	return c.height
}

func (c *Chunk) data() ([]byte, error) {
	if c.cache != nil {
		return c.cache, nil
	}

	uncompressedData := append(c.blockTypes, c.blockMetadata...)
	uncompressedData = append(uncompressedData, c.blockLight.data...)
	uncompressedData = append(uncompressedData, c.skyLight.data...)

	var compressedData bytes.Buffer
	writer := zlib.NewWriter(&compressedData)
	_, err := writer.Write(uncompressedData)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	c.cache = compressedData.Bytes()
	return c.cache, nil
}

type chunkBlock struct {
	owner    *Chunk
	pos      BlockPos
	material *material.Block
	data     byte
}

func (c *chunkBlock) Position() BlockPos {
	return c.pos
}

func (c *chunkBlock) Material() *material.Block {
	return c.material
}

func (c *chunkBlock) Data() byte {
	return c.data
}

func (c *chunkBlock) Set(block *material.Block, data byte) error {
	c.material = block
	c.data = data
	c.owner.cache = nil // invalidate cached chunk data

	chunkX := uint32(c.pos.X) % ChunkSize
	chunkY := uint32(c.pos.Y)
	chunkZ := uint32(c.pos.Z) % ChunkSize
	index := blockIndex(chunkX, chunkY, chunkZ, c.owner.height)

	c.owner.blockTypes[index] = byte(block.Id())

	metaIndex := index / 2
	if index%2 == 0 {
		c.owner.blockMetadata[metaIndex] = (c.owner.blockMetadata[metaIndex] & 0xF0) | (data & 0x0F)
	} else {
		c.owner.blockMetadata[metaIndex] = (c.owner.blockMetadata[metaIndex] & 0x0F) | ((data & 0x0F) << 4)
	}

	if !c.owner.generated {
		return nil // chunk is not yet generated, no need to send block updates
	}
	return c.owner.updater(c)
}

func blockIndex(x, y, z uint32, height byte) uint32 {
	height32 := uint32(height)
	return y + (z * height32) + (x * height32 * ChunkSize)
}

func blockPos(i uint32, height byte) (x, y, z uint32) {
	height32 := uint32(height)
	layerSize := height32 * ChunkSize
	x = i / layerSize
	remainder := i % layerSize
	z = remainder / height32
	y = remainder % height32
	return x, y, z
}
