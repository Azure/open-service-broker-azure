package slice

// ContainsString returns a bool indicating whether the given []string contained
// the given string
func ContainsString(slice []string, value string) bool {
	for _, v := range slice {
		if value == v {
			return true
		}
	}
	return false
}

// ContainsInt returns a bool indicating whether the given []int contained the
// given int
func ContainsInt(slice []int, value int) bool {
	for _, v := range slice {
		if value == v {
			return true
		}
	}
	return false
}

// ContainsInt64 returns a bool indicating whether the given []int64 contained
// the given int64
func ContainsInt64(slice []int64, value int64) bool {
	for _, v := range slice {
		if value == v {
			return true
		}
	}
	return false
}
