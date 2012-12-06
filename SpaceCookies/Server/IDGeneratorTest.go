package Server

import "runtime"
import "time"
import "log"
import "math/rand"

func TestGenerator() {
	runtime.GOMAXPROCS(8)
	defer log.Println("Done!")
	size := 500000
	ch := make(chan int, size+1)

	generator := NewIDGenerator(size / 2)
	log.Println(len(generator.ids))
	for i := 0; i < size; i++ {
		go func() {
			time.Sleep(time.Millisecond)
			id := generator.NextID()
			time.Sleep(time.Duration(float32(time.Second) * rand.Float32()))
			generator.PutID(id)
			log.Println(len(generator.queuedIDs))
			ch <- id
		}()
	}
	log.Println(len(generator.queuedIDs))
	for {
		if len(ch) != size {
			time.Sleep(time.Second)
			continue
		}
		for id := range ch {
			log.Println(id, len(ch))
			if len(ch) == 0 {
				return
			}
		}
	}

}
