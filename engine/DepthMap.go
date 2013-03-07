package engine

var depthMap = make(map[int8][]*GameObject)

func depthMapAdd(depth int8, object *GameObject) (index int) {
	arr := depthMap[depth]
	arr = append(arr, object)
	depthMap[depth] = arr
	return len(arr) - 1
}

func depthMapRemove(depth int8, index int) (deleted bool) {
	if index >= 0 {
		arr := depthMap[depth]
		if index > len(arr)-1 {
			return false
		}
		arr[len(arr)-1], arr[index], depthMap[depth] = nil, arr[len(arr)-1], arr[:len(arr)-1]
		return true
	}
	return false
}
