package Engine

//import "log"

type SceneData struct {
	name	string
	gameObjects []*GameObject
}

type Scene interface {
	New() Scene
	Load()
	SceneBase() *SceneData
} 

func NewScene(name string) *SceneData{
	return &SceneData{name, make([]*GameObject, 0)}
}

func (s *SceneData) Name() string {
	return s.name
}

func (s *SceneData) AddGameObject(g *GameObject) {
	s.gameObjects = append(s.gameObjects, g)
}

func (s *SceneData) RemoveGameObject(g *GameObject) {
	for i,c := range s.gameObjects {
		if g == c {
			s.gameObjects = append(s.gameObjects[:i], s.gameObjects[i+1:]...)
			break
		}
	}
}

