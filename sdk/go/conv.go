package sourcetool

func convertPathToInt32Slice(p path) []int32 {
	result := make([]int32, len(p))
	for i, v := range p {
		result[i] = int32(v)
	}
	return result
}
