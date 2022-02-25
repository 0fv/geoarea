package geoarea

func euqal(x float64, y float64) bool {
	v := x - y
	const delta float64 = 1e-8
	if v < delta && v > -delta {
		return true
	}
	return false

}

func little(x float64, y float64) bool {
	if euqal(x, y) {
		return false
	}
	return x < y
}
func little_equal(x float64, y float64) bool {
	if euqal(x, y) {
		return true
	}
	return x < y
}
