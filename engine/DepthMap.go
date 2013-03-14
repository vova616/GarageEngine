package engine

type DepthMap map[int8][]*GameObject

var depthMap = make(DepthMap)

func (this *DepthMap) Add(depth int, object *GameObject) {
	depth8 := int8(depth)
	arr := depthMap[depth8]
	arr = append(arr, object)
	depthMap[depth8] = arr
}

func (this *DepthMap) Remove(depth int, object *GameObject) bool {
	depth8 := int8(depth)
	if object != nil {
		arr := depthMap[depth8]
		for i, ob := range arr {
			if ob == object {
				arr[len(arr)-1], arr[i], depthMap[depth8] = nil, arr[len(arr)-1], arr[:len(arr)-1]
				return true
			}
		}
	}
	return false
}

func (this *DepthMap) Iter(fnc func(*GameObject)) {
	for i := int8(-127); ; i++ {
		drawArr, exists := depthMap[i]
		if exists && len(drawArr) > 0 {
			for j := len(drawArr) - 1; j >= 0; j-- {
				obj := drawArr[j]
				if obj != nil {
					fnc(obj)
					//Checks if the drawArr array has been changed
					if obj != drawArr[j] {
						j++
					}
				}
			}
		}
		if i == 127 {
			break
		}
	}
}
