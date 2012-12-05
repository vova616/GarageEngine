package Input

var (
	MousePosition func() (int, int) = nil
)

func ButtonPress(btn, state int) {
	switch state {
	case Key_Release:
		mouseState[btn] &= 2
	case Key_Press:
		if mouseState[btn] == 0 {
			mouseState[btn] = 3
		} else {
			mouseState[btn] |= 1
		}
	}
}

func MouseDown(key int) bool {
	return mouseState[key]&1 != 0
}

func MouseUp(key int) bool {
	return mouseState[key]&1 == 0
}

func MousePress(key int) bool {
	return mouseState[key]&2 != 0
}


