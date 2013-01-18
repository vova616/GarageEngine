package Tween

import (
	"github.com/vova616/GarageEngine/Engine"
	"math"
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
		t.startTime = Engine.GameTime().Add(time.Duration(float64(t.progress) * float64(t.Time)))
	}
	return false
}
func PingPong(t *Tween) bool {
	if t.progress >= 1 && !t.reverse {
		t.reverse = true
		t.progress = t.progress - 1
		t.startTime = Engine.GameTime().Add(time.Duration(float64(t.progress) * float64(t.Time)))
	} else if t.progress <= 0 && t.reverse {
		t.reverse = false
		t.progress = -t.progress
		t.startTime = Engine.GameTime().Add(time.Duration(float64(t.progress) * float64(t.Time)))
	}
	return false
}
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

func Linear(start, end, value float32) float32 {
	return (end-start)*value + start
}
func Clerp(start, end, value float32) float32 {
	var max, half, retval, diff float32 = 360.0, 180.0, 0.0, 0.0
	if (end - start) < -half {
		diff = ((max - start) + end) * value
		retval = start + diff
	} else if (end - start) > half {
		diff = -((max - end) + start) * value
		retval = start + diff
	} else {
		retval = start + (end-start)*value
	}
	return retval
}
func EaseInQuad(start, end, value float32) float32 {
	return (end-start)*value*value + start
}
func EaseOutQuad(start, end, value float32) float32 {
	return (end-start)*value*(value-2) + start
}
func EaseInOutQuad(start, end, value float32) float32 {
	value *= 2
	end -= start
	if value < 1 {
		return end/2*value*value + start
	}
	value--
	return -end/2*(value*(value-2)-1) + start
}
func EaseInCubic(start, end, value float32) float32 {
	return (end-start)*value*value*value + start
}
func EaseOutCubic(start, end, value float32) float32 {
	value--
	return (end-start)*(value*value*value+1) + start
}
func EaseInOutCubic(start, end, value float32) float32 {
	value *= 2
	end -= start
	if value < 1 {
		return end/2*value*value*value + start
	}
	value -= 2
	return end/2*(value*value*value+2) + start
}
func EaseInQuart(start, end, value float32) float32 {
	return (end-start)*value*value*value*value + start
}
func EeaseOutQuart(start, end, value float32) float32 {
	value--
	return -(end-start)*(value*value*value*value-1) + start
}
func EaseOutInQuart(start, end, value float32) float32 {
	value *= 2
	end -= start
	if value < 1 {
		return end/2*value*value*value*value + start
	}
	return -end/2*(value*value*value*value-2) + start
}

func Spring(start, end, value float32) float32 {
	if value > 1 {
		value = 1
	} else if value < 0 {
		value = 0
	}
	v := float64(value)
	v = (math.Sin(v*math.Pi*(0.2+2.5*v*v*v))*math.Pow(1.0-v, 2.2) + v) * (1.0 + (1.2 * (1.0 - v)))
	value = float32(v)
	return start + (end-start)*value
}
func EaseInQuint(start, end, value float32) float32 {
	return (end-start)*value*value*value*value*value + start
}

/* 13, 22

	private float easeOutQuint(float start, float end, float value){
		value--;
		end -= start;
		return end * (value * value * value * value * value + 1) + start;
	}

	private float easeInOutQuint(float start, float end, float value){
		value /= .5f;
		end -= start;
		if (value < 1) return end / 2 * value * value * value * value * value + start;
		value -= 2;
		return end / 2 * (value * value * value * value * value + 2) + start;
	}

	private float easeInSine(float start, float end, float value){
		end -= start;
		return -end * Mathf.Cos(value / 1 * (Mathf.PI / 2)) + end + start;
	}

	private float easeOutSine(float start, float end, float value){
		end -= start;
		return end * Mathf.Sin(value / 1 * (Mathf.PI / 2)) + start;
	}

	private float easeInOutSine(float start, float end, float value){
		end -= start;
		return -end / 2 * (Mathf.Cos(Mathf.PI * value / 1) - 1) + start;
	}

	private float easeInExpo(float start, float end, float value){
		end -= start;
		return end * Mathf.Pow(2, 10 * (value / 1 - 1)) + start;
	}

	private float easeOutExpo(float start, float end, float value){
		end -= start;
		return end * (-Mathf.Pow(2, -10 * value / 1) + 1) + start;
	}

	private float easeInOutExpo(float start, float end, float value){
		value /= .5f;
		end -= start;
		if (value < 1) return end / 2 * Mathf.Pow(2, 10 * (value - 1)) + start;
		value--;
		return end / 2 * (-Mathf.Pow(2, -10 * value) + 2) + start;
	}

	private float easeInCirc(float start, float end, float value){
		end -= start;
		return -end * (Mathf.Sqrt(1 - value * value) - 1) + start;
	}

	private float easeOutCirc(float start, float end, float value){
		value--;
		end -= start;
		return end * Mathf.Sqrt(1 - value * value) + start;
	}

	private float easeInOutCirc(float start, float end, float value){
		value /= .5f;
		end -= start;
		if (value < 1) return -end / 2 * (Mathf.Sqrt(1 - value * value) - 1) + start;
		value -= 2;
		return end / 2 * (Mathf.Sqrt(1 - value * value) + 1) + start;
	}

	private float easeInBounce(float start, float end, float value){
		end -= start;
		float d = 1f;
		return end - easeOutBounce(0, end, d-value) + start;
	}

	//private float bounce(float start, float end, float value){
	private float easeOutBounce(float start, float end, float value){
		value /= 1f;
		end -= start;
		if (value < (1 / 2.75f)){
			return end * (7.5625f * value * value) + start;
		}else if (value < (2 / 2.75f)){
			value -= (1.5f / 2.75f);
			return end * (7.5625f * (value) * value + .75f) + start;
		}else if (value < (2.5 / 2.75)){
			value -= (2.25f / 2.75f);
			return end * (7.5625f * (value) * value + .9375f) + start;
		}else{
			value -= (2.625f / 2.75f);
			return end * (7.5625f * (value) * value + .984375f) + start;
		}
	}

	private float easeInOutBounce(float start, float end, float value){
		end -= start;
		float d = 1f;
		if (value < d/2) return easeInBounce(0, end, value*2) * 0.5f + start;
		else return easeOutBounce(0, end, value*2-d) * 0.5f + end*0.5f + start;
	}

	private float easeInBack(float start, float end, float value){
		end -= start;
		value /= 1;
		float s = 1.70158f;
		return end * (value) * value * ((s + 1) * value - s) + start;
	}

	private float easeOutBack(float start, float end, float value){
		float s = 1.70158f;
		end -= start;
		value = (value / 1) - 1;
		return end * ((value) * value * ((s + 1) * value + s) + 1) + start;
	}

	private float easeInOutBack(float start, float end, float value){
		float s = 1.70158f;
		end -= start;
		value /= .5f;
		if ((value) < 1){
			s *= (1.525f);
			return end / 2 * (value * value * (((s) + 1) * value - s)) + start;
		}
		value -= 2;
		s *= (1.525f);
		return end / 2 * ((value) * value * (((s) + 1) * value + s) + 2) + start;
	}

	private float punch(float amplitude, float value){
		float s = 9;
		if (value == 0){
			return 0;
		}
		if (value == 1){
			return 0;
		}
		float period = 1 * 0.3f;
		s = period / (2 * Mathf.PI) * Mathf.Asin(0);
		return (amplitude * Mathf.Pow(2, -10 * value) * Mathf.Sin((value * 1 - s) * (2 * Mathf.PI) / period));
    }

	private float easeInElastic(float start, float end, float value){
		end -= start;

		float d = 1f;
		float p = d * .3f;
		float s = 0;
		float a = 0;

		if (value == 0) return start;

		if ((value /= d) == 1) return start + end;

		if (a == 0f || a < Mathf.Abs(end)){
			a = end;
			s = p / 4;
			}else{
			s = p / (2 * Mathf.PI) * Mathf.Asin(end / a);
		}

		return -(a * Mathf.Pow(2, 10 * (value-=1)) * Mathf.Sin((value * d - s) * (2 * Mathf.PI) / p)) + start;
	}		

	private float easeOutElastic(float start, float end, float value){
		end -= start;

		float d = 1f;
		float p = d * .3f;
		float s = 0;
		float a = 0;

		if (value == 0) return start;

		if ((value /= d) == 1) return start + end;

		if (a == 0f || a < Mathf.Abs(end)){
			a = end;
			s = p / 4;
			}else{
			s = p / (2 * Mathf.PI) * Mathf.Asin(end / a);
		}

		return (a * Mathf.Pow(2, -10 * value) * Mathf.Sin((value * d - s) * (2 * Mathf.PI) / p) + end + start);
	}		

	private float easeInOutElastic(float start, float end, float value){
		end -= start;

		float d = 1f;
		float p = d * .3f;
		float s = 0;
		float a = 0;

		if (value == 0) return start;

		if ((value /= d/2) == 2) return start + end;

		if (a == 0f || a < Mathf.Abs(end)){
			a = end;
			s = p / 4;
			}else{
			s = p / (2 * Mathf.PI) * Mathf.Asin(end / a);
		}

		if (value < 1) return -0.5f * (a * Mathf.Pow(2, 10 * (value-=1)) * Mathf.Sin((value * d - s) * (2 * Mathf.PI) / p)) + start;
		return a * Mathf.Pow(2, -10 * (value-=1)) * Mathf.Sin((value * d - s) * (2 * Mathf.PI) / p) * 0.5f + end + start;
	}		*/
