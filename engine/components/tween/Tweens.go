package tween

import (
	//"fmt"
	"github.com/vova616/GarageEngine/engine"
	"time"
)

type Tweens struct {
	engine.BaseComponent
	TweensArr []*Tween
}

func (this *Tweens) AddTween(tween *Tween) {
	this.TweensArr = append(this.TweensArr, tween)
}

func (this *Tweens) RemoveTween(tween *Tween) {
	for i, t := range this.TweensArr {
		if t == tween {
			this.TweensArr[i], this.TweensArr = this.TweensArr[len(this.TweensArr)-1], this.TweensArr[:len(this.TweensArr)-1]
			break
		}
	}
}

func (this *Tweens) Update() {
	for _, tween := range this.TweensArr {

		if tween.updateProgress() {
			this.RemoveTween(tween)
			tween.EndCallback()
		}
		//fmt.Println(tween.progress, tween.reverse)

		if tween.Type != nil {
			tween.Value()
		}
	}
}

func newTweens(t *Tween) *Tweens {
	ts := &Tweens{engine.NewComponent(), make([]*Tween, 0, 2)}
	ts.AddTween(t)
	return ts
}

func Create(t *Tween) *Tween {
	if t.To == nil || (t.To != nil && t.From == nil && t.Type == nil) {
		panic("Not possible tween")
	}
	if t.Target != nil {
		t.Target.AddComponent(newTweens(t))
	}
	if t.Algo == nil {
		t.Algo = Linear
	}
	if t.Loop == nil {
		t.Loop = None
	}
	t.startTime = engine.GameTime()
	t.progress = 0
	return t
}

func CreateHelper(target *engine.GameObject, typef TypeFunc, from []float32, to []float32, time time.Duration) *Tween {
	return Create(&Tween{Target: target, Type: typef, From: from, To: to, Time: time})
}

func CreateHelper2(target *engine.GameObject, typef TypeFunc, from []float32, to []float32, time time.Duration, algo Algorithm) *Tween {
	return Create(&Tween{Target: target, Type: typef, From: from, To: to, Time: time, Algo: algo})
}

func CreateHelper3(target *engine.GameObject, typef TypeFunc, from []float32, to []float32, time time.Duration, algo Algorithm, loop LoopFunc) *Tween {
	return Create(&Tween{Target: target, Type: typef, From: from, To: to, Time: time, Algo: algo, Loop: loop})
}

func CreateHelper4(target *engine.GameObject, typef TypeFunc, from []float32, to []float32, time time.Duration, algo Algorithm, loop LoopFunc, format string) *Tween {
	return Create(&Tween{Target: target, Type: typef, From: from, To: to, Time: time, Algo: algo, Loop: loop, Format: format})
}

func CreateHelper5(target *engine.GameObject, typef TypeFunc, from []float32, to []float32, time time.Duration, format string) *Tween {
	return Create(&Tween{Target: target, Type: typef, From: from, To: to, Time: time, Format: format})
}

func CreateHelper6(target *engine.GameObject, typef TypeFunc, from []float32, to []float32, time time.Duration, algo Algorithm, format string) *Tween {
	return Create(&Tween{Target: target, Type: typef, From: from, To: to, Time: time, Algo: algo, Format: format})
}
