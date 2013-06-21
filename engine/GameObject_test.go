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

type ActiveComponent struct {
	BaseComponent
	Enabled int
}

func (this *ActiveComponent) OnEnable() {
	this.Enabled++
}

func (this *ActiveComponent) OnDisable() {
	this.Enabled--
}

func TestGameObjects_Activity(t *testing.T) {
	LoadTestScene()
	a, b, c := NewGameObject("A"), NewGameObject("B"), NewGameObject("C")
	aa := &ActiveComponent{NewComponent(), 1}
	ba := &ActiveComponent{NewComponent(), 1}
	ca := &ActiveComponent{NewComponent(), 1}
	a.AddComponent(aa)
	b.AddComponent(ba)
	c.AddComponent(ca)
	b.SetActive(false)
	a.AddToScene()
	b.AddToScene()
	c.AddToScene()

	c.SetActive(false)

	if !a.IsActive() {
		t.Fatal("Object A is not active when it should be")
	}

	if b.IsActive() {
		t.Fatal("Object B is active when it shoudln't be")
	}

	if c.IsActive() {
		t.Fatal("Object C is active when it shoudln't be")
	}

	if ca.Enabled != 0 {
		t.Fatalf("Object C is component enable status is %d need 0", ca.Enabled)
	}

	if ba.Enabled != 0 {
		t.Fatalf("Object B is component enable status is %d need 0", ba.Enabled)
	}

	if aa.Enabled != 1 {
		t.Fatalf("Object A is component enable status is %d need 0", aa.Enabled)
	}

	b.SetActive(true)

	if !b.IsActive() {
		t.Fatal("Object B is not active when it should be")
	}

	c.Transform().SetParent2(a)
	b.Transform().SetParent2(a)

	if c.IsActive() {
		t.Fatal("Object C is active when it shoudln't be")
	}

	if !b.IsActive() {
		t.Fatal("Object B is not active when it should be")
	}

	a.SetActive(false)

	if c.IsActive() {
		t.Fatal("Object C is active when it shoudln't be")
	}

	if b.IsActive() {
		t.Fatal("Object B is active when it shoudln't be")
	}

	if ca.Enabled != 0 {
		t.Fatalf("Object C is component enable status is %d need 0", ca.Enabled)
	}

	if ba.Enabled != 0 {
		t.Fatalf("Object B is component enable status is %d need 0", ba.Enabled)
	}

	if aa.Enabled != 0 {
		t.Fatalf("Object A is component enable status is %d need 0", aa.Enabled)
	}

	a.SetActive(true)

	if c.IsActive() {
		t.Fatal("Object C is active when it shoudln't be")
	}

	if !b.IsActive() {
		t.Fatal("Object B is not active when it should be")
	}

	if ca.Enabled != 0 {
		t.Fatalf("Object C is component enable status is %d need 0", ca.Enabled)
	}

	if ba.Enabled != 1 {
		t.Fatalf("Object B is component enable status is %d need 1", ba.Enabled)
	}

	if aa.Enabled != 1 {
		t.Fatalf("Object A is component enable status is %d need 1", aa.Enabled)
	}

	a.SetActive(false)

	if c.IsActive() {
		t.Fatal("Object C is active when it shoudln't be")
	}

	if b.IsActive() {
		t.Fatal("Object B is active when it shoudln't be")
	}

	if ca.Enabled != 0 {
		t.Fatalf("Object C is component enable status is %d need 0", ca.Enabled)
	}

	if ba.Enabled != 0 {
		t.Fatalf("Object B is component enable status is %d need 0", ba.Enabled)
	}

	if aa.Enabled != 0 {
		t.Fatalf("Object A is component enable status is %d need 0", aa.Enabled)
	}
}
