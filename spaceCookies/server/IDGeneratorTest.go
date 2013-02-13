package server

import "runtime"
import "time"
import "log"

func TestGenerator() {
	runtime.GOMAXPROCS(8)
	defer log.Println("Done!")
	size := 100000
	ch := make(chan ID, size+1)

	generator := NewIDGenerator(size/2, true)
	for i := 0; i < size; i++ {
		go func() {
			time.Sleep(time.Millisecond)
			id := generator.NextID()
			generator.PutID(id)
			ch <- id
		}()
	}
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
