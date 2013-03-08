package engine

import (
	"errors"
	"image"
)

type AtlasNode struct {
	Child     [2]*AtlasNode
	NodeRect  image.Rectangle
	imageRect image.Rectangle
	ImageID   ID
}

func NewAtlasNode(width, height int) *AtlasNode {
	return &AtlasNode{NodeRect: image.Rect(0, 0, width, height)}
}

func (node *AtlasNode) ImageRect() image.Rectangle {
	return node.imageRect.Add(node.NodeRect.Min)
}

func (node *AtlasNode) NodeByID(id ID) *AtlasNode {
	if node == nil {
		return nil
	}
	if node.ImageID == id {
		return node
	}
	n := node.Child[0].NodeByID(id)
	if n != nil {
		return n
	}
	return node.Child[1].NodeByID(id)
}

func (node *AtlasNode) Insert(img image.Rectangle, id ID) (*AtlasNode, error) {
	if node.Child[0] != nil {
		newNode, _ := node.Child[1].Insert(img, id)
		if newNode != nil {
			newNode.imageRect = img
			newNode.ImageID = id
			return newNode, nil
		}
		newNode, _ = node.Child[0].Insert(img, id)
		if newNode != nil {
			newNode.imageRect = img
			newNode.ImageID = id
			return newNode, nil
		}
	} else {
		if node.ImageID != nil {
			return nil, errors.New("Not enough space in atlas.")
		}

		dw := node.NodeRect.Dx() - (img.Dx() + Padding)
		dh := node.NodeRect.Dy() - (img.Dy() + Padding)

		if dw < 0 || dh < 0 {
			return nil, errors.New("Not enough space in atlas.")
		}

		if dw == 0 && dh == 0 {
			node.imageRect = img
			node.ImageID = id
			return node, nil
		}

		node.Child[0] = &AtlasNode{}
		node.Child[1] = &AtlasNode{}

		if dw > dh {
			node.Child[0].NodeRect = image.Rect(
				node.NodeRect.Min.X, node.NodeRect.Min.Y,
				node.NodeRect.Min.X+dw, node.NodeRect.Max.Y)

			node.Child[1].NodeRect = image.Rect(
				node.NodeRect.Min.X+dw, node.NodeRect.Min.Y,
				node.NodeRect.Max.X, node.NodeRect.Max.Y)
		} else {
			node.Child[0].NodeRect = image.Rect(
				node.NodeRect.Min.X, node.NodeRect.Min.Y,
				node.NodeRect.Max.X, node.NodeRect.Min.Y+dh)

			node.Child[1].NodeRect = image.Rect(
				node.NodeRect.Min.X, node.NodeRect.Min.Y+dh,
				node.NodeRect.Max.X, node.NodeRect.Max.Y)
		}
		return node.Child[1].Insert(img, id)
	}
	return nil, errors.New("Not enough space in atlas.")
}

func (node *AtlasNode) Count() int {
	if node == nil {
		return 0
	}
	if node.ImageID == nil {
		return node.Child[0].Count() + node.Child[1].Count()
	}
	return node.Child[0].Count() + node.Child[1].Count() + 1
}

func (node *AtlasNode) Nodes() []*AtlasNode {
	nodes := make([]*AtlasNode, 0, node.Count())
	node.nodes(&nodes)
	return nodes
}

func (node *AtlasNode) nodes(nodes *[]*AtlasNode) {
	if node == nil {
		return
	}
	if node.Child[0] != nil && node.Child[0].ImageID != nil {
		*nodes = append(*nodes, node.Child[0])
	}
	if node.Child[1] != nil && node.Child[1].ImageID != nil {
		*nodes = append(*nodes, node.Child[1])
	}
	node.Child[0].nodes(nodes)
	node.Child[1].nodes(nodes)
}

func (node *AtlasNode) Rebuild() error {
	nodes := node.Nodes()

	newNode := NewAtlasNode(node.NodeRect.Dx(), node.NodeRect.Dy())

	for {
		biggestArea := 0
		biggestIndex := 0
		found := false
		for i, n := range nodes {
			if n != nil {
				found = true
				area := n.imageRect.Dx() * n.imageRect.Dy()
				if area > biggestArea {
					biggestIndex = i
					biggestArea = area
				}
			}
		}
		if !found {
			break
		}
		n := nodes[biggestIndex]
		nodes[biggestIndex] = nil
		_, e := newNode.Insert(n.imageRect, n.ImageID)
		if e != nil {
			return e
		}
	}

	*node = *newNode
	return nil
}
