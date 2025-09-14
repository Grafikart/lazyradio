package utils

func CastSlice[From, To any](slice []From) []To {
	result := make([]To, len(slice))
	for i, item := range slice {
		result[i] = any(item).(To)
	}
	return result
}
