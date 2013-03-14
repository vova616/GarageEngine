package input

import "testing"

func TestKeyboard(t *testing.T) {
	if KeyDown('B') {
		t.Errorf("Key B is down need up.")
	}
	OnKey('B', key_Press)
	if KeyUp('B') {
		t.Errorf("Key B is up need down.")
	}
	if !KeyPress('B') {
		t.Errorf("Key B isn't pressed need pressed.")
	}
	UpdateInput()
	if KeyPress('B') {
		t.Errorf("Key B is pressed need not pressed.")
	}
}

func TestMouse(t *testing.T) {
	if MouseDown(Mouse1) {
		t.Errorf("Mouse1 is down need up.")
	}
	ButtonPress(Mouse1, key_Press)
	if MouseUp(Mouse1) {
		t.Errorf("Mouse1 is up need down.")
	}
	if !MousePress(Mouse1) {
		t.Errorf("Mouse1 isn't pressed need pressed.")
	}
	UpdateInput()
	if MousePress(Mouse1) {
		t.Errorf("Mouse1 is pressed need not pressed.")
	}
}

func TestText(t *testing.T) {
	test := 'A'
	test2 := 'A'
	key := AddCharCallback(func(char rune) {
		test = char
	})
	key2 := AddCharCallback(func(char rune) {
		test2 = char
	})
	OnChar('B', key_Press)
	if test != 'B' {
		t.Errorf("test need B have %v", string(test))
	}
	if test2 != 'B' {
		t.Errorf("test2 need B have %v", string(test))
	}
	RemoveCharCallback(key)
	OnChar('A', key_Press)
	if test != 'B' {
		t.Errorf("test need B have %v", string(test))
	}
	if test2 != 'A' {
		t.Errorf("test2 need A have %v", string(test))
	}
	RemoveCharCallback(key2)
	OnChar('B', key_Press)
	if test2 != 'A' {
		t.Errorf("test2 need A have %v", string(test))
	}
	if test != 'B' {
		t.Errorf("test need B have %v", string(test))
	}

	test = 'A'
	test2 = 'A'
	key = AddCharCallback(func(char rune) {
		test = char
	})
	key2 = AddCharCallback(func(char rune) {
		test2 = char
	})
	RemoveCharCallback(key)
	RemoveCharCallback(key2)
	OnChar('B', key_Press)
	if test != 'A' {
		t.Errorf("test need A have %v", string(test))
	}
	if test2 != 'A' {
		t.Errorf("test2 need A have %v", string(test))
	}
}
