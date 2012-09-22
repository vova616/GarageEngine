package Engine

import (
	"time"
)

type Timer map[interface{}]time.Time

const defaultKey = "DefTimer"

func NewTimer() Timer {
	return make(map[interface{}]time.Time)
}

func (timer Timer) Start() {
	timer[defaultKey] = time.Now()
}

func (timer Timer) StartCustom(key interface{}) {
	if key == defaultKey {
		panic("Cannot use " + defaultKey + " as a key.")
	}
	timer[key] = time.Now()
}

func (timer Timer) StopCustom(key interface{}) time.Duration {
	startTime, exist := timer[key]
	if !exist {
		panic("No such custom key")
	}
	now := time.Now()
	return now.Sub(startTime)
}

func (timer Timer) Stop() time.Duration {
	startTime, exist := timer[defaultKey]
	if !exist {
		panic("You must start the timer before stopping it")
	}
	now := time.Now()
	return now.Sub(startTime)
}
