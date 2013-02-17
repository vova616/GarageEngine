package engine

var depthMap = make(map[int8][]*GameObject)

func depthMapAdd(depth int8, object *GameObject) (index int) {
	arr := depthMap[depth]
	for i, obj := range arr {
		if obj == nil {
			index = i
			arr[i] = object
			return
		}
	}
	arr = append(arr, object)
	depthMap[depth] = arr
	return len(arr) - 1
}

func depthMapRemove(depth int8, index int) {
	if index >= 0 {
		arr := depthMap[depth]
		if index < len(arr) {
			arr[index] = nil
		}
	}
}
