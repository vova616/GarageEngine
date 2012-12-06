package Server

import (
	"sync"
	"sync/atomic"
)

type IDGenerator struct {
	lastID     int
	generating int32
	ids        chan int
	queuedIDs  []int
	queueMutex *sync.Mutex
}

func NewIDGenerator(buffer int) *IDGenerator {
	id := &IDGenerator{0, 0, make(chan int, buffer), make([]int, 0, 0), &sync.Mutex{}}
	id.GenIDs()
	id.GenIDs()
	return id
}

func (gen *IDGenerator) GenIDs() {
	defer atomic.StoreInt32(&gen.generating, 0)
	for i := 0; i < cap(gen.ids)/2; i++ {
		gen.genFromQueue()
		select {
		case gen.ids <- gen.lastID:
			gen.lastID++
		default:
			return
		}
	}
}

func (gen *IDGenerator) genFromQueue() {
	if len(gen.queuedIDs) > 0 {
		gen.queueMutex.Lock()
		defer gen.queueMutex.Unlock()
		if len(gen.queuedIDs) > 0 {
			for len(gen.queuedIDs) > 0 {
				var id int
				id, gen.queuedIDs = gen.queuedIDs[len(gen.queuedIDs)-1], gen.queuedIDs[:len(gen.queuedIDs)-1]
				select {
				case gen.ids <- id:
				default:
					gen.queueMutex.Unlock()
					return
				}
			}
		}
	}
}

func (gen *IDGenerator) NextID() int {
	id := 0

	if len(gen.queuedIDs) > 0 {
		gen.queueMutex.Lock()
		if len(gen.queuedIDs) > 0 {
			id, gen.queuedIDs = gen.queuedIDs[len(gen.queuedIDs)-1], gen.queuedIDs[:len(gen.queuedIDs)-1]
		} else {
			id = <-gen.ids
		}
		gen.queueMutex.Unlock()
	} else {
		id = <-gen.ids
	}

	if len(gen.ids) == 0 {
		if atomic.CompareAndSwapInt32(&gen.generating, 0, 1) {
			go gen.GenIDs()
		}
	}
	return id
}

func (gen *IDGenerator) PutID(id int) {
	select {
	case gen.ids <- id:
	default:
		gen.queueMutex.Lock()
		gen.queuedIDs = append(gen.queuedIDs, id)
		gen.queueMutex.Unlock()
	}
}
