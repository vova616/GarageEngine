package engine

import (
	"image"
	"testing"
)

func TestAtlasTree(t *testing.T) {
	tree := NewAtlasNode(20, 20)
	n, e := tree.Insert(image.Rect(0, 0, 2, 2), 1)
	if e != nil {
		t.Fatal(e)
	}
	if n != tree.NodeByID(1) {
		t.Fatal("NodeByID returns wrong node", n)
	}
	n, e = tree.Insert(image.Rect(0, 0, 8, 8), 2)
	if e != nil {
		t.Fatal(e)
	}
	if n != tree.NodeByID(2) {
		t.Fatal("NodeByID returns wrong node", n)
	}
	n, e = tree.Insert(image.Rect(0, 0, 3, 3), 3)
	if e != nil {
		t.Fatal(e)
	}
	if n != tree.NodeByID(3) {
		t.Fatal("NodeByID returns wrong node", n)
	}
	n, e = tree.Insert(image.Rect(0, 0, 5, 5), 4)
	if e != nil {
		t.Fatal(e)
	}
	if n != tree.NodeByID(4) {
		t.Fatal("NodeByID returns wrong node", n)
	}
	if tree.Count() != 4 {
		t.Fatal("Tree count is wrong Result:", tree.Count(), "need 4")
	}
	nodes := tree.Nodes()
	if len(nodes) != tree.Count() {
		t.Fatal("Tree Nodes count is wrong Result:", len(nodes), "need 4")
	}
	for _, node := range nodes {
		if node == nil {
			t.Fatal("Tree Nodes functions returns nil node")
		}
		if node.ImageID != 1 && node.ImageID != 2 && node.ImageID != 3 && node.ImageID != 4 {
			t.Fatal("Tree Nodes functions returns wrong nodes", node.ImageID)
		}
	}
	e = tree.Rebuild()
	if e != nil {
		t.Fatal("Rebuild error", e)
	}
	if tree.NodeByID(1) == nil {
		t.Fatal("NodeByID returns wrong node nil")
	}
	if tree.NodeByID(2) == nil {
		t.Fatal("NodeByID returns wrong node nil")
	}
	if tree.NodeByID(3) == nil {
		t.Fatal("NodeByID returns wrong node nil")
	}
	if tree.NodeByID(4) == nil {
		t.Fatal("NodeByID returns wrong node nil")
	}
	if tree.Count() != 4 {
		t.Fatal("Tree count is wrong Result:", tree.Count(), "need 4")
	}
	for _, node := range nodes {
		if node == nil {
			t.Fatal("Tree Nodes functions returns nil node")
		}
		if node.ImageID != 1 && node.ImageID != 2 && node.ImageID != 3 && node.ImageID != 4 {
			t.Fatal("Tree Nodes functions returns wrong nodes", node.ImageID)
		}
	}
}
