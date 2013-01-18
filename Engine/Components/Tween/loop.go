package Tween

import (
	"github.com/vova616/GarageEngine/Engine"
	//"math"
	"time"
)

func None(t *Tween) bool {
	if t.progress >= 1 {
		t.progress = 1
		return true
	}
	return false
}
func Loop(t *Tween) bool {
	if t.progress >= 1 {
		t.progress = t.progress - 1
		t.startTime = Engine.GameTime().Add(time.Duration(float64(t.progress) * float64(t.Time)))
	}
	return false
}
func PingPong(t *Tween) bool {
	if t.progress >= 1 && !t.reverse {
		t.reverse = true
		t.progress = t.progress - 1
		t.startTime = Engine.GameTime().Add(time.Duration(float64(t.progress) * float64(t.Time)))
	} else if t.progress <= 0 && t.reverse {
		t.reverse = false
		t.progress = -t.progress
		t.startTime = Engine.GameTime().Add(time.Duration(float64(t.progress) * float64(t.Time)))
	}
	return false
}
