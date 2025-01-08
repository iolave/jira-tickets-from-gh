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
