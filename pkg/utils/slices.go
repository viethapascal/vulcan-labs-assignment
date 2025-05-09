package utils

func SliceWithSize(size int32) []int32 {
	r := make([]int32, size)
	var i int32
	for i = 0; i < size; i++ {
		r[i] = i
	}
	return r
}

func FillMap[V any](m map[int32]V, size int32, value V) {
	var i int32
	for i = 0; i < size; i++ {
		m[i] = value
	}
}
