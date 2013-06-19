package cr

import (
	"fmt"
	"runtime"
	"time"
)

type Command byte

const (
	Continue = Command(iota)
	Close    = Command(iota)
	Idle     = Command(iota)
)

var (
	coroutines        []*Coroutine = make([]*Coroutine, 0, 100)
	current           *Coroutine
	index             int       = 0
	end               chan bool = make(chan bool)
	runningCoroutines bool      = false
)

type Coroutine struct {
	in       chan Command
	State    Command
	UserData interface{}
}

func (gr *Coroutine) WaitForCommand() {
	gr.State = Idle
	act := <-gr.in
	gr.State = act
	switch act {
	case Continue:
		break
	case Close:
		coroutines[len(coroutines)-1], coroutines[index], coroutines = nil, coroutines[len(coroutines)-1], coroutines[:len(coroutines)-1]
	}
}

func Start(fnc func()) *Coroutine {
	gr := &Coroutine{make(chan Command, 1), Idle, nil}
	coroutines = append(coroutines, gr)
	go start(fnc, gr)
	return gr
}

func runNext() {
	index++
	if index < len(coroutines) {
		current = coroutines[index]
		current.in <- Continue
	} else {
		end <- true
	}
}

func start(fnc func(), gr *Coroutine) {
	defer errorFunc()
	defer gr.endFunc()
	gr.WaitForCommand()
	fnc()
}

func (gr *Coroutine) endFunc() {
	gr.State = Close
	coroutines[len(coroutines)-1], coroutines[index], coroutines = nil, coroutines[len(coroutines)-1], coroutines[:len(coroutines)-1]
	runNext()
}

func errorFunc() {
	if p := recover(); p != nil {
		fmt.Println(p, PanicPath())
	}
}

func Clear() {
	coroutines = coroutines[:0]
}

func Skip() {
	if !runningCoroutines {
		return
	}
	c := current
	runNext()
	c.WaitForCommand()
}

func YieldCoroutine(gr *Coroutine) {
	if !runningCoroutines {
		return
	}
	for gr.State != Close {
		Skip()
	}
}

func YieldUntil(done <-chan bool) {
	if !runningCoroutines {
		return
	}
	for {
		select {
		case <-done:
			return
		default:
			Skip()
		}
	}
}

func Sleep(seconds float32) {
	if !runningCoroutines {
		return
	}
	start := time.Now()
	for {
		Skip()
		if time.Since(start).Seconds() >= float64(seconds) {
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

func Run() {
	runningCoroutines = true
	index = 0
	if len(coroutines) > 0 {
		current = coroutines[0]
		coroutines[0].in <- Continue
		<-end
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
