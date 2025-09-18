package world

import (
	"bytes"
	"compress/zlib"
	"fmt"

	"github.com/Pesekjak/173go/pkg/world/material"
)

const (
	ChunkHeight uint32 = 128
	ChunkSize   uint32 = 16
)

type Chunk struct {
	pos ChunkPos

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

func newChunk(pos ChunkPos, updater func(block Block) error) (*Chunk, error) {
	bCount := ChunkSize * ChunkHeight * ChunkSize
	c := &Chunk{
		pos: pos,

		valid: true,

		blocks: make([]Block, bCount),

		blockTypes:    make([]byte, bCount),
		blockMetadata: make([]byte, bCount/2),
		blockLight:    newLight(),
		skyLight:      newLight(),

		cache:     nil,
		updater:   updater,
		generated: false,
	}

	for i := uint32(0); i < bCount; i++ {
		x, y, z := blockPos(i)
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
		block, err := c.GetBlock(x, y, z)
		if err != nil {
			return nil, err
		}
		err = block.Set(mat.(*material.Block), data)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Chunk) Pos() ChunkPos {
	return c.pos
}

func (c *Chunk) GetBlock(x, y, z uint32) (Block, error) {
	if _, err := inChunkBounds(x, y, z); err != nil {
		return nil, err
	}
	i := blockIndex(x, y, z)
	if block := c.blocks[i]; block != nil {
		return block, nil
	} else {
		block = &chunkBlock{
			owner:    c,
			pos:      NewBlockPos(c.Pos().X*16+int32(x), int32(y), c.Pos().Z*16+int32(z)),
			material: material.Air,
			data:     0,
		}
		c.blocks[i] = block
		return block, nil
	}
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

func (b *chunkBlock) Position() BlockPos {
	return b.pos
}

func (b *chunkBlock) Material() *material.Block {
	return b.material
}

func (b *chunkBlock) Data() byte {
	return b.data
}

func (b *chunkBlock) Set(block *material.Block, data byte) error {
	b.material = block
	b.data = data
	b.owner.cache = nil // invalidate cached chunk data

	chunkX := uint32(b.pos.X) % ChunkSize
	chunkY := uint32(b.pos.Y)
	chunkZ := uint32(b.pos.Z) % ChunkSize
	index := blockIndex(chunkX, chunkY, chunkZ)

	b.owner.blockTypes[index] = byte(block.Id())

	metaIndex := index / 2
	if index%2 == 0 {
		b.owner.blockMetadata[metaIndex] = (b.owner.blockMetadata[metaIndex] & 0xF0) | (data & 0x0F)
	} else {
		b.owner.blockMetadata[metaIndex] = (b.owner.blockMetadata[metaIndex] & 0x0F) | ((data & 0x0F) << 4)
	}

	if emission := block.LightEmission(); emission != 0 {
		if err := b.owner.blockLight.propagate(b.owner, chunkX, chunkY, chunkZ, emission); err != nil {
			return err
		}
	}

	if !b.owner.generated {
		return nil // chunk is not yet generated, no need to send block updates
	}
	return b.owner.updater(b)
}

func blockIndex(x, y, z uint32) uint32 {
	return y + (z * ChunkHeight) + (x * ChunkHeight * ChunkSize)
}

func blockPos(i uint32) (x, y, z uint32) {
	layerSize := ChunkHeight * ChunkSize
	x = i / layerSize
	remainder := i % layerSize
	z = remainder / ChunkHeight
	y = remainder % ChunkHeight
	return x, y, z
}

func inChunkBounds(x, y, z uint32) (bool, error) {
	if x >= ChunkSize || z >= ChunkSize || y >= ChunkHeight {
		return false, fmt.Errorf("coodinates %v;%v;%v are out of bounds of a chunk", x, y, z)
	}
	return true, nil
}
