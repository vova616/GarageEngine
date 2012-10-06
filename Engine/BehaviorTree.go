package Engine

/*
TODO:
WRITE IT LOL





*/

import (
	//"log"
	"math/rand"
	"time"
)

var (
	Routines []Routiner = make([]Routiner, 0)
)

type RoutineFunc func() Command

type Routine struct {
	CurrentFunc int
	Status      Command
	Funcs       []RoutineFunc
}

func (r *Routine) Run() bool {
	b := r.Funcs[r.CurrentFunc]()
	switch b {
	case Continue:
		r.CurrentFunc++
	case Restart:
		r.CurrentFunc = 0
	case Yield:
		break
	case Close:
		r.CurrentFunc = len(r.Funcs)
	}
	if r.CurrentFunc >= len(r.Funcs) {
		return false
	}
	return true
}

func init() {
	return
	StartBehavior(Sleep(5), func() Command { println("asd"); return Continue }, Sleep(5), func() Command { println("asd"); return Restart })
	go func() {
		for {
			RunBT(10)
			time.Sleep(time.Millisecond * 200)
		}
	}()
}

type Routiner interface {
	Run() bool
}

func WaitContinue(fnc RoutineFunc, child RoutineFunc, secTimeout float32) RoutineFunc {
	started := false
	var start time.Time
	return func() Command {
		if !started {
			started = true
			start = time.Now()
		}
		now := time.Now()
		if now.Sub(start).Seconds() > float64(secTimeout) {
			started = false
			return Continue
		}
		if fnc != nil && child != nil && fnc() == Continue {
			c := child()
			if c != Yield {
				started = false
			}
			return c
		}
		return Yield
	}
}

func Sleep(secs float32) RoutineFunc {
	return WaitContinue(func() Command { return Yield }, nil, secs)
}

func SleepRand(secs float32) RoutineFunc {
	started := false
	originalValue := secs
	var start time.Time
	return WaitContinue(func() Command { return Continue },
		func() Command {
			if !started {
				started = true
				start = time.Now()
				secs = rand.Float32() * originalValue
			}
			now := time.Now()
			if now.Sub(start).Seconds() > float64(secs) {
				started = false
				return Continue
			}
			return Yield
		}, secs)
}

func StartBehavior(funcs ...RoutineFunc) *Routine {
	r := &Routine{0, 0, funcs}
	found := false
	for i, ch := range Routines {
		if ch == nil {
			Routines[i] = r
			found = true
			break
		}
	}
	if !found {
		Routines = append(Routines, r)
	}
	return r
}

func RunBT(ticks int) {
	for i := 0; i < ticks; i++ {
		for index, r := range Routines {
			if r != nil {
				if !r.Run() {
					Routines[index] = nil
				}
			}
		}
	}
}
