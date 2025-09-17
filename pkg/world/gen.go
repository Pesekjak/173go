package world

import (
	"fmt"

	"github.com/Pesekjak/173go/pkg/world/material"
)

type Generator interface {
	GenerateBlocks(chunk *Chunk) error
}

type FlatGenerator struct {
	Layers []struct {
		Material *material.Block
		Data     byte
		Height   uint32
	}
}

func MakeStandardFlatGenerator() (*FlatGenerator, error) {
	return MakeFlatGenerator(
		material.Bedrock, byte(0), uint32(1),
		material.Dirt, byte(0), uint32(3),
		material.GrassBlock, byte(0), uint32(1),
	)
}

func MakeFlatGenerator(params ...interface{}) (*FlatGenerator, error) {
	if len(params)%3 != 0 {
		return nil, fmt.Errorf("invalid parameters for flat gen: %v", params)
	}
	numLayers := len(params) / 3
	layers := make([]struct {
		Material *material.Block
		Data     byte
		Height   uint32
	}, numLayers)
	for i := 0; i < numLayers; i++ {
		paramIndex := i * 3

		mat, ok := params[paramIndex].(*material.Block)
		if !ok {
			return nil, fmt.Errorf("expected *material.Block at index %d, got %T", paramIndex, params[paramIndex])
		}

		data, ok := params[paramIndex+1].(byte)
		if !ok {
			return nil, fmt.Errorf("expected byte at index %d, got %T", paramIndex+1, params[paramIndex+1])
		}

		height, ok := params[paramIndex+2].(uint32)
		if !ok {
			return nil, fmt.Errorf("expected uint32 at index %d, got %T", paramIndex+2, params[paramIndex+2])
		}

		layers[i].Material = mat
		layers[i].Data = data
		layers[i].Height = height
	}

	return &FlatGenerator{Layers: layers}, nil
}

func (g *FlatGenerator) GenerateBlocks(chunk *Chunk) error {
	y := uint32(0)
	cHeight := uint32(chunk.Height())

	for _, layer := range g.Layers {
		for level := uint32(0); level < layer.Height; level++ {
			if y == cHeight {
				return nil // the flat generator reached max world height
			}

			for x := uint32(0); x < ChunkSize; x++ {
				for z := uint32(0); z < ChunkSize; z++ {
					block, err := chunk.GetBlock(x, y, z)
					if err != nil {
						return err
					}
					err = block.Set(layer.Material, layer.Data)
					if err != nil {
						return err
					}
				}
			}

			y++
		}
	}
	return nil
}
