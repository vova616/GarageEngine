package Tween

import (
//"github.com/vova616/GarageEngine/Engine"
//"math"
//"time"
)

func Scale(t *Tween, arr []float32) []float32 {
	scale := t.Target.Transform().Scale()
	if arr == nil || len(arr) == 0 {
		return []float32{scale.X, scale.Y, scale.Z}
	}
	if len(arr) > 2 {
		scale.X = arr[0]
		scale.Y = arr[1]
		scale.Z = arr[2]
	} else if len(arr) > 1 {
		scale.X = arr[0]
		scale.Y = arr[1]
	} else {
		scale.X = arr[0]
	}
	t.Target.Transform().SetScale(scale)
	return []float32{scale.X, scale.Y, scale.Z}
}
