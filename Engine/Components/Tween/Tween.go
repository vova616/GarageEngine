package Tween

import (
	"github.com/vova616/GarageEngine/Engine"
	"time"
)

type TweenType byte
type TypeFunc func(*Tween, []float32) []float32
type LoopType func(*Tween) bool

const (
	FromTo = TweenType(0)
	To     = TweenType(1)
)

type Tween struct {
	Target *Engine.GameObject
	From   []float32
	To     []float32
	Time   time.Duration
	Algo   func(float32, float32, float32) float32
	LoopF  LoopType

	StartCallback func()
	EndCallback   func()
	Update        func([]float32)

	startTime time.Time
	progress  float32
	Type      TypeFunc

	reverse bool
}

func (this *Tween) SetFunc(typeFunc TypeFunc) {
	this.Type = typeFunc
}

func (t *Tween) updateProgress() bool {
	delta := Engine.GameTime().Sub(t.startTime)
	if t.reverse {
		t.progress = 1 - float32(float64(delta)/float64(t.Time))
	} else {
		t.progress = float32(float64(delta) / float64(t.Time))
	}
	return t.LoopF(t)
}

func (t *Tween) Value() []float32 {
	if t.From != nil && t.To != nil {
		m := make([]float32, len(t.From))
		for i, _ := range t.From {
			m[i] = t.Algo(t.From[i], t.To[i], t.progress)
		}
		if t.Type != nil {
			t.Type(t, m)
		}
		return m
	} else if t.From == nil && t.To != nil && t.Type != nil {
		from := t.Type(t, nil)
		m := make([]float32, len(from))
		for i, _ := range t.From {
			m[i] = t.Algo(t.From[i], t.To[i], t.progress)
		}
		t.Type(t, m)
		return m
	} else {
		panic("Cannot tween")
	}
	return nil
}
