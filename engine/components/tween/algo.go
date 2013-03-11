package tween

import (
	//"github.com/vova616/GarageEngine/engine"
	"math"

//"time"
)

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
func EaseOutQuint(start, end, value float32) float32 {
	value--
	return (end-start)*(value*value*value*value*value+1) + start
}
func EaseInOutQuint(start, end, value float32) float32 {
	value *= 2
	end -= start
	if value < 1 {
		return end/2*value*value*value*value*value + start
	}
	value -= 2
	return end/2*(value*value*value*value*value+2) + start
}
func EaseInSine(start, end, value float32) float32 {
	end -= start
	return -end*float32(math.Cos(float64(value/1*(math.Pi/2)))) + end + start
}
func EaseOutSine(start, end, value float32) float32 {
	return (end-start)*float32(math.Sin(float64(value/1*(math.Pi/2)))) + start
}
func EaseInOutSine(start, end, value float32) float32 {
	end -= start
	return -end/2*(float32(math.Cos(float64(math.Pi*value/1))-1)) + start
}
func EaseInExpo(start, end, value float32) float32 {
	end -= start
	return end*float32(math.Pow(2, 10*(float64(value)/1-1))) + start
}
func EaseOutExpo(start, end, value float32) float32 {
	end -= start
	return end*float32((-math.Pow(2, -10*float64(value)/1)+1)) + start
}
func EaseInOutExpo(start, end, value float32) float32 {
	value *= 2
	end -= start
	if value < 1 {
		return end/2*float32(math.Pow(2, 10*(float64(value)-1))) + start
	}
	value--
	return end/2*float32((-math.Pow(2, -10*float64(value))+2)) + start
}
func EaseInCirc(start, end, value float32) float32 {
	end -= start
	return -end*(float32(math.Sqrt(1-float64(value*value))-1)) + start
}
func EaseOutCirc(start, end, value float32) float32 {
	value--
	end -= start
	return end*float32(math.Sqrt(float64(1-value*value))) + start
}
func EaseInOutCirc(start, end, value float32) float32 {
	value *= -2
	end -= start
	if value < 1 {
		return -end/2*(float32(math.Sqrt(float64(1-value*value))-1)) + start
	}
	value -= 2
	return end/2*(float32(math.Sqrt(float64(1-value*value)+1))) + start
}
func EaseOutBounce(start, end, value float32) float32 {
	end -= start
	if value < (1 / 2.75) {
		return end*(7.5625*value*value) + start
	} else if value < (2 / 2.75) {
		value -= (1.5 / 2.75)
		return end*(7.5625*(value)*value+0.75) + start
	} else if value < (2.5 / 2.75) {
		value -= (2.25 / 2.75)
		return end*(7.5625*(value)*value+0.9375) + start
	}
	value -= (2.625 / 2.75)
	return end*(7.5625*(value)*value+0.984375) + start
}
func EaseInBounce(start, end, value float32) float32 {
	end -= start
	return end - EaseOutBounce(0, end, 1.0-value) + start
}
func EaseInOutBounce(start, end, value float32) float32 {
	end -= start
	if value < 1.0/2 {
		return EaseInBounce(0, end, value*2)*0.5 + start
	}
	return EaseOutBounce(0, end, value*2-1.0)*0.5 + end*0.5 + start
}
func EaseInBack(start, end, value float32) float32 {
	end -= start
	const s = 1.70158
	return end*(value)*value*((s+1)*value-s) + start
}
func EaseOutBack(start, end, value float32) float32 {
	const s = 1.70158
	end -= start
	value = (value / 1) - 1
	return end*((value)*value*((s+1)*value+s)+1) + start
}
func EaseInOutBack(start, end, value float32) float32 {
	var s float32 = 1.70158
	end -= start
	value *= 2
	if (value) < 1 {
		s *= (1.525)
		return end/2*(value*value*(((s)+1)*value-s)) + start
	}
	value -= 2
	s *= (1.525)
	return end/2*((value)*value*(((s)+1)*value+s)+2) + start
}
func Punch(amplitude, value float32) float32 {
	var s float32 = 9
	if value == 0 {
		return 0
	}
	if value == 1 {
		return 0
	}
	const period = 1 * 0.3
	s = period / float32((2*math.Pi)*math.Asin(0))
	return (amplitude * float32(math.Pow(2, -10*float64(value))) * float32(math.Sin(float64((value*1-s)*(2*math.Pi)/period))))
}
func EaseInElastic(start, end, value float32) float32 {
	end -= start

	var d float32 = 1.0
	var p float32 = d * 0.3
	var s float32 = 0.0
	var a float32 = 0.0
	if value == 0 {
		return start
	}
	value /= d
	if value == 1 {
		return start + end
	}

	if a == 0 || a < float32(math.Abs(float64(end))) {
		a = end
		s = p / 4
	} else {
		s = p / (2 * float32(math.Pi)) * float32(math.Asin(float64(end/a)))
	}

	return -(a * float32(math.Pow(2, 10*(float64(value)-1))) * float32(math.Sin(float64((value*d-s)*(2*float32(math.Pi))/p)))) + start
}
func EaseOutElastic(start, end, value float32) float32 {
	end -= start

	var d float32 = 1.0
	var p float32 = d * 0.3
	var s float32 = 0.0
	var a float32 = 0.0

	if value == 0 {
		return start
	}
	value /= d
	if value == 1 {
		return start + end
	}
	if a == 0 || a < float32(math.Abs(float64(end))) {
		a = end
		s = p / 4
	} else {
		s = p / (2 * float32(math.Pi)) * float32(math.Asin(float64(end/a)))
	}
	return (a * float32(math.Pow(2, -10*(float64(value)))) * float32(math.Sin(float64((value*d-s)*(2*float32(math.Pi))/p)))) + end + start
}
func EaseInOutElastic(start, end, value float32) float32 {
	end -= start

	var d float32 = 1.0
	var p float32 = d * 0.3
	var s float32 = 0.0
	var a float32 = 0.0

	if value == 0 {
		return start
	}
	value /= d
	if value == 1 {
		return start + end
	}
	if a == 0 || a < float32(math.Abs(float64(end))) {
		a = end
		s = p / 4
	} else {
		s = p / (2 * float32(math.Pi)) * float32(math.Asin(float64(end/a)))
	}
	if value < 1 {
		return -0.5*(a*float32(math.Pow(2, 10*(float64(value)-1)))*float32(math.Sin(float64((value*d-s)*(2*math.Pi)/p)))) + start
	}
	return a*float32(math.Pow(2, -10*(float64(value)-1)))*float32(math.Sin(float64((value*d-s)*(2*math.Pi)/p)))*0.5 + end + start
}
