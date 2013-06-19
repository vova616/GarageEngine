package engine

import (
	"fmt"
	"reflect"

//"github.com/teomat/mater/collision"
)

type GameObject struct {
	name         string
	transform    *Transform
	components   []Component
	valid        bool
	active       bool
	silentActive bool
	destoryMark  bool

	Tag     string
	Physics *Physics
	Sprite  *Sprite
}

func init() {
	fmt.Print()
}

func NewGameObject(name string) *GameObject {
	g := new(GameObject)
	g.name = name
	g.components = make([]Component, 0)
	g.valid = true
	g.active = true
	g.AddComponent(NewTransform())
	return g
}

func (g *GameObject) Components() []Component {
	arr := make([]Component, len(g.components))
	copy(arr, g.components)
	return arr
}

func (g *GameObject) ComponentTypeOf(component Component) Component {
	typ := reflect.TypeOf(component)
	for _, c := range g.components {
		if typ == reflect.TypeOf(c) {
			return c
		}
	}
	return nil
}

func (g *GameObject) ComponentImplements(intrfce interface{}) Component {
	typ := reflect.TypeOf(intrfce).Elem()
	for _, c := range g.components {
		t := reflect.TypeOf(c)
		if t.Implements(typ) {
			return c
		}
	}
	return nil
}

func (c *GameObject) GameObject() *GameObject {
	if c.IsValid() == false {
		return nil
	}
	return c
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

func (g *GameObject) SetActive(active bool) {
	g.active = active
	g.silentActive = active
	if active {
		for _, c := range g.components {
			c.OnEnable()
		}
	} else {
		for _, c := range g.components {
			c.OnDisable()
		}
	}
}

func (g *GameObject) SetActiveRecursive(active bool) {
	g.SetActive(active)
	for _, c := range g.transform.children {
		c.GameObject().SetActiveRecursive(active)
	}
}

//Used to call OnEnable & OnDisable on object which leave the scene
func (g *GameObject) setActiveRecursiveSilent(active bool) {
	g.setActiveSilent(active)
	for _, c := range g.transform.children {
		c.GameObject().setActiveRecursiveSilent(active)
	}
}

//Used to call OnEnable & OnDisable on object which leave the scene
func (g *GameObject) setActiveSilent(active bool) {
	if !g.active || g.silentActive == active {
		return
	}
	g.silentActive = active
	if active {
		for _, c := range g.components {
			if c != nil {
				c.OnEnable()
			}
		}
	} else {
		for _, c := range g.components {
			if c != nil {
				c.OnDisable()
			}
		}
	}
}

//Removed object from Scene if hes in one
func (g *GameObject) RemoveFromScene() {
	if g.transform.InScene() {
		g.transform.SetParent(nil)
		GetScene().SceneBase().removeGameObject(g)
		g.setActiveRecursiveSilent(false)
		g.transform.removeFromDepthMapRecursive()
	}
}

func (g *GameObject) AddToScene() {
	if !g.transform.childOfScene {
		g.transform.SetParent(nil)
	}
}

func (g *GameObject) IsActive() bool {
	return g.active
}

func (g *GameObject) Destroy() {
	g.destoryMark = true
	g.active = false
	for _, c := range g.transform.children {
		c.gameObject.Destroy()
	}
}

func (g *GameObject) destroy() {
	l := len(g.components)
	for i := l - 1; i >= 0; i-- {
		g.components[i].OnDestroy()
		g.components[i] = nil
	}

	chs := g.transform.children
	l = len(chs)
	for i := l - 1; i >= 0; i-- {
		chs[i].GameObject().destroy()
	}

	if g.transform.childOfScene {
		g.RemoveFromScene()
	} else if g.transform.parent != nil {
		t := g.transform
		for i, c := range t.parent.children {
			if t == c {
				t.parent.children = append(t.parent.children[:i], t.parent.children[i+1:]...)
				break
			}
		}
		g.transform.parent = nil
	}

	g.name = ""
	g.transform = nil
	g.components = nil
	g.valid = false
	g.active = false
	g.Sprite = nil
	g.Physics = nil
}

func (g *GameObject) Clone() *GameObject {
	ng := new(GameObject)
	ng.valid = true
	ng.active = true
	ng.name = g.name + ""
	ng.Tag = g.Tag
	ng.components = make([]Component, 0)

	for _, c := range g.components {
		v := reflect.ValueOf(c).Elem()
		n := reflect.New(v.Type())
		n.Elem().Set(v)
		nc := n.Interface().(Component)
		nc.setGameObject(ng)
		nc.setStarted(false)
		ng.AddComponent(nc)
		nc.Clone()
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
