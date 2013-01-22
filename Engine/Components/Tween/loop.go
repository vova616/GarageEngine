package Tween

import (
//"github.com/vova616/GarageEngine/Engine"
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
<<<<<<< HEAD
func Loop(t *Tween) bool {
	if t.progress >= 1 {
		t.progress = t.progress - 1
=======

func Loop(t *Tween, progress float32) (newProgress float32, destroy bool) {
	if progress >= 1 {
		progress = t.progress - 1
>>>>>>> upstream/master
		t.startTime = t.startTime.Add(t.Time)
	}
	return progress, false
}

func PingPong(t *Tween, progress float32) (newProgress float32, destroy bool) {
	if progress >= 1 && !t.reverse {
		t.reverse = true
<<<<<<< HEAD
		t.progress = 2 - t.progress
		t.startTime = t.startTime.Add(t.Time)
	} else if t.progress <= 0 && t.reverse {
		t.reverse = false
		t.progress = -t.progress
=======
		progress = 2 - progress
		t.startTime = t.startTime.Add(t.Time)
	} else if progress <= 0 && t.reverse {
		t.reverse = false
		progress = -t.progress
>>>>>>> upstream/master
		t.startTime = t.startTime.Add(t.Time)
	}

	return progress, false
}
