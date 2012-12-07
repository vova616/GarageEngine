package Engine

import (
	"fmt"
	"runtime"
	"time"
)

type Command byte
type Signal chan Command

const (
	Continue = Command(1)
	Close    = Command(2)

	Running = Command(4)
	Ended   = Command(8)

	Yield   = Command(16)
	Restart = Command(32)
)

var (
	coroutines        []*Coroutine = make([]*Coroutine, 0, 100)
	current           *Coroutine
	runningCoroutines bool = false
)

type Coroutine struct {
	in       chan Command
	out      chan Command
	State    Command
	UserData interface{}
}

func (gr *Coroutine) WaitForCommand() {
	act := <-gr.in
	switch act {
	case Continue:
		break
	case Close:
		for i, ch := range coroutines {
			if ch == gr {
				coroutines[i] = nil
				break
			}
		}
		panic("")
	}
}

func StartCoroutine(fnc func()) *Coroutine {
	gr := &Coroutine{make(chan Command), make(chan Command), Running, nil}
	found := false
	for i, ch := range coroutines {
		if ch == nil {
			coroutines[i] = gr
			found = true
			break
		}
	}
	if !found {
		coroutines = append(coroutines, gr)
	}
	go startCoroutine(fnc, gr)
	return gr
}

func startCoroutine(fnc func(), gr *Coroutine) {
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

func CoYieldSkip() {
	if !runningCoroutines {
		return
	}
	c := current
	c.out <- Running
	c.WaitForCommand()
}

func CoYieldCoroutine(gr *Coroutine) {
	if !runningCoroutines {
		return
	}
	for gr.State != Ended {
		CoYieldSkip()
	}
}

func CoYieldUntil(Out <-chan Command) {
	if !runningCoroutines {
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

func CoSleep(seconds float32) {
	if !runningCoroutines {
		return
	}
	start := time.Now()
	for {
		CoYieldSkip()
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

func RunCoroutines() {
	runningCoroutines = true
	for i, ch := range coroutines {
		if ch != nil {
			current = ch
			ch.in <- Continue
			state := <-ch.out
			if state == Ended {
				coroutines[i] = nil
				ch.State = Ended
			}
		}
	}
	runningCoroutines = false
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
