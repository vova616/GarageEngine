package engine

import "image"

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

func (node *AtlasNode) Insert(img image.Image, id ID) *AtlasNode {
	if node.Child[0] != nil {
		newNode := node.Child[1].Insert(img, id)
		if newNode != nil {
			newNode.imageRect = img.Bounds()
			newNode.ImageID = id
			return newNode
		}
		newNode = node.Child[0].Insert(img, id)
		if newNode != nil {
			newNode.imageRect = img.Bounds()
			newNode.ImageID = id
			return newNode
		}
	} else {
		if node.ImageID != nil {
			return nil
		}

		dw := node.NodeRect.Dx() - (img.Bounds().Dx() + Padding)
		dh := node.NodeRect.Dy() - (img.Bounds().Dy() + Padding)

		if dw < 0 ||
			dh < 0 {
			return nil
		}

		if dw == 0 && dh == 0 {
			node.imageRect = img.Bounds()
			node.ImageID = id
			return node
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
	return nil
}
