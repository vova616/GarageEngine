package input

var (
	MousePosition          func() (int, int) = nil
	MouseWheelPosition     func() int        = nil
	MouseWheelDelta        int
	lastMouseWheelPosition int
)

func ButtonPress(btn, state int) {
	switch state {
	case key_Release:
		mouseState[btn] = idle
	case key_Press:
		mouseState[btn] = pressed | wasPressed
	}
}

func MouseWheelCallback(position int) {
	MouseWheelDelta = position - lastMouseWheelPosition
	lastMouseWheelPosition = position
}

func MouseDown(key int) bool {
	return mouseState[key]&pressed != 0
}

func MouseUp(key int) bool {
	return mouseState[key]&pressed == 0
}

func MousePress(key int) bool {
	return mouseState[key]&wasPressed != 0
}
