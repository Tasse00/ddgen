package utils

func ContainsString(arr []string, val string) bool {
	for _, elem := range arr {
		if elem == val {
			return true
		}
	}
	return false
}
