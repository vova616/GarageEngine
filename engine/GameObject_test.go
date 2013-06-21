package engine

import (
	"testing"
)

func TestGameObjects(t *testing.T) {
	LoadTestScene()
	a, b, c := NewGameObject("A"), NewGameObject("B"), NewGameObject("C")
	a.AddToScene()
	b.AddToScene()
	c.AddToScene()

	objs := &mainScene.SceneBase().gameObjects

	contains := func(g *GameObject) bool {
		for _, gs := range *objs {
			if gs == g {
				return true
			}
		}
		return false
	}

	if !contains(a) {
		t.Fatal("Object A does not exists")
	}
	if !contains(b) {
		t.Fatal("Object B does not exists")
	}
	if !contains(c) {
		t.Fatal("Object C does not exists")
	}

	if len(*objs) != 3 {
		t.Fatal("Len of scene is", len(*objs), "need 3")
	}

	a.RemoveFromScene()

	if contains(a) {
		t.Fatal("Object A does exists after delete")
	}

	bc := NewGameObject("Bc")
	bc.transform.SetParent2(b)

	if !contains(bc) {
		t.Fatal("Object Bc does not exists")
	}

	ac := NewGameObject("Ac")
	ac.transform.SetParent2(a)

	if contains(ac) {
		t.Fatal("Object Ac does exists when his parent not")
	}

	mainScene.SceneBase().cleanNil()

	if len(*objs) != 3 {
		t.Fatal("Len of scene is", len(*objs), "need 3")
	}

	bc.RemoveFromScene()

	if contains(bc) {
		t.Fatal("Object Bc does exists after delete")
	}
}
