package engine

import (
	"fmt"
	"reflect"

//"github.com/teomat/mater/collision"
)

type GameObject struct {
	name       string
	transform  *Transform
	components []Component
	valid      bool

	active          bool
	selfActive      bool
	componentActive bool

	destoryMark bool

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
	g.selfActive = true
	g.componentActive = true
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
	if active == g.active {
		return
	}

	if active {
		if g.transform.parent == nil || (g.transform.parent != nil && g.transform.parent.gameObject.active) {
			g.active = true
		}
	} else {
		g.active = false
	}

	if g.selfActive != g.active {
		g.selfActive = g.active
		g.componentActive = g.active
		if g.componentActive {
			for _, c := range g.components {
				c.OnEnable()
			}
		} else {
			for _, c := range g.components {
				c.OnDisable()
			}
		}
	}

	for _, t := range g.transform.children {
		t.gameObject.setChildrenActive(active)
	}
}

func (g *GameObject) setChildrenActive(active bool) {
	if active == g.active {
		return
	}

	if active && g.selfActive {
		g.active = true
	} else {
		g.active = false
	}

	if g.componentActive != g.active {
		g.componentActive = g.active
		if g.componentActive {
			for _, c := range g.components {
				c.OnEnable()
			}
		} else {
			for _, c := range g.components {
				c.OnDisable()
			}
		}
	}

	for _, t := range g.transform.children {
		t.gameObject.setChildrenActive(active)
	}
}

func (g *GameObject) silentActive(active bool) {
	if !g.IsValid() {
		return
	}
	if g.selfActive != active {
		g.selfActive = active
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

	for _, t := range g.transform.children {
		t.gameObject.silentActive(active)
	}
}

//Removed object from Scene if hes in one
func (g *GameObject) RemoveFromScene() {
	if g.transform.InScene() {
		g.transform.SetParent(nil)
		GetScene().SceneBase().removeGameObject(g)
		g.silentActive(false)
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

func (g *GameObject) IsSelfActive() bool {
	return g.selfActive
}

func (g *GameObject) Destroy() {
	g.destoryMark = true
	g.active = false
	for _, c := range g.transform.children {
		c.gameObject.Destroy()
	}
}

func (g *GameObject) destroy() {
	//Remove this gameobject from his parent children.
	//RemoveFromScene is doing it internally.
	if g.transform.childOfScene {
		g.RemoveFromScene()
	} else if g.transform.parent != nil {
		g.transform.removeFromParent()
	}

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
