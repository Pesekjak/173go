package world

import (
	"bytes"
	"compress/zlib"
	"fmt"

	"github.com/Pesekjak/173go/pkg/net"
	"github.com/Pesekjak/173go/pkg/prot"
	"github.com/Pesekjak/173go/pkg/world/material"
)

type Block interface {
	Position() BlockPos
	Material() *material.Block
	Data() byte

	Set(block *material.Block, data byte) error
}

type Chunk struct {
	pos    ChunkPos
	height byte
	blocks map[uint16]Block

	updater func(blocks []Block) error

	cached []byte
}

func NewChunk(pos ChunkPos, height byte, updater func(blocks []Block) error) *Chunk {
	return &Chunk{
		pos:     pos,
		height:  height,
		blocks:  make(map[uint16]Block, 16*16*int(height)),
		updater: updater,
		cached:  nil,
	}
}

func (c *Chunk) GetBlock(x, y, z int32) (Block, error) {
	if x >= 16 || z >= 16 || y < 0 || y > int32(c.height) { // TODO height should be exclusive
		return nil, fmt.Errorf("coodinates %v;%v;%v are out of bounds of a chunk", x, y, z)
	}
	i := blockIndex(x, y, z)
	if b, ok := c.blocks[i]; ok {
		return b, nil
	}

	b := c.createBlock(x, y, z)
	c.blocks[i] = b
	return b, nil
}

func (c *Chunk) Load(conn *net.Connection) error {
	err := conn.WritePacket(&prot.PacketOutPreChunk{
		X:    c.pos.X,
		Z:    c.pos.Z,
		Load: true,
	}, false)
	if err != nil {
		return err
	}
	dataPacket := prot.PacketOutMapChunk{
		X:     c.pos.X * 16,
		Y:     0,
		Z:     c.pos.Z * 16,
		SizeX: 15,
		SizeY: 127,
		SizeZ: 15,
		Data:  c.cached,
	}

	if c.cached == nil {
		fmt.Println("creating chunk data")
		sizeX := 16
		sizeY := 128
		sizeZ := 16

		blockTypes := make([]byte, sizeX*sizeY*sizeZ)
		blockMetadata := make([]byte, (sizeX*sizeY*sizeZ)/2)
		blockLight := make([]byte, (sizeX*sizeY*sizeZ)/2)
		skyLight := make([]byte, (sizeX*sizeY*sizeZ)/2)

		for x := 0; x < sizeX; x++ {
			for z := 0; z < sizeZ; z++ {
				for y := 0; y < sizeY; y++ {
					index := y + (z * sizeY) + (x * sizeY * sizeZ)
					block, err := c.GetBlock(int32(x), int32(y), int32(z))
					if block == nil {
						fmt.Println("wtf:", err)
						continue
					}

					blockTypes[index] = byte(block.Material().Id())

					if index%2 == 0 {
						blockMetadata[index/2] = block.Data() & 0x0F
						blockLight[index/2] = 0x00
						skyLight[index/2] = 0x0F
					} else {
						blockMetadata[index/2] |= (block.Data() & 0x0F) << 4
						blockLight[index/2] |= 0x00 << 4
						skyLight[index/2] |= 0x0F << 4
					}
				}
			}
		}

		uncompressedData := append(blockTypes, blockMetadata...)
		uncompressedData = append(uncompressedData, blockLight...)
		uncompressedData = append(uncompressedData, skyLight...)

		var compressedData bytes.Buffer
		writer := zlib.NewWriter(&compressedData)
		writer.Write(uncompressedData)
		writer.Close()

		c.cached = compressedData.Bytes()
		dataPacket.Data = c.cached
	}

	return conn.WritePacket(&dataPacket, true)
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
	c.owner.cached = nil // invalidate cached chunk data
	return c.owner.updater([]Block{c})
}

func (c *Chunk) createBlock(x, y, z int32) Block {
	return &chunkBlock{
		owner:    c,
		pos:      NewBlockPos(c.pos.X*16+x, y, c.pos.Z*16+z),
		material: material.Air,
		data:     0,
	}
}

func blockIndex(x, y, z int32) uint16 {
	return uint16(y) | (uint16(z) << 8) | (uint16(x) << 12)
}
