package Tween

import (
	"github.com/vova616/GarageEngine/Engine"
	//"time"
)

type Tweens struct {
	Engine.BaseComponent
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
		}

		if tween.Type != nil {
			tween.Value()
		}
	}
}

func newTweens(t *Tween) *Tweens {
	ts := &Tweens{Engine.NewComponent(), make([]*Tween, 0, 2)}
	ts.AddTween(t)
	return ts
}

func CreateTween(t *Tween) *Tween {
	if t.To == nil || (t.To != nil && t.From == nil && t.Type == nil) {
		panic("Not possible tween")
	}
	if t.Target != nil {
		t.Target.AddComponent(newTweens(t))
	}
	t.startTime = Engine.GameTime()
	t.progress = 0
	return t
}
