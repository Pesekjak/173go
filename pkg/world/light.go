package world

import "fmt"

type light struct {
	data []byte
}

func newLight() *light {
	return &light{data: make([]byte, (ChunkSize*ChunkHeight*ChunkSize)/2)}
}

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

type lightPropagation struct {
	x, y, z uint32
	value   byte
}

func (qe *lightPropagation) neighbours() []lightPropagation {
	neighbours := make([]lightPropagation, 0, 6)
	if qe.x > 0 {
		neighbours = append(neighbours, lightPropagation{qe.x - 1, qe.y, qe.z, qe.value - 1})
	}
	if qe.x < ChunkSize-1 {
		neighbours = append(neighbours, lightPropagation{qe.x + 1, qe.y, qe.z, qe.value - 1})
	}
	if qe.z > 0 {
		neighbours = append(neighbours, lightPropagation{qe.x, qe.y, qe.z - 1, qe.value - 1})
	}
	if qe.z < ChunkSize-1 {
		neighbours = append(neighbours, lightPropagation{qe.x, qe.y, qe.z + 1, qe.value - 1})
	}
	if qe.y > 0 {
		neighbours = append(neighbours, lightPropagation{qe.x, qe.y - 1, qe.z, qe.value - 1})
	}
	if qe.y < ChunkHeight-1 {
		neighbours = append(neighbours, lightPropagation{qe.x, qe.y + 1, qe.z, qe.value - 1})
	}
	return neighbours
}

func lightChunkBorders(world *World, chunk *Chunk) error {
	chunkPos := chunk.Pos()
	var incomingLight []lightPropagation

	// stage 1: gather all light sources from  neighbours
	// we also consider diagonally positioned chunks neighbours
	// because lighting from the chunk can reach there too
	for dz := -1; dz <= 1; dz++ {
		for dx := -1; dx <= 1; dx++ {

			// skip "this" chunk
			if dx == 0 && dz == 0 {
				continue
			}

			neighbourPos := chunkPos.Add(ChunkPos{int32(dx), int32(dz)})
			neighbour, ok := world.Chunk(neighbourPos)
			if !ok {
				continue // neighbour is not loaded
			}

			var chunkX, neighbourX, chunkZ, neighbourZ uint32
			if dx == -1 {
				chunkX, neighbourX = 0, ChunkSize-1
			} else if dx == 1 {
				chunkX, neighbourX = ChunkSize-1, 0
			}
			if dz == -1 {
				chunkZ, neighbourZ = 0, ChunkSize-1
			} else if dz == 1 {
				chunkZ, neighbourZ = ChunkSize-1, 0
			}

			// loop along the shared border or edge
			for y := uint32(0); y < ChunkHeight; y++ {
				// case 1: diagonal neighbour
				if dx != 0 && dz != 0 {
					val, err := neighbour.blockLight.get(neighbourX, y, neighbourZ)
					if err != nil {
						return err
					}
					if val > 1 {
						incomingLight = append(incomingLight, lightPropagation{chunkX, y, chunkZ, val - 1})
					}
					// case 2: cardinal neighbour
				} else if dx != 0 { // x face
					for z := uint32(0); z < ChunkSize; z++ {
						val, err := neighbour.blockLight.get(neighbourX, y, z)
						if err != nil {
							return err
						}
						if val > 1 {
							incomingLight = append(incomingLight, lightPropagation{chunkX, y, z, val - 1})
						}
					}
				} else { // z face
					for x := uint32(0); x < ChunkSize; x++ {
						val, err := neighbour.blockLight.get(x, y, neighbourZ)
						if err != nil {
							return err
						}
						if val > 1 {
							incomingLight = append(incomingLight, lightPropagation{x, y, chunkZ, val - 1})
						}
					}
				}
			}
		}
	}

	// stage 2: apply all collected light sources to the current chunk
	for _, lp := range incomingLight {
		if err := chunk.blockLight.propagate(chunk, lp.x, lp.y, lp.z, lp.value); err != nil {
			return err
		}
	}

	// stage 3: propagate updated light values outwards to all 8 neighbours (including collected light sources plus
	// what was already there)
	for dz := -1; dz <= 1; dz++ {
		for dx := -1; dx <= 1; dx++ {

			// skip "this" chunk
			if dx == 0 && dz == 0 {
				continue
			}

			neighbourPos := chunkPos.Add(ChunkPos{int32(dx), int32(dz)})
			neighbour, ok := world.Chunk(neighbourPos)
			if !ok {
				continue // neighbour is not loaded
			}

			var chunkX, neighbourX, chunkZ, neighbourZ uint32
			if dx == -1 {
				chunkX, neighbourX = 0, ChunkSize-1
			} else if dx == 1 {
				chunkX, neighbourX = ChunkSize-1, 0
			}
			if dz == -1 {
				chunkZ, neighbourZ = 0, ChunkSize-1
			} else if dz == 1 {
				chunkZ, neighbourZ = ChunkSize-1, 0
			}

			for y := uint32(0); y < ChunkHeight; y++ {
				// case 1: diagonal neighbour
				if dx != 0 && dz != 0 {
					val, err := chunk.blockLight.get(chunkX, y, chunkZ)
					if err != nil {
						return err
					}
					if val > 1 {
						if err := neighbour.blockLight.propagate(neighbour, neighbourX, y, neighbourZ, val-1); err != nil {
							return err
						}
					}
					// case 2: cardinal neighbour
				} else if dx != 0 { // x face
					for z := uint32(0); z < ChunkSize; z++ {
						val, err := chunk.blockLight.get(chunkX, y, z)
						if err != nil {
							return err
						}
						if val > 1 {
							if err := neighbour.blockLight.propagate(neighbour, neighbourX, y, z, val-1); err != nil {
								return err
							}
						}
					}
				} else { // z face
					for x := uint32(0); x < ChunkSize; x++ {
						val, err := chunk.blockLight.get(x, y, chunkZ)
						if err != nil {
							return err
						}
						if val > 1 {
							if err := neighbour.blockLight.propagate(neighbour, x, y, neighbourZ, val-1); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (l *light) propagate(chunk *Chunk, x, y, z uint32, value byte) error {
	if value > 15 {
		return fmt.Errorf("light value out of bounds: %v", value)
	}
	if _, err := inChunkBounds(x, y, z); err != nil {
		return err
	}
	if value == 0 {
		return nil
	}

	if previous, err := l.get(x, y, z); err != nil {
		return err
	} else if previous >= value {
		return nil // there is higher light value already present
	}

	if err := l.set(x, y, z, value); err != nil {
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
			current, err := l.get(neighbour.x, neighbour.y, neighbour.z)
			if err != nil {
				return err
			}
			block, err := chunk.GetBlock(neighbour.x, neighbour.y, neighbour.z)
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

			if err = l.set(neighbour.x, neighbour.y, neighbour.z, neighbour.value); err != nil {
				return err
			}

			queue = append(queue, neighbour)
		}
	}

	return nil
}

func (l *light) remove(chunk *Chunk, x, y, z uint32) error {
	if _, err := inChunkBounds(x, y, z); err != nil {
		return err
	}

	previous, err := l.get(x, y, z)
	if err != nil {
		return err
	} else if previous == 0 {
		return nil // no need to remove, there is no light present
	}

	if err = l.set(x, y, z, 0); err != nil {
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
			current, err := l.get(neighbour.x, neighbour.y, neighbour.z)
			if err != nil {
				return err
			}

			if current == 0 {
				continue // we already removed light from that neighbour
			}

			block, err := chunk.GetBlock(neighbour.x, neighbour.y, neighbour.z)
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

			if err = l.set(neighbour.x, neighbour.y, neighbour.z, 0); err != nil {
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
		if err = l.propagate(chunk, propagate.x, propagate.y, propagate.z, propagate.value); err != nil {
			return err
		}
	}

	for _, recalculate := range toRecalculate {
		if err = l.propagate(chunk, recalculate.x, recalculate.y, recalculate.z, recalculate.value); err != nil {
			return err
		}
	}

	return nil
}
