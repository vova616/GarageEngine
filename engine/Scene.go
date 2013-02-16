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
	s.gameObjects = append(s.gameObjects, gameObject...)
}

func (s *SceneData) AddGameObject(gameObject ...*GameObject) {
	for _, obj := range gameObject {
		obj.AddToScene()
	}
}

func (s *SceneData) removeGameObject(g *GameObject) {
	if g == nil {
		return
	}
	for i, c := range s.gameObjects {
		if g == c {
			s.gameObjects[i] = nil
			s.gameObjects = s.gameObjects[:i+copy(s.gameObjects[i:], s.gameObjects[i+1:])]
			break
		}
	}
}

func (s *SceneData) RemoveGameObject(g *GameObject) {
	if g == nil {
		return
	}
	g.RemoveFromScene()
}
