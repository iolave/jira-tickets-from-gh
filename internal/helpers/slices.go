package helpers

func FilterSlice[T any](slice []T, keep func(T) bool) []T {
	var newSlice []T

	for _, item := range slice {
		if keep(item) {
			newSlice = append(newSlice, item)
		}
	}

	return newSlice
}

func MapSlice[T any](slice []T, mapfn func(T) T) []T {
	var newSlice []T

	for _, item := range slice {
		newSlice = append(newSlice, mapfn(item))
	}

	return newSlice
}
