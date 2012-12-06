package Input

type CharCallback func(char rune)

var (
	keyState   = make(map[int]byte)
	mouseState = make(map[int]byte)

	charCallbacks = []CharCallback{}
)

const (
	idle       = byte(0)
	pressed    = byte(1)
	wasPressed = byte(2)
)

func OnKey(key, state int) {
	switch state {
	case Key_Release:
		keyState[key] = idle
	case Key_Press:
		keyState[key] = pressed | wasPressed
	}
}

func OnChar(key, state int) {
	for _, callback := range charCallbacks {
		callback(rune(key))
	}
}

func AddCharCallback(callback CharCallback) {
	charCallbacks = append(charCallbacks, callback)
}

func UpdateInput() {
	for i, v := range keyState {
		keyState[i] = v & ^wasPressed
	}
	for i, v := range mouseState {
		mouseState[i] = v & ^wasPressed
	}
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
