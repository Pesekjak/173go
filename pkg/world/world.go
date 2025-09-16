package world

type World struct {
	Dimension  Dimension
	SpawnPoint BlockPos
	Time       int64
}

func NewWorld() *World {
	return &World{
		Dimension:  Overworld,
		SpawnPoint: NewBlockPos(0, 60, 0),
		Time:       0,
	}
}
