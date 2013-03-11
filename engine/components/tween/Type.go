package tween

import (
	"github.com/vova616/GarageEngine/engine"

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
func Color(t *Tween, arr []float32) []float32 {
	col := t.Target.Sprite.Color
	if t.Target.Sprite == nil {
		panic("Cannot run Color tween on none Sprite GameObjects")
	}
	if arr == nil || len(arr) == 0 {
		return []float32{col.R, col.G, col.B, col.A}
	}
	col = ColorFmt(col, arr, t.Format)
	t.Target.Sprite.Color = col
	return []float32{col.R, col.G, col.B, col.A}
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

func ColorFmt(v engine.Color, arr []float32, s string) engine.Color {
	if len(s) == 0 {
		if len(arr) > 3 {
			v.R = arr[0]
			v.G = arr[1]
			v.B = arr[2]
			v.A = arr[3]
		} else if len(arr) > 2 {
			v.R = arr[0]
			v.G = arr[1]
			v.B = arr[2]
		} else if len(arr) > 1 {
			v.R = arr[0]
			v.G = arr[1]
		} else {
			v.R = arr[0]
		}
		return v
	}
	if (len(s) <= 4 && len(s) >= 1) && len(arr) == 1 {
		for _, r := range s {
			switch r {
			case 'r', 'R':
				v.R = arr[0]
			case 'g', 'G':
				v.G = arr[0]
			case 'b', 'B':
				v.B = arr[0]
			case 'a', 'A':
				v.A = arr[0]
			}
		}
		return v
	}
	for i, r := range s {
		if i >= len(arr) {
			break
		}
		switch r {
		case 'r', 'R':
			v.R = arr[i]
		case 'g', 'G':
			v.G = arr[i]
		case 'b', 'B':
			v.B = arr[i]
		case 'a', 'A':
			v.A = arr[i]
		}
	}
	return v
}

func VectorFmt(v engine.Vector, arr []float32, s string) engine.Vector {
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

func VectorFmtRotation(v engine.Vector, arr []float32, s string) engine.Vector {
	if len(s) == 0 && len(arr) == 1 {
		v.Z = arr[0]
		return v
	}
	return VectorFmt(v, arr, s)
}
