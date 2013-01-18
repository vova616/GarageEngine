package Tween

import (
	"github.com/vova616/GarageEngine/Engine"

//"math"
//"time"
)

func Scale(t *Tween, arr []float32) []float32 {
	scale := t.Target.Transform().Scale()
	if arr == nil || len(arr) == 0 {
		return []float32{scale.X, scale.Y, scale.Z}
	}
	scale = VectorFmt(scale, arr, t.Format)
	t.Target.Transform().SetScale(scale)
	return []float32{scale.X, scale.Y, scale.Z}
}

func Position(t *Tween, arr []float32) []float32 {
	pos := t.Target.Transform().Position()
	if arr == nil || len(arr) == 0 {
		return []float32{pos.X, pos.Y, pos.Z}
	}
	pos = VectorFmt(pos, arr, t.Format)
	t.Target.Transform().SetPosition(pos)
	return []float32{pos.X, pos.Y, pos.Z}
}

func Rotation(t *Tween, arr []float32) []float32 {
	rot := t.Target.Transform().Rotation()
	if arr == nil || len(arr) == 0 {
		return []float32{rot.X, rot.Y, rot.Z}
	}
	rot = VectorFmtRotation(rot, arr, t.Format)
	t.Target.Transform().SetRotation(rot)
	return []float32{rot.X, rot.Y, rot.Z}
}

func WorldScale(t *Tween, arr []float32) []float32 {
	scale := t.Target.Transform().WorldScale()
	if arr == nil || len(arr) == 0 {
		return []float32{scale.X, scale.Y, scale.Z}
	}
	scale = VectorFmt(scale, arr, t.Format)
	t.Target.Transform().SetWorldScale(scale)
	return []float32{scale.X, scale.Y, scale.Z}
}

func WorldPosition(t *Tween, arr []float32) []float32 {
	pos := t.Target.Transform().WorldPosition()
	if arr == nil || len(arr) == 0 {
		return []float32{pos.X, pos.Y, pos.Z}
	}
	pos = VectorFmt(pos, arr, t.Format)
	t.Target.Transform().SetWorldPosition(pos)
	return []float32{pos.X, pos.Y, pos.Z}
}

func WorldRotation(t *Tween, arr []float32) []float32 {
	rot := t.Target.Transform().WorldRotation()
	if arr == nil || len(arr) == 0 {
		return []float32{rot.X, rot.Y, rot.Z}
	}
	rot = VectorFmtRotation(rot, arr, t.Format)
	t.Target.Transform().SetWorldRotation(rot)
	return []float32{rot.X, rot.Y, rot.Z}
}

func VectorFmt(v Engine.Vector, arr []float32, s string) Engine.Vector {
	if len(s) == 0 {
		if len(arr) > 2 {
			v.X = arr[0]
			v.Y = arr[1]
			v.Z = arr[2]
		} else if len(arr) > 1 {
			v.X = arr[0]
			v.Y = arr[1]
		} else {
			v.X = arr[0]
		}
		return v
	}
	if (len(s) <= 3 && len(s) >= 1) && len(arr) == 1 {
		for _, r := range s {
			switch r {
			case 'x', 'X':
				v.X = arr[0]
			case 'y', 'Y':
				v.Y = arr[0]
			case 'z', 'Z':
				v.Z = arr[0]
			}
		}
		return v
	}
	for i, r := range s {
		if i >= len(arr) {
			break
		}
		switch r {
		case 'x', 'X':
			v.X = arr[i]
		case 'y', 'Y':
			v.Y = arr[i]
		case 'z', 'Z':
			v.Z = arr[i]
		}
	}
	return v
}

func VectorFmtRotation(v Engine.Vector, arr []float32, s string) Engine.Vector {
	if len(s) == 0 && len(arr) == 1 {
		v.Z = arr[0]
		return v
	}
	return VectorFmt(v, arr, s)
}
