package base

import (
	"fmt"
	"math"
)

const ChunkWidth = 16

type Location struct {
	X     float64
	Y     float64
	Z     float64
	Yaw   float32
	Pitch float32
}

func NewLocation(x, y, z float64, yaw, pitch float32) Location {
	return Location{X: x, Y: y, Z: z, Yaw: yaw, Pitch: pitch}
}

func (l Location) ToBlockPos() BlockPos {
	return BlockPos{
		X: int(math.Floor(l.X)),
		Y: int(math.Floor(l.Y)),
		Z: int(math.Floor(l.Z)),
	}
}

func (l Location) Add(x, y, z float64) Location {
	return Location{
		X:     l.X + x,
		Y:     l.Y + y,
		Z:     l.Z + z,
		Yaw:   l.Yaw,
		Pitch: l.Pitch,
	}
}

func (l Location) Subtract(other Location) Location {
	return Location{
		X:     l.X - other.X,
		Y:     l.Y - other.Y,
		Z:     l.Z - other.Z,
		Yaw:   l.Yaw,
		Pitch: l.Pitch,
	}
}

func (l Location) DistanceTo(other Location) float64 {
	return math.Sqrt(l.DistanceToSquared(other))
}

func (l Location) DistanceToSquared(other Location) float64 {
	dx := l.X - other.X
	dy := l.Y - other.Y
	dz := l.Z - other.Z
	return dx*dx + dy*dy + dz*dz
}

func (l Location) DirectionVector() Location {

	yawRad := float64(l.Yaw * math.Pi / 180.0)
	pitchRad := float64(l.Pitch * math.Pi / 180.0)

	x := -math.Sin(yawRad) * math.Cos(pitchRad)
	y := -math.Sin(pitchRad)
	z := math.Cos(yawRad) * math.Cos(pitchRad)

	return Location{X: x, Y: y, Z: z}
}

func (l Location) String() string {
	return fmt.Sprintf("Location(X: %.2f, Y: %.2f, Z: %.2f, Yaw: %.1f, Pitch: %.1f)", l.X, l.Y, l.Z, l.Yaw, l.Pitch)
}

type BlockPos struct {
	X int
	Y int
	Z int
}

func NewBlockPos(x, y, z int) BlockPos {
	return BlockPos{X: x, Y: y, Z: z}
}

func (p BlockPos) Add(other BlockPos) BlockPos {
	return BlockPos{X: p.X + other.X, Y: p.Y + other.Y, Z: p.Z + other.Z}
}

func (p BlockPos) Offset(dx, dy, dz int) BlockPos {
	return BlockPos{X: p.X + dx, Y: p.Y + dy, Z: p.Z + dz}
}

func (p BlockPos) Up(n int) BlockPos {
	return p.Offset(0, n, 0)
}

func (p BlockPos) Down(n int) BlockPos {
	return p.Offset(0, -n, 0)
}

func (p BlockPos) North(n int) BlockPos {
	return p.Offset(0, 0, -n)
}

func (p BlockPos) South(n int) BlockPos {
	return p.Offset(0, 0, n)
}

func (p BlockPos) East(n int) BlockPos {
	return p.Offset(n, 0, 0)
}

func (p BlockPos) West(n int) BlockPos {
	return p.Offset(-n, 0, 0)
}

func (p BlockPos) ToChunkPos() ChunkPos {
	chunkX := int(math.Floor(float64(p.X) / float64(ChunkWidth)))
	chunkZ := int(math.Floor(float64(p.Z) / float64(ChunkWidth)))
	return ChunkPos{X: chunkX, Z: chunkZ}
}

func (p BlockPos) String() string {
	return fmt.Sprintf("BlockPos(X: %d, Y: %d, Z: %d)", p.X, p.Y, p.Z)
}

type ChunkPos struct {
	X int
	Z int
}

func (c ChunkPos) String() string {
	return fmt.Sprintf("ChunkPos(X: %d, Z: %d)", c.X, c.Z)
}
