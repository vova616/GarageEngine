package Server

import (
	"sync"
)

type ID int32

type IDGenerator struct {
	lastID ID
	ids    []ID
	locker *sync.Mutex
}

func NewIDGenerator(buffer int) *IDGenerator {
	gen := &IDGenerator{0, make([]ID, 0, buffer), &sync.Mutex{}}
	gen.locker.Lock()
	gen.genIDs()
	gen.locker.Unlock()
	return gen
}

func (gen *IDGenerator) genIDs() {
	for i := 0; i < cap(gen.ids)/2; i++ {
		gen.ids = append(gen.ids, gen.lastID)
		gen.lastID++
	}
}

func (gen *IDGenerator) NextID() ID {
	id := ID(0)

	gen.locker.Lock()
	id, gen.ids = gen.ids[len(gen.ids)-1], gen.ids[:len(gen.ids)-1]
	if len(gen.ids) == 0 {
		gen.genIDs()
	}
	gen.locker.Unlock()

	return id
}

func (gen *IDGenerator) PutID(id ID) {
	gen.locker.Lock()
	gen.ids = append(gen.ids, id)
	gen.locker.Unlock()
}
