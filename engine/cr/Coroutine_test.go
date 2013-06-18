package cr

import (
	"runtime"
	"testing"
)

func Benchmark_Coroutine(b *testing.B) {
	b.StopTimer()
	n := 0
	fastShoot := func() {
		for i := 0; i < 100; i++ {
			n++
			Skip()
		}
	}

	for i := 0; i < b.N; i++ {
		Start(fastShoot)
	}

	b.StartTimer()
	for len(coroutines) > 0 {
		Run()
	}

	if n != 100*b.N {
		b.Fatalf("n is %d need %d", n, b.N*100)
	}
}

func Benchmark_Coroutine_MaxProcs(b *testing.B) {
	b.StopTimer()
	runtime.GOMAXPROCS(runtime.NumCPU())
	n := 0
	fastShoot := func() {
		for i := 0; i < 100; i++ {
			n++
			Skip()
		}
	}

	for i := 0; i < b.N; i++ {
		Start(fastShoot)
	}

	b.StartTimer()
	for len(coroutines) > 0 {
		Run()
	}

	if n != 100*b.N {
		b.Fatalf("n is %d need %d", n, b.N*100)
	}
}

func Test_Coroutine(t *testing.T) {
	n := 0
	fastShoot := func() {
		for i := 0; i < 100; i++ {
			n++
			Skip()
		}
	}

	Start(fastShoot)

	for i := 0; i < 101; i++ {
		Run()
	}

	if len(coroutines) > 0 {
		t.Fatalf("coroutines len is %d need 0", len(coroutines))
	}

	if n != 100 {
		t.Fatalf("n is %d need 100", n)
	}
}
