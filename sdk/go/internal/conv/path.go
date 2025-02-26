package conv

func PathToInt32Slice(path []int) []int32 {
	result := make([]int32, len(path))
	for i, v := range path {
		result[i] = int32(v)
	}
	return result
}
