package utils

func PositivePair(x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}
	return true
}
func IntAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func PointsWithinManhattan(x, y, d int, cols int32) []int32 {
	result := make([]int32, 0)
	for dx := -d; dx <= d; dx++ {
		rem := d - IntAbs(dx)
		for dy := -rem; dy <= rem; dy++ {
			tx := x + dx
			ty := y + dy
			if PositivePair(tx, ty) {
				tp := tx*int(cols) + ty
				result = append(result, int32(tp))
				//fmt.Printf("Point: (%d,%d) - I: %d\n", tx, ty, tp)
			}

		}
	}
	return result
}
