package engine

type DepthMap map[int8][]*GameObject

var depthMap = make(DepthMap)

func (this *DepthMap) Add(depth int, object *GameObject) int {
	depth8 := int8(depth)
	arr := depthMap[depth8]
	arr = append(arr, object)
	depthMap[depth8] = arr
	return len(arr) - 1
}

func (this *DepthMap) Remove(depth int, index int) bool {
	depth8 := int8(depth)
	if index >= 0 {
		arr := depthMap[depth8]
		if index > len(arr)-1 {
			return false
		}
		arr[len(arr)-1], arr[index], depthMap[depth8] = nil, arr[len(arr)-1], arr[:len(arr)-1]
		return true
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
