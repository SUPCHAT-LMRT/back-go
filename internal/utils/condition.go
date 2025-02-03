package utils

func IfThenElse[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}
