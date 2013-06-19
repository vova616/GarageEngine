package bt

import (
	//"log"
	"math/rand"
	"time"
)

var (
	Routines []Routiner = make([]Routiner, 0)
)

type Command byte

const (
	Continue = Command(iota)
	Close    = Command(iota)
	Yield    = Command(iota)
	Restart  = Command(iota)
)

type RoutineFunc func() Command

type Routine struct {
	CurrentFunc int
	Status      Command
	Funcs       []RoutineFunc
}

func Clear() {
	Routines = Routines[:0]
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
	r := New(funcs...)
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
		New(func() Command {
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

func New(funcs ...RoutineFunc) *Routine {
	r := &Routine{0, 0, funcs}
	return r
}

func Start(funcs ...RoutineFunc) *Routine {
	r := New(funcs...)
	Routines = append(Routines, r)
	return r
}

func Run(ticks int) {
	for i := 0; i < ticks; i++ {
		for index := 0; index < len(Routines); index++ {
			r := Routines[index]
			if r != nil {
				_, delete := r.Run()
				if delete {
					Routines[len(Routines)-1], Routines[index], Routines = nil, Routines[len(Routines)-1], Routines[:len(Routines)-1]
					index--
				}
			}
		}
	}
}
