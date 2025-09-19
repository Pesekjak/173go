package world

import (
	"fmt"
)

// light represents encoded lighting data with a chunk
type light struct {
	// encoded chunk light data
	data []byte
}

// newLight provides new empty light instance for a chunk
func newLight() *light {
	return &light{data: make([]byte, (ChunkSize*ChunkHeight*ChunkSize)/2)}
}

// get gets a light value at given chunk coordinates
func (l *light) get(x, y, z uint32) (byte, error) {
	if _, err := inChunkBounds(x, y, z); err != nil {
		return 0, err
	}
	i := blockIndex(x, y, z)
	val := l.data[i/2]
	if i%2 == 0 {
		return val & 0x0F, nil
	} else {
		return val >> 4, nil
	}
}

// set updates a light value at given chunk coordinates
func (l *light) set(x, y, z uint32, value byte) error {
	if value > 15 {
		return fmt.Errorf("light value out of bounds: %v", value)
	}
	if _, err := inChunkBounds(x, y, z); err != nil {
		return err
	}
	i := blockIndex(x, y, z)
	val := l.data[i/2]
	if i%2 == 0 {
		l.data[i/2] = (val & 0xF0) | value
	} else {
		l.data[i/2] = (val & 0x0F) | (value << 4)
	}
	return nil
}

// lightPropagation is data used during the light propagation updates
type lightPropagation struct {
	x, y, z int32
	value   byte
}

// neighbours returns lightPropagation around itself.
// For each neighbour block within the chunk new lightPropagation instance is
// returned with its light value decreased by 1
func (lp *lightPropagation) neighbours() []lightPropagation {
	newValue := lp.value - 1
	neighbours := []lightPropagation{
		{lp.x - 1, lp.y, lp.z, newValue},
		{lp.x + 1, lp.y, lp.z, newValue},
		{lp.x, lp.y, lp.z - 1, newValue},
		{lp.x, lp.y, lp.z + 1, newValue},
	}
	if lp.y > 0 {
		neighbours = append(neighbours, lightPropagation{lp.x, lp.y - 1, lp.z, newValue})
	}
	if lp.y < int32(ChunkHeight)-1 {
		neighbours = append(neighbours, lightPropagation{lp.x, lp.y + 1, lp.z, newValue})
	}
	return neighbours
}

// lightChunkBorders populates the given chunk with lightning from its neighbours
func lightChunkBorders(world *World, chunk *Chunk) error {
	chunkPos := chunk.Pos()
	var incomingLight []lightPropagation

	// stage 1: gather all light sources from  neighbours
	for _, neighbourPos := range chunkPos.Neighbours() {
		neighbour, ok := world.Chunk(neighbourPos)
		if !ok {
			continue // neighbour is not loaded
		}

		dx := neighbourPos.X - chunkPos.X
		dz := neighbourPos.Z - chunkPos.Z

		var neighbourX, neighbourZ uint32 // relative neighbour chunk coordinates to access the light to share
		switch {
		case dx == -1:
			neighbourX = ChunkSize - 1
		case dx == 1:
			neighbourX = 0
		case dz == -1:
			neighbourZ = ChunkSize - 1
		case dz == 1:
			neighbourZ = 0
		}

		// loop along the shared border or edge
		for y := uint32(0); y < ChunkHeight; y++ {
			if dx != 0 { // x face
				for z := uint32(0); z < ChunkSize; z++ {
					val, err := neighbour.blockLight.get(neighbourX, y, z)
					if err != nil {
						return err
					}
					if val > 1 {
						lp := lightPropagation{neighbourPos.X*16 + int32(neighbourX), int32(y), neighbourPos.Z*16 + int32(z), val}
						incomingLight = append(incomingLight, lp)
					}
				}
			} else { // z face
				for x := uint32(0); x < ChunkSize; x++ {
					val, err := neighbour.blockLight.get(x, y, neighbourZ)
					if err != nil {
						return err
					}
					if val > 1 {
						lp := lightPropagation{neighbourPos.X*16 + int32(x), int32(y), neighbourPos.Z*16 + int32(neighbourZ), val}
						incomingLight = append(incomingLight, lp)
					}
				}
			}
		}
	}

	// stage 2: apply all collected light sources to the current chunk
	for _, lp := range incomingLight {
		if err := propagateLight(chunk.world, lp.x, lp.y, lp.z, lp.value); err != nil {
			return err
		}
	}
	return nil
}

// propagateLight propagates in a world at given coordinates
func propagateLight(world *World, x, y, z int32, value byte) error {
	if value > 15 {
		return fmt.Errorf("light value out of bounds: %v", value)
	}
	if value == 0 {
		return nil
	}

	cp, cx, cy, cz := WorldToChunkLocal(x, y, z)
	chunk, ok := world.Chunk(cp)
	if !ok {
		return fmt.Errorf("chunk at coordinates %v;%v is not loaded", x>>4, z>>4)
	}

	if previous, err := chunk.blockLight.get(cx, cy, cz); err != nil {
		return err
	} else if previous > value { // if the value is equal we still propagate, this is needed when refreshing light sources
		return nil // there is higher light value already present
	}

	if err := chunk.blockLight.set(cx, cy, cz, value); err != nil {
		return err
	}

	var queue = []lightPropagation{{x, y, z, value}}
	for {
		if len(queue) == 0 {
			break
		}

		next := queue[0]
		queue = queue[1:]

		if next.value <= 1 {
			continue // the light would not reach neighbours
		}

		for _, neighbour := range next.neighbours() {
			cp, cx, cy, cz := WorldToChunkLocal(neighbour.x, neighbour.y, neighbour.z)
			chunk, ok := world.Chunk(cp)
			if !ok {
				continue
			}

			current, err := chunk.blockLight.get(cx, cy, cz)
			if err != nil {
				return err
			}
			block, err := chunk.GetBlock(cx, cy, cz)
			if err != nil {
				return err
			}

			if neighbour.value <= block.Material().LightOpacity() {
				continue // the block would block all the light
			}
			neighbour.value -= block.Material().LightOpacity()

			if current >= neighbour.value {
				continue // there is higher or same light value already present
			}

			if err = chunk.blockLight.set(cx, cy, cz, neighbour.value); err != nil {
				return err
			}

			queue = append(queue, neighbour)
		}
	}

	return nil
}

// removeLight removes light at given chunk coordinates
func removeLight(world *World, x, y, z int32) error {
	cp, cx, cy, cz := WorldToChunkLocal(x, y, z)
	chunk, ok := world.Chunk(cp)
	if !ok {
		return fmt.Errorf("chunk at coordinates %v;%v is not loaded", x>>4, z>>4)
	}
	previous, err := chunk.blockLight.get(cx, cy, cz)
	if err != nil {
		return err
	} else if previous == 0 {
		return nil // no need to remove, there is no light present
	}

	if err = chunk.blockLight.set(cx, cy, cz, 0); err != nil {
		return err
	}

	var queue = []lightPropagation{{x, y, z, previous}}
	var toPropagate = make([]lightPropagation, 0)
	var toRecalculate = make([]lightPropagation, 0)
	for {
		if len(queue) == 0 {
			break
		}

		next := queue[0]
		queue = queue[1:]

		if next.value <= 1 {
			continue // the light would not reach neighbours
		}

		for _, neighbour := range next.neighbours() {
			cp, cx, cy, cz := WorldToChunkLocal(neighbour.x, neighbour.y, neighbour.z)
			chunk, ok := world.Chunk(cp)
			if !ok {
				continue
			}

			current, err := chunk.blockLight.get(cx, cy, cz)
			if err != nil {
				return err
			}

			if current == 0 {
				continue // we already removed light from that neighbour
			}

			block, err := chunk.GetBlock(cx, cy, cz)
			if err != nil {
				return err
			}

			if neighbour.value <= block.Material().LightOpacity() {
				continue // the block would block all the light
			}
			neighbour.value -= block.Material().LightOpacity()

			if current > neighbour.value {
				// light from different source
				toPropagate = append(toPropagate, lightPropagation{neighbour.x, neighbour.y, neighbour.z, current})
				continue
			}

			if err = chunk.blockLight.set(cx, cy, cz, 0); err != nil {
				return err
			}

			queue = append(queue, neighbour)
			emission := block.Material().LightEmission()
			if emission != 0 {
				toRecalculate = append(queue, lightPropagation{neighbour.x, neighbour.y, neighbour.z, emission})
			}
		}
	}

	for _, propagate := range toPropagate {
		if err = propagateLight(world, propagate.x, propagate.y, propagate.z, propagate.value); err != nil {
			return err
		}
	}

	for _, recalculate := range toRecalculate {
		if err = propagateLight(world, recalculate.x, recalculate.y, recalculate.z, recalculate.value); err != nil {
			return err
		}
	}

	return nil
}
