package world

import (
	"fmt"

	"github.com/Pesekjak/173go/pkg/prot"
)

type World struct {
	dimension  Dimension
	SpawnPoint BlockPos
	time       int64
	height     byte

	generator Generator

	chunks map[ChunkPos]*Chunk

	entities map[int32]Entity
}

func NewWorld() (*World, error) {
	dimension := Overworld
	spawnPoint := NewBlockPos(0, 60, 0)
	time := int64(0)
	height := DefaultChunkHeight

	generator, err := MakeStandardFlatGenerator()
	if err != nil {
		return nil, err
	}

	chunks := make(map[ChunkPos]*Chunk)

	entities := make(map[int32]Entity)

	return &World{
		dimension:  dimension,
		SpawnPoint: spawnPoint,
		time:       time,
		height:     height,

		generator: generator,

		chunks: chunks,

		entities: entities,
	}, nil
}

func (w *World) Dimension() Dimension {
	return w.dimension
}

func (w *World) Time() int64 {
	return w.time
}

func (w *World) Height() byte {
	return w.height
}

func (w *World) SpawnPlayer(player PlayerEntity) error {
	if _, ok := w.entities[player.Id()]; ok {
		return fmt.Errorf("there is already an entity with id %v in this world", player.Id())
	}

	connection := player.Connection()

	for pos, chunk := range w.chunks {
		err := connection.WritePacket(&prot.PacketOutPreChunk{
			X:    pos.X,
			Z:    pos.Z,
			Load: true,
		}, false)
		if err != nil {
			return err
		}
		chunkData, err := chunk.data()
		if err != nil {
			return err
		}
		dataPacket := prot.PacketOutMapChunk{
			X:     pos.X * 16,
			Y:     0,
			Z:     pos.Z * 16,
			SizeX: byte(ChunkSize) - 1,
			SizeY: chunk.height - 1,
			SizeZ: byte(ChunkSize) - 1,
			Data:  chunkData,
		}
		err = connection.WritePacket(&dataPacket, false)
		if err != nil {
			return err
		}
	}

	return connection.Flush()
}

func (w *World) LoadChunk(pos ChunkPos) (*Chunk, error) {
	if loaded, ok := w.chunks[pos]; ok {
		return loaded, nil
	}

	chunk, err := newChunk(pos, w.height, func(block Block) error {
		return nil // no block updates for now
	})
	if err != nil {
		return nil, err
	}
	if err = w.generator.GenerateBlocks(chunk); err != nil {
		return nil, err
	}
	w.chunks[pos] = chunk
	return chunk, nil
}
