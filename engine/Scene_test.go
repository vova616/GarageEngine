package engine

type TestScene struct {
	*SceneData
}

func (s *TestScene) Load() {

}

func (s *TestScene) New() Scene {
	gs := new(TestScene)
	gs.SceneData = NewScene("TestScene")
	return gs
}

func LoadTestScene() Scene {
	testScene := &TestScene{}
	LoadScene(testScene)
	return mainScene
}
