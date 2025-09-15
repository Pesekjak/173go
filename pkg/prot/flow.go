package prot

import "github.com/Pesekjak/173go/pkg/buff"

type puller struct {
	buff *buff.MCReader
	err  error
}

func newPuller(buff *buff.MCReader) *puller {
	return &puller{buff: buff}
}

func (p *puller) pull(f func()) {
	if p.err == nil {
		f()
	}
}

type pusher struct {
	buff *buff.MCWriter
	err  error
}

func newPusher(buff *buff.MCWriter) *pusher {
	return &pusher{buff: buff}
}

func (p *pusher) push(f func() error) {
	if p.err == nil {
		p.err = f()
	}
}
