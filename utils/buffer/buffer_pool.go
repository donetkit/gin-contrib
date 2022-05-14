package buffer

import (
	"bytes"
	"sync"
)

// Pool is Pool of *bytes.Buffer
type Pool struct {
	pool sync.Pool
}

// Get a bytes.Buffer pointer
func (p *Pool) Get() *bytes.Buffer {
	buf := p.pool.Get()
	if buf == nil {
		return &bytes.Buffer{}
	}
	return buf.(*bytes.Buffer)
}

// Put a bytes.Buffer pointer to BufferPool
func (p *Pool) Put(buf *bytes.Buffer) {
	p.pool.Put(buf)
}
