package engine

//import "log"

type SceneData struct {
	name        string
	gameObjects []*GameObject
	Camera      *Camera
}

type Scene interface {
	New() Scene
	Load()
	SceneBase() *SceneData
}

func NewScene(name string) *SceneData {
	return &SceneData{name: name, gameObjects: make([]*GameObject, 0)}
}

func (s *SceneData) SetCamera(Camera *Camera) {
	s.Camera = Camera
}

func (s *SceneData) Name() string {
	return s.name
}

func (s *SceneData) SceneBase() *SceneData {
	return s
}

func (s *SceneData) addGameObject(gameObject ...*GameObject) {
	for _, obj := range gameObject {
		s.gameObjects = append(s.gameObjects, obj)
		obj.transform.childOfScene = true
		for _, t := range obj.transform.children {
			s.addGameObject(t.gameObject)
		}
	}
}

func (s *SceneData) AddGameObject(gameObject ...*GameObject) {
	if s == GetScene().SceneBase() {
		for _, obj := range gameObject {
			obj.AddToScene()
		}
	} else {
		s.addGameObject(gameObject...)
	}
}

func (s *SceneData) removeGameObject(g *GameObject) {
	if g == nil {
		return
	}
	for i, c := range s.gameObjects {
		if g == c {
			s.gameObjects[i].transform.childOfScene = false
			for _, t := range g.transform.children {
				s.removeGameObject(t.gameObject)
			}
			s.gameObjects[i] = nil
			break
		}
	}
}

func (s *SceneData) cleanNil() {
	for i := 0; i < len(s.gameObjects); i++ {
		if s.gameObjects[i] == nil {
			s.gameObjects[i], s.gameObjects = s.gameObjects[len(s.gameObjects)-1], s.gameObjects[:len(s.gameObjects)-1]
			i--
		}
	}
}

func (s *SceneData) RemoveGameObject(g *GameObject) {
	if g == nil {
		return
	}
	if s == GetScene().SceneBase() {
		g.RemoveFromScene()
	} else {
		g.transform.removeFromParent()
		s.removeGameObject(g)
	}
}
