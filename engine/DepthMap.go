package engine

import (
	"fmt"
	"sort"
)

type DepthMap []depthData

var depthMap DepthMap

type depthData struct {
	array []*GameObject
	depth int
}

func init() {
	depthMap = make([]depthData, 0, 10)
}

func (this DepthMap) Len() int {
	return len(this)
}

func (this DepthMap) Less(i, j int) bool {
	return this[i].depth < this[j].depth
}

func (this DepthMap) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this *DepthMap) getDepth(depth int, create bool) *depthData {
	arr := *this
	for i := 0; i < len(arr); i++ {
		dData := &arr[i]
		if dData.depth == depth {
			return dData
		} else if dData.depth > depth {
			break
		}
	}
	if !create {
		return nil
	}
	*this = append(*this, depthData{nil, depth})
	sort.Sort(this)
	return this.getDepth(depth, true)
}

func (this *DepthMap) Add(depth int, object *GameObject) {
	dData := this.getDepth(depth, true)
	dData.array = append(dData.array, object)
}

func (this *DepthMap) Remove(depth int, object *GameObject) bool {
	if object != nil {
		dData := this.getDepth(depth, false)
		if dData == nil {
			return false
		}
		arr := dData.array
		for i, ob := range dData.array {
			if ob == object {
				arr[len(arr)-1], arr[i], dData.array = nil, arr[len(arr)-1], arr[:len(arr)-1]
				return true
			}
		}
	}
	return false
}

func (this DepthMap) String() string {
	s := fmt.Sprintln("Depths", len(this))
	for _, dData := range this {
		s += fmt.Sprintln("Depth:", dData.depth, "Items:", len(dData.array))
	}
	return s
}

func (this DepthMap) Iter(fnc func(*GameObject)) {
	for i := 0; i < len(this); i++ {
		dData := &this[i]
		arr := dData.array
		if len(arr) > 0 {
			for j := len(arr) - 1; j >= 0; j-- {
				obj := arr[j]
				if obj != nil {
					fnc(obj)
					//Checks if the drawArr array has been changed
					if obj != arr[j] {
						j++
					}
				}
			}
		}
	}
}
