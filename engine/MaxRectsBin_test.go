package engine

import (
	"image"
	"testing"
)

func TestMaxRectsBin(t *testing.T) {
	b := NewBin(10, 10, 1)
	_, e := b.Insert(image.Rect(0, 0, 5, 5))
	if e != nil {
		t.Error(e)
		return
	}
	_, e = b.Insert(image.Rect(0, 0, 2, 2))
	if e != nil {
		t.Error(e)
		return
	}
	_, e = b.Insert(image.Rect(0, 0, 2, 2))
	if e != nil {
		t.Error(e)
		return
	}
	_, e = b.Insert(image.Rect(0, 0, 2, 2))
	if e != nil {
		t.Error(e)
		return
	}
	_, e = b.Insert(image.Rect(0, 0, 1, 1))
	if e != nil {
		t.Error(e)
		return
	}
	_, e = b.Insert(image.Rect(0, 0, 1, 1))
	if e != nil {
		t.Error(e)
		return
	}
	_, e = b.Insert(image.Rect(0, 0, 1, 1))
	if e != nil {
		t.Error(e)
		return
	}
	_, e = b.Insert(image.Rect(0, 0, 1, 1))
	if e != nil {
		t.Error(e)
		return
	}
	_, e = b.Insert(image.Rect(0, 0, 1, 1))
	if e != nil {
		t.Error(e)
		return
	}
	_, e = b.Insert(image.Rect(0, 0, 1, 1))
	if e != nil {
		t.Error(e)
		return
	}
	_, e = b.Insert(image.Rect(0, 0, 1, 1))
	if e == nil {
		t.Error("Should return error got nil")
		return
	}
}
