package ptr

// ToInt returns a pointer to an int. This exists for convenience elsewhere
// since getting the address of a literal isn't possible.
func ToInt(val int) *int {
	return &val
}

// ToInt64 returns a pointer to an int64. This exists for convenience elsewhere
// since getting the address of a literal isn't possible.
func ToInt64(val int64) *int64 {
	return &val
}

// ToFloat64 returns a pointer to an float64. This exists for convenience
// elsewhere since getting the address of a literal isn't possible.
func ToFloat64(val float64) *float64 {
	return &val
}
