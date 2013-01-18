package Tween

import (
//"github.com/vova616/GarageEngine/Engine"
//"math"
//"time"
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
		t.startTime = t.startTime.Add(t.Time)
	}
	return false
}
func PingPong(t *Tween) bool {
	if t.progress >= 1 && !t.reverse {
		t.reverse = true
		t.progress = 2 - t.progress
		t.startTime = t.startTime.Add(t.Time)
	} else if t.progress <= 0 && t.reverse {
		t.reverse = false
		t.progress = -t.progress
		t.startTime = t.startTime.Add(t.Time)
	}
	return false
}
