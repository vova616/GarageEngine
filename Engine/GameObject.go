package Engine

import (
	"fmt"
	"reflect"

//"github.com/teomat/mater/collision"
)

type GameObject struct {
	name        string
	transform   *Transform
	components  []Component
	valid       bool
	active      bool
	destoryMark bool

	Physics *Physics
	Sprite  *Sprite
}

func init() {
	fmt.Print()
}

func NewGameObject(name string) *GameObject {
	g := new(GameObject)
	g.name = name
	g.transform = NewTransform(g)
	g.components = make([]Component, 0)
	g.valid = true
	g.active = true
	return g
}

func (g *GameObject) Components() []Component {
	arr := make([]Component, len(g.components))
	copy(arr, g.components)
	return arr
}

func (g *GameObject) ComponentTypeOf(typ reflect.Type) Component {
	for _, c := range g.components {
		if typ == reflect.TypeOf(c) {
			return c
		}
	}
	return nil
}

func (g *GameObject) SetName(name string) {
	g.name = name
}

func (g *GameObject) Name() string {
	return g.name
}

func (g *GameObject) Transform() *Transform {
	return g.transform
}

func (g *GameObject) IsValid() bool {
	return g.valid
}

func (g *GameObject) SetActive(a bool) {
	g.active = a
}

func (g *GameObject) IsActive() bool {
	return g.active
}

func (g *GameObject) Destroy() {
	g.destoryMark = false
	g.active = false
}

func (g *GameObject) destroy() {
	g.name = ""

	for _, c := range g.transform.children {
		c.GameObject().destroy()
	}
	g.transform = nil
	for _, c := range g.components {
		c.Destroy()
	}
	g.Transform().SetParent(nil)
	g.components = nil
	g.valid = false
	g.active = false
}

func (g *GameObject) Clone() *GameObject {
	ng := new(GameObject)
	ng.valid = true
	ng.active = true
	ng.transform = g.transform.clone(ng)
	ng.name = g.name + ""
	ng.components = make([]Component, 0)
	for _, c := range g.components {
		v := reflect.ValueOf(c).Elem()
		n := reflect.New(v.Type())
		n.Elem().Set(v)
		nc := n.Interface().(Component)
		nc.setGameObject(ng)
		ng.AddComponent(nc)
	}
	return ng
}

func (g *GameObject) AddComponent(com Component) Component {
	com.onAdd(com, g)
	com.setStarted(false)
	g.components = append(g.components, com)
	return com
}

func (g *GameObject) RemoveComponent(com Component) bool {
	t := reflect.TypeOf(com)
	for i, c := range g.components {
		if t == reflect.TypeOf(c) {
			g.components = append(g.components[:i], g.components[i+1:]...)
			return true
		}
	}
	return false
}

func (g *GameObject) RemoveComponentOfType(typ reflect.Type) bool {
	for i, c := range g.components {
		if typ == reflect.TypeOf(c) {
			g.components = append(g.components[:i], g.components[i+1:]...)
			return true
		}
	}
	return false
}

func (g *GameObject) RemoveComponentsOfType(typ reflect.Type) {
	for i, c := range g.components {
		if typ == reflect.TypeOf(c) {
			g.components = append(g.components[:i], g.components[i+1:]...)
		}
	}
}
