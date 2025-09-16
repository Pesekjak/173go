package buff

// Puller is a utility for reading multiple elements from MCReader while respecting the possibly returned errors.
// Example usage: puller.Pull(func() { object.Data, puller.Err = buf.ReadByte() })
// Once the data pulling is finished, the possibly returned error is stored in Err, if not nil, the pulling failed
// at some point.
type Puller struct {
	buff *MCReader
	// Error that might have occurred during the continuous data pulling
	Err error
}

// NewPuller creates new Puller wrapped around MCReader
func NewPuller(buff *MCReader) *Puller {
	return &Puller{buff: buff}
}

// Pull pulls data from the wrapped MCReader
func (p *Puller) Pull(f func()) {
	if p.Err == nil {
		f()
	}
}

// Pusher is a utility for writing multiple elements to MCWriter while respecting the possibly returned errors.
// Example usage: pusher.Push(func() error { return buf.WriteByte(object.Data) })
// Once the data pushing is finished, the possibly returned error is stored in Err, if not nil, the pushing failed
// at some point.
type Pusher struct {
	buff *MCWriter
	// Error that might have occurred during the continuous data pushing
	Err error
}

// NewPusher creates new Pusher wrapped around MCWriter
func NewPusher(buff *MCWriter) *Pusher {
	return &Pusher{buff: buff}
}

// Push pushes data to the wrapped MCWriter
func (p *Pusher) Push(f func() error) {
	if p.Err == nil {
		p.Err = f()
	}
}
