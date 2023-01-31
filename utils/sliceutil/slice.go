package sliceutil

func InSlice[E comparable](needle E, haystack []E) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

func SliceUnique[E comparable](arr []E) []E {
	result := make([]E, 0, len(arr))
	tmpSet := map[E]struct{}{}
	for _, item := range arr {
		if _, ok := tmpSet[item]; !ok {
			tmpSet[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func SplitSliceToGroup[E any](arr []E, size int64) [][]E {
	max := int64(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= size {
		return [][]E{arr}
	}
	//获取应该数组分割为多少份
	var quantity int64
	if max%size == 0 {
		quantity = max / size
	} else {
		quantity = (max / size) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]E, 0)
	//声明分割数组的截止下标
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * size
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * size
	}
	return segments
}
