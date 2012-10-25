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

func (r *Routine) Run() (Command, bool) {
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
		r.CurrentFunc = 0
		return b, true
	}
	if r.CurrentFunc >= len(r.Funcs) {
		r.CurrentFunc = 0
		return b, false
	}
	return b, false
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
	Run() (Command, bool)
}

func WaitContinue(fnc RoutineFunc, child Routiner, secTimeout float32) RoutineFunc {
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
		if fnc != nil && fnc() == Continue && child != nil {
			c, _ := child.Run()
			if c != Yield {
				started = false
			}
			return c
		}
		return Yield
	}
}

func Sequence(funcs ...RoutineFunc) RoutineFunc {
	r := NewBehavior(funcs...)
	return func() Command {
		_, stop := r.Run()
		if !stop {
			return Yield
		}
		return Close
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
		NewBehavior(func() Command {
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
		}), secs)
}

func NewBehavior(funcs ...RoutineFunc) *Routine {
	r := &Routine{0, 0, funcs}
	return r
}

func StartBehavior(funcs ...RoutineFunc) *Routine {
	r := NewBehavior(funcs...)
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
				_, delete := r.Run()
				if delete {
					Routines[index] = nil
				}
			}
		}
	}
}
