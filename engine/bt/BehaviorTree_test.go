package bt

import (
	//"runtime"
	"testing"
)

/*
	StartBehavior(Sleep(5), func() Command { println("asd"); return Continue }, Sleep(5), func() Command { println("asd"); return Restart })
	go func() {
		for {
			RunBT(10)
			time.Sleep(time.Millisecond * 200)
		}
	}()
*/

func Benchmark_BT(b *testing.B) {
	b.StopTimer()

	n := 0
	fastShoot := func() Command {
		if n < 100*b.N {
			n++
			return Yield
		}
		return Close
	}

	for i := 0; i < b.N; i++ {
		Start(fastShoot)
	}

	b.StartTimer()
	for len(Routines) > 0 {
		Run(1)
	}

	if n != 100*b.N {
		b.Fatalf("n is %d need %d", n, b.N*100)
	}
}

func Test_BT(t *testing.T) {
	n := 0
	fastShoot := func() Command {
		if n < 100 {
			n++
			return Yield
		}
		return Close
	}

	Start(fastShoot)

	Run(101)

	if len(Routines) != 0 {
		t.Fatalf("Routines len is %d need 0", len(Routines))
	}

	if n != 100 {
		t.Fatalf("n is %d need %d", n, 100)
	}
}
