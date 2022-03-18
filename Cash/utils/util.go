package utils

func Clamp(x, min, max int) int {
	if x > max {
		x = max
	} else if x < min {
		x = min
	}
	return x
}
