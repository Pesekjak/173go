package world

import "fmt"

type light struct {
	height byte
	data   []byte
}

func newLight(height byte) *light {
	return &light{height: height, data: make([]byte, ChunkSize*uint32(height)*ChunkSize)}
}

func (l *light) get(x, y, z uint32) (byte, error) {
	if x >= 16 || z >= 16 || y >= uint32(l.height) {
		return 0, fmt.Errorf("coordinates %v;%v;%v out of bounds: ", x, y, z)
	}
	i := blockIndex(x, y, z, l.height)
	val := l.data[i/2]
	if i%2 == 0 {
		return val & 0x0F, nil
	} else {
		return val & 0xF0, nil
	}
}

func (l *light) set(x, y, z uint32, value byte) error {
	if value > 15 {
		return fmt.Errorf("light value out of bounds: %v", value)
	}
	if x >= 16 || z >= 16 || y >= uint32(l.height) {
		return fmt.Errorf("coordinates %v;%v;%v out of bounds: ", x, y, z)
	}
	i := blockIndex(x, y, z, l.height)
	val := l.data[i/2]
	if i%2 == 0 {
		l.data[i/2] = (val & 0xF0) | value
	} else {
		l.data[i/2] = (val & 0x0F) | (value << 4)
	}
	return nil
}
