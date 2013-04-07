package engine

import (
	"errors"
	"image"
	"sort"
)

type RectID struct {
	Rect image.Rectangle
	ID   ID
}

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

func FindOptimalSizeFast(totalSize int64) (w, h int) {
	ww, hh := int64(1), int64(1)
	sw := true
	for ww*hh < totalSize {
		if sw {
			hh *= 2
		} else {
			ww *= 2
		}
		sw = !sw
	}
	return int(ww), int(hh)
}

type RectSortable []RectID

func (this RectSortable) Len() int {
	return len(this)
}

func (this RectSortable) Less(i, j int) bool {
	return (this[i].Rect.Dx() * this[i].Rect.Dy()) > (this[j].Rect.Dx() * this[j].Rect.Dy())
}

func (this RectSortable) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

/*
This needs to be smarter, but it does work great for images like fonts
*/
func FindOptimalSize(tries int, rects ...RectID) (w, h int, node *AtlasNode, err error) {
	totalSize := int64(0)
	for _, rect := range rects {
		totalSize += int64(rect.Rect.Dx() * rect.Rect.Dy())
	}

	w, h = FindOptimalSizeFast(totalSize)
	sw := true
	if w < h {
		sw = false
	}

	ww, hh := int64(w), int64(h)

	sort.Sort(RectSortable(rects))

Top:
	for i := 0; i < tries; i++ {
		atlas := NewAtlasNode(int(ww), int(hh))
		for _, rect := range rects {
			_, e := atlas.Insert(rect.Rect, rect.ID)
			if e != nil {
				if sw {
					hh *= 2
				} else {
					ww *= 2
				}
				continue Top
			}
		}
		return int(ww), int(hh), atlas, nil
	}

	return -1, -1, nil, errors.New("Cannot find optimal size")
}
