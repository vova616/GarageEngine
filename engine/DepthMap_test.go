package engine

import "testing"

func BenchmarkDepthMap(bb *testing.B) {
	LoadTestScene()
	a, b, c := NewGameObject("A"), NewGameObject("B"), NewGameObject("C")
	a.AddToScene()
	b.AddToScene()
	c.AddToScene()
	a.Transform().SetDepth(-1)
	b.Transform().SetDepth(0)
	c.Transform().SetDepth(1)

	for i := 0; i < bb.N; i++ {
		depthMap.Iter(func(g *GameObject) {

		})
	}
}

func TestDepthMap(t *testing.T) {
	LoadTestScene()
	a, b, c := NewGameObject("A"), NewGameObject("B"), NewGameObject("C")
	a.AddToScene()
	b.AddToScene()
	c.AddToScene()
	a.Transform().SetDepth(-1)
	b.Transform().SetDepth(0)
	c.Transform().SetDepth(1)

	depthMap.Iter(func(g *GameObject) {
		if !((g.Name() == "A" && g.Transform().Depth() == -1) ||
			(g.Name() == "B" && g.Transform().Depth() == 0) ||
			(g.Name() == "C" && g.Transform().Depth() == 1)) {
			t.Errorf("Unkown gameobject %s %d", g.Name(), g.Transform().Depth())
		}
	})

	b.Transform().SetDepth(1)

	depthMap.Iter(func(g *GameObject) {
		if !((g.Name() == "A" && g.Transform().Depth() == -1) ||
			(g.Name() == "B" && g.Transform().Depth() == 1) ||
			(g.Name() == "C" && g.Transform().Depth() == 1)) {
			t.Errorf("Unkown gameobject %s %d", g.Name(), g.Transform().Depth())
		}
	})

	b.RemoveFromScene()

	depthMap.Iter(func(g *GameObject) {
		if !((g.Name() == "A" && g.Transform().Depth() == -1) ||
			(g.Name() != "B") ||
			(g.Name() == "C" && g.Transform().Depth() == 1)) {
			t.Errorf("Unkown gameobject %s %d", g.Name(), g.Transform().Depth())
		}
	})

}
