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

func (s *SceneData) AddGameObject(gameObject ...*GameObject) {
	s.gameObjects = append(s.gameObjects, gameObject...)
}

func (s *SceneData) RemoveGameObject(g *GameObject) {
	for i, c := range s.gameObjects {
		if g == c {
			s.gameObjects = append(s.gameObjects[:i], s.gameObjects[i+1:]...)
			break
		}
	}
}
