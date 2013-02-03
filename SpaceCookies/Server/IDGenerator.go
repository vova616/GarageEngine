package server

import (
	"sync"
)

type ID int32

type IDGenerator struct {
	lastID     ID
	ids        []ID
	locker     *sync.Mutex
	threadSafe bool
}

func NewIDGenerator(buffer int, threadSafe bool) *IDGenerator {
	gen := &IDGenerator{0, make([]ID, 0, buffer), &sync.Mutex{}, threadSafe}
	if threadSafe {
		gen.locker.Lock()
		defer gen.locker.Unlock()
	}

	gen.genIDs()
	return gen
}

func (gen *IDGenerator) genIDs() {
	gen.lastID += ID(cap(gen.ids) / 2)
	id := gen.lastID - 1
	for i := 0; i < cap(gen.ids)/2; i++ {
		gen.ids = append(gen.ids, id)
		id--
	}
}

func (gen *IDGenerator) NextID() ID {
	id := ID(0)

	if gen.threadSafe {
		gen.locker.Lock()
		defer gen.locker.Unlock()
	}
	id, gen.ids = gen.ids[len(gen.ids)-1], gen.ids[:len(gen.ids)-1]
	if len(gen.ids) == 0 {
		gen.genIDs()
	}

	return id
}

func (gen *IDGenerator) PutID(id ID) {
	if gen.threadSafe {
		gen.locker.Lock()
		defer gen.locker.Unlock()
	}
	gen.ids = append(gen.ids, id)
}
