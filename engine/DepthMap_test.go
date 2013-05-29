package engine

import (
	"testing"
)

const minDepth = -50
const maxDepth = 50
const Objects = 500

func BenchmarkDepthMap_Iter(bb *testing.B) {
	bb.StopTimer()
	LoadTestScene()

	for i, j := 0, 0; i < Objects; i++ {
		j %= maxDepth - minDepth
		a := NewGameObject("A")
		a.AddToScene()
		a.Transform().SetDepth(j + minDepth)
		j++
	}

	bb.StartTimer()
	for i := 0; i < bb.N; i++ {
		depthMap.Iter(func(g *GameObject) {

		})
	}
}

func BenchmarkDepthMap_getDepth(bb *testing.B) {
	bb.StopTimer()
	LoadTestScene()

	for i, j := 0, 0; i < Objects; i++ {
		j %= maxDepth - minDepth
		a := NewGameObject("A")
		a.AddToScene()
		a.Transform().SetDepth(j + minDepth)
		j++
	}

	bb.StartTimer()
	for i, j := 0, 0; i < bb.N; i++ {
		j %= maxDepth - minDepth
		depthMap.getDepth(j+minDepth, false)
		j++
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

	i := -1
	depthMap.Iter(func(g *GameObject) {
		if g.Transform().Depth() != i {
			t.Errorf("bad depth %s %d %d", g.Name(), g.Transform().Depth(), i)
		}
		i++
	})

	d := depthMap.getDepth(-1, false)
	if d == nil || len(d.array) != 1 || d.depth != -1 {
		t.Errorf("bad children or depth need (1,-1) have ", d.depth, len(d.array))
	}
	d = depthMap.getDepth(0, false)
	if d == nil || len(d.array) != 1 || d.depth != 0 {
		t.Errorf("bad children or depth need (1,0) have ", d.depth, len(d.array))
	}
	d = depthMap.getDepth(1, false)
	if d == nil || len(d.array) != 1 || d.depth != 1 {
		t.Errorf("bad children or depth need (1,1) have ", d.depth, len(d.array))
	}

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
