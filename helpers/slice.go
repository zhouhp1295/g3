package helpers

func IndexOf[T comparable](source []T, search T) int {
	for i, v := range source {
		if v == search {
			return i
		}
	}
	return -1
}
