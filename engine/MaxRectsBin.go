/*
Taken from https://github.com/juj/RectangleBinPack.
*/
package engine

import (
	"errors"
	"image"
)

const maxInt = int(^uint(0) >> 1)

type MaxRectsBin struct {
	Width, Height  int
	Padding        int
	usedRectangles []image.Rectangle
	freeRectangles []image.Rectangle
}

func NewBin(width, height, padding int) *MaxRectsBin {
	return &MaxRectsBin{
		Width:          width,
		Height:         height,
		Padding:        padding,
		usedRectangles: make([]image.Rectangle, 0, 1),
		freeRectangles: []image.Rectangle{image.Rect(0, 0, width, height)},
	}
}

func (this *MaxRectsBin) Insert(rect image.Rectangle) (image.Rectangle, error) {
	r, _, _ := this.FindPositionForNewNodeBestShortSideFit(rect.Dx()+this.Padding, rect.Dy()+this.Padding)
	if r.Dx() == 0 {
		return r, errors.New("Not enough space in atlas.")
	}

	this.placeRect(r)

	r.Max.X -= this.Padding
	r.Max.Y -= this.Padding

	return r, nil
}

func (this *MaxRectsBin) InsertArray(rects []image.Rectangle) ([]image.Rectangle, error) {
	r := make([]image.Rectangle, len(rects))
	numRects := len(rects)
	for numRects != 0 {
		bestScore1 := maxInt
		bestScore2 := maxInt
		bestRectIndex := -1
		var bestNode image.Rectangle

		for i, rect := range rects {
			if r[i] != image.ZR {
				continue
			}
			newNode, score1, score2 := this.FindPositionForNewNodeBestShortSideFit(rect.Dx()+this.Padding, rect.Dy()+this.Padding)

			if score1 < bestScore1 || (score1 == bestScore1 && score2 < bestScore2) {
				bestScore1 = score1
				bestScore2 = score2
				bestNode = newNode
				bestRectIndex = i
			}
		}

		if bestRectIndex == -1 {
			return nil, errors.New("Not enough space in atlas.")
		} else {
			this.placeRect(bestNode)
			bestNode.Max.X -= this.Padding
			bestNode.Max.Y -= this.Padding
			r[bestRectIndex] = bestNode
			numRects--
		}
	}
	return r, nil
}

func (this *MaxRectsBin) placeRect(r image.Rectangle) {
	l := len(this.freeRectangles)
	for i := 0; i < l; i++ {
		if this.SplitFreeNode(this.freeRectangles[i], r) {
			this.freeRectangles = append(this.freeRectangles[:i], this.freeRectangles[i+1:]...)
			i--
			l--
		}
	}
	this.PruneFreeList()
	this.usedRectangles = append(this.usedRectangles, r)
}

/// Computes the ratio of used surface area.
func (this *MaxRectsBin) Occupancy() float32 {
	usedSurfaceArea := uint64(0)
	for _, rect := range this.usedRectangles {
		usedSurfaceArea += uint64(rect.Dx()) * uint64(rect.Dy())
	}

	return float32(float64(usedSurfaceArea) / float64(this.Width*this.Height))
}

func (this *MaxRectsBin) PruneFreeList() {
	/*
		///  Would be nice to do something like this, to avoid a Theta(n^2) loop through each pair.
		///  But unfortunately it doesn't quite cut it, since we also want to detect containment.
		///  Perhaps there's another way to do this faster than Theta(n^2).

		if (freeRectangles.size() > 0)
			clb::sort::QuickSort(&freeRectangles[0], freeRectangles.size(), NodeSortCmp);

		for(size_t i = 0; i < freeRectangles.size()-1; ++i)
			if (freeRectangles[i].x == freeRectangles[i+1].x &&
			    freeRectangles[i].y == freeRectangles[i+1].y &&
			    freeRectangles[i].width == freeRectangles[i+1].width &&
			    freeRectangles[i].height == freeRectangles[i+1].height)
			{
				freeRectangles.erase(freeRectangles.begin() + i);
				--i;
			}
	*/

	/// Go through each pair and remove any rectangle that is redundant.
	for i := 0; i < len(this.freeRectangles); i++ {
		for j := i + 1; j < len(this.freeRectangles); j++ {

			if this.freeRectangles[i].In(this.freeRectangles[j]) {
				this.freeRectangles = append(this.freeRectangles[:i], this.freeRectangles[i+1:]...)
				i--
				break
			}
			if this.freeRectangles[j].In(this.freeRectangles[i]) {
				this.freeRectangles = append(this.freeRectangles[:j], this.freeRectangles[j+1:]...)
				j--
			}
		}
	}
}

func (this *MaxRectsBin) SplitFreeNode(freeNode, usedNode image.Rectangle) bool {
	// Test with SAT if the rectangles even intersect.
	if usedNode.Min.X >= freeNode.Max.X || usedNode.Max.X <= freeNode.Min.X ||
		usedNode.Min.Y >= freeNode.Max.Y || usedNode.Max.Y <= freeNode.Min.Y {
		return false
	}

	if usedNode.Min.X < freeNode.Max.X && usedNode.Max.X > freeNode.Min.X {
		// New node at the top side of the used node.
		if usedNode.Min.Y > freeNode.Min.Y && usedNode.Min.Y < freeNode.Max.Y {
			newNode := freeNode
			newNode.Max.Y = usedNode.Min.Y
			this.freeRectangles = append(this.freeRectangles, newNode)
		}

		// New node at the bottom side of the used node.
		if usedNode.Max.Y < freeNode.Max.Y {
			newNode := freeNode
			newNode.Min.Y = usedNode.Max.Y
			this.freeRectangles = append(this.freeRectangles, newNode)
		}
	}

	if usedNode.Min.Y < freeNode.Max.Y && usedNode.Max.Y > freeNode.Min.Y {
		// New node at the left side of the used node.
		if usedNode.Min.X > freeNode.Min.X && usedNode.Min.X < freeNode.Max.X {
			newNode := freeNode
			newNode.Max.X = usedNode.Min.X
			this.freeRectangles = append(this.freeRectangles, newNode)
		}

		// New node at the right side of the used node.
		if usedNode.Max.X < freeNode.Max.X {
			newNode := freeNode
			newNode.Min.X = usedNode.Max.X
			this.freeRectangles = append(this.freeRectangles, newNode)
		}
	}

	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (this *MaxRectsBin) String() string {
	d := make([]rune, (this.Width+1)*this.Height)

	for x := 0; x < this.Width; x++ {
		for y := 0; y < this.Height; y++ {
			d[y*(this.Width+1)+x] = '0'
		}
	}

	for y := 0; y < this.Height; y++ {
		d[(y*(this.Width+1))+this.Width] = '\n'
	}

	for _, rect := range this.usedRectangles {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			for y := rect.Min.Y; y < rect.Max.Y; y++ {
				d[y*(this.Width+1)+x] = '#'
			}
		}
	}

	return string(d)
}

func (this *MaxRectsBin) FindPositionForNewNodeBestShortSideFit(width, height int) (bestNode image.Rectangle, bestShortSideFit, bestLongSideFit int) {
	bestShortSideFit = maxInt

	for _, r := range this.freeRectangles {
		rW := r.Dx()
		rH := r.Dy()
		// Try to place the rectangle in upright (non-flipped) orientation.
		if rW >= width && rH >= height {
			leftoverHoriz := abs(rW - width)
			leftoverVert := abs(rH - height)
			shortSideFit := min(leftoverHoriz, leftoverVert)
			longSideFit := max(leftoverHoriz, leftoverVert)

			if shortSideFit < bestShortSideFit || (shortSideFit == bestShortSideFit && longSideFit < bestLongSideFit) {
				bestNode.Min = r.Min
				bestNode.Max.X = bestNode.Min.X + width
				bestNode.Max.Y = bestNode.Min.Y + height
				bestShortSideFit = shortSideFit
				bestLongSideFit = longSideFit
			}
		}

		/* Disable rotation
		if (rW >= height && rH >= width)
		{
			flippedLeftoverHoriz := abs(rW - height);
			flippedLeftoverVert := abs(rH - width);
			flippedShortSideFit := min(flippedLeftoverHoriz, flippedLeftoverVert);
			flippedLongSideFit := max(flippedLeftoverHoriz, flippedLeftoverVert);

			if (flippedShortSideFit < bestShortSideFit || (flippedShortSideFit == bestShortSideFit && flippedLongSideFit < bestLongSideFit)) {
				bestNode.x = freeRectangles[i].x;
				bestNode.y = freeRectangles[i].y;
				bestNode.width = height;
				bestNode.height = width;
				bestShortSideFit = flippedShortSideFit;
				bestLongSideFit = flippedLongSideFit;
			}
		}
		*/
	}
	return
}
