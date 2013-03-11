package tween

import (
	"github.com/vova616/GarageEngine/engine"
	"time"
)

type TypeFunc func(*Tween, []float32) []float32
type Algorithm func(float32, float32, float32) float32
type LoopFunc func(t *Tween, progress float32) (newProgress float32, destroy bool)

type Tween struct {
	Target *engine.GameObject
	From   []float32
	To     []float32
	Time   time.Duration
	Algo   Algorithm
	Loop   LoopFunc

	StartCallback func()
	EndCallback   func()
	Update        func([]float32)

	startTime time.Time
	progress  float32
	Type      TypeFunc

	reverse bool
	Format  string
}

func (this *Tween) SetFunc(typeFunc TypeFunc) {
	this.Type = typeFunc
}

func (this *Tween) Progress() float32 {
	return this.progress
}

func (t *Tween) updateProgress() bool {
	delta := engine.GameTime().Sub(t.startTime)
	if t.reverse {
		t.progress = 1 - float32(float64(delta)/float64(t.Time))
	} else {
		t.progress = float32(float64(delta) / float64(t.Time))
	}
	destroy := false
	t.progress, destroy = t.Loop(t, t.progress)
	return destroy
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
