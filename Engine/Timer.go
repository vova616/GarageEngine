package engine

import (
	"time"
)

type Timer map[interface{}]time.Time

const defaultKey = "0-=|DefTimer|=-0"

func NewTimer() Timer {
	return make(map[interface{}]time.Time)
}

func (timer Timer) Start() {
	timer.StartCustom(defaultKey)
}

func (timer Timer) Stop() time.Duration {
	return timer.StopCustom(defaultKey)
}

func (timer Timer) Defer(result *time.Duration) func() {
	return timer.DeferCustom(defaultKey, result)
}

func (timer Timer) DeferCustom(key interface{}, result *time.Duration) func() {
	if result == nil {
		panic("Result is nil.")
	}
	timer.StartCustom(key)
	return func() {
		*result = timer.StopCustom(key)
	}
}

func (timer Timer) StartCustom(key interface{}) {
	timer[key] = time.Now()
}

func (timer Timer) StopCustom(key interface{}) time.Duration {
	startTime, exist := timer[key]
	if !exist {
		panic("No such custom key")
	}
	return time.Since(startTime)
}
