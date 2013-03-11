package tween

import (
//"github.com/vova616/GarageEngine/engine"
//"math"
//"time"
)

func None(t *Tween, progress float32) (newProgress float32, destroy bool) {
	if progress >= 1 {
		progress = 1
		return progress, true
	}
	return progress, false
}

func Loop(t *Tween, progress float32) (newProgress float32, destroy bool) {
	if progress >= 1 {
		progress = t.progress - 1
		t.startTime = t.startTime.Add(t.Time)
	}
	return progress, false
}

func PingPong(t *Tween, progress float32) (newProgress float32, destroy bool) {
	if progress >= 1 && !t.reverse {
		t.reverse = true
		progress = 2 - progress
		t.startTime = t.startTime.Add(t.Time)
	} else if progress <= 0 && t.reverse {
		t.reverse = false
		progress = -t.progress
		t.startTime = t.startTime.Add(t.Time)
	}

	return progress, false
}
