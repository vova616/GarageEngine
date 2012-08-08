package Input

import (
    "github.com/jteeuwen/glfw"
	//"log"
)

func ButtonPress(btn,state int) {
	switch state {
		case glfw.KeyRelease:
			mouseState[btn] &= 2
		case glfw.KeyPress:
			if mouseState[btn] == 0 {
				mouseState[btn] = 3
			} else {
				mouseState[btn] |= 1
			}
	} 
}

func MouseDown(key int) bool {
	return mouseState[key] & 1 != 0
}

func MouseUp(key int) bool {
	return mouseState[key] & 1 == 0
}

func MousePress(key int) bool {
	return mouseState[key] & 2 != 0
}


func MousePosition() (int,int) {
	return glfw.MousePos()
}