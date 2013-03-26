package input

type CharCallback func(char rune)
type ChatCallbackKey *CharCallback

var (
	keyState   = make(map[int]byte)
	mouseState = make(map[int]byte)

	charCallbacks = []*CharCallback{}
)

const (
	idle       = byte(0)
	pressed    = byte(1)
	wasPressed = byte(2)
)

func OnKey(key, state int) {
	switch state {
	case key_Release:
		keyState[key] = idle
	case key_Press:
		keyState[key] = pressed | wasPressed
	}
}

func OnChar(key, state int) {
	for i, callback := range charCallbacks {
		if callback != nil && *callback != nil {
			(*callback)(rune(key))
		} else {
			charCallbacks[len(charCallbacks)-1], charCallbacks[i], charCallbacks = nil, charCallbacks[len(charCallbacks)-1], charCallbacks[:len(charCallbacks)-1]
			if callback != nil {
				*callback = nil
			}
		}
	}
}

func AddCharCallback(callback CharCallback) (key ChatCallbackKey) {
	if callback == nil {
		return
	}
	c := &callback
	charCallbacks = append(charCallbacks, c)
	return c
}

func RemoveCharCallback(key ChatCallbackKey) (deleted bool) {
	if key == nil || *key == nil {
		return false
	}
	for i, c := range charCallbacks {
		if c == key {
			charCallbacks[len(charCallbacks)-1], charCallbacks[i], charCallbacks = nil, charCallbacks[len(charCallbacks)-1], charCallbacks[:len(charCallbacks)-1]
			*key = nil
			return true
		}
	}
	return false
}

func UpdateInput() {
	for i, v := range keyState {
		keyState[i] = v & ^wasPressed
	}
	for i, v := range mouseState {
		mouseState[i] = v & ^wasPressed
	}
	MouseWheelDelta = 0
}

func ClearInput() {
	for i, _ := range keyState {
		keyState[i] = idle
	}
	for i, _ := range mouseState {
		mouseState[i] = idle
	}
}

func KeyDown(key int) bool {
	return keyState[key]&pressed != 0
}

func KeyUp(key int) bool {
	return keyState[key]&pressed == 0
}

func KeyPress(key int) bool {
	return keyState[key]&wasPressed != 0
}
