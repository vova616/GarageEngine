package Engine

import (
	"fmt"
	"runtime"
	"time"
)

type Command byte
type Signal chan Command

const (
	Continue = Command(0)
	Close    = Command(1)

	Running = Command(2)
	Ended   = Command(3)
)

type Goroutine struct {
	in    chan Command
	out   chan Command
	State Command
}

func (gr *Goroutine) WaitForCommand() {
	act := <-gr.in
	switch act {
	case Continue:
		break
	case Close:
		for i, ch := range goroutines {
			if ch == gr {
				goroutines[i] = nil
				break
			}
		}
		panic("")
	}
}

var (
	goroutines        []*Goroutine = make([]*Goroutine, 0, 100)
	current           *Goroutine
	runningGoroutines bool = false
)

func StartGoroutine(fnc func()) *Goroutine {
	gr := &Goroutine{make(chan Command), make(chan Command), Running}
	found := false
	for i, ch := range goroutines {
		if ch == nil {
			goroutines[i] = gr
			found = true
			break
		}
	}
	if !found {
		goroutines = append(goroutines, gr)
	}
	go startGoroutine(fnc, gr)
	return gr
}

func startGoroutine(fnc func(), gr *Goroutine) {
	defer errorFunc()
	gr.WaitForCommand()
	fnc()
	gr.out <- Ended
}

func errorFunc() {
	if p := recover(); p != nil && p != "" {
		fmt.Println(p, PanicPath())
	}
}

func YieldSkip() {
	if !runningGoroutines {
		return
	}
	c := current
	c.out <- Running
	c.WaitForCommand()
}

func Yield(Out <-chan Command) {
	if !runningGoroutines {
		return
	}
	c := current
	c.out <- Running

	for {
		select {
		case out := <-Out:
			if out == Ended {
				goto work
			}
		case in := <-c.in:
			if in == Close {
				panic("")
			} else {
				c.out <- Running
			}
		}
	}
work:

	c.WaitForCommand()
}

func NewSignal() Signal {
	return make(chan Command)
}

func (signal Signal) SendEnd() {
	signal <- Ended
}

func Wait(seconds float32) {
	if !runningGoroutines {
		return
	}
	start := time.Now()
	for {
		YieldSkip()
		now := time.Now()
		if now.Sub(start).Seconds() >= float64(seconds) {
			break
		}
	}
	/*
		Signal := NewSignal()
		go func() {
			<-time.After(time.Second * 3)
			Signal.SendEnd()
		}() 
		Yield(Signal)
	*/
}

func RunGoroutines() {
	runningGoroutines = true
	for i, ch := range goroutines {
		if ch != nil {
			current = ch
			ch.in <- Continue
			state := <-ch.out
			if state == Ended {
				goroutines[i] = nil
				ch.State = Ended
			}
		}
	}
	runningGoroutines = false
	current = nil
}

func PanicPath() string {
	fullPath := ""
	skip := 3
	for i := skip; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if i > skip {
			fullPath += ", "
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		fullPath += fmt.Sprintf("%s:%d", file, line)
	}
	return fullPath
}
