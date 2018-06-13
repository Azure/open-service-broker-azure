package ptr

// ToString returns a pointer to a string. This exists for convenience elsewhere
// since getting the address of a literal isn't possible.
func ToString(val string) *string {
	return &val
}

// ToInt returns a pointer to an int. This exists for convenience elsewhere
// since getting the address of a literal isn't possible.
func ToInt(val int) *int {
	return &val
}

// ToInt32 returns a pointer to an int32. This exists for convenience elsewhere
// since getting the address of a literal isn't possible.
func ToInt32(val int32) *int32 {
	return &val
}

// ToInt64 returns a pointer to an int64. This exists for convenience elsewhere
// since getting the address of a literal isn't possible.
func ToInt64(val int64) *int64 {
	return &val
}

// ToFloat32 returns a pointer to a float32. This exists for convenience
// elsewhere since getting the address of a literal isn't possible.
func ToFloat32(val float32) *float32 {
	return &val
}

// ToFloat64 returns a pointer to a float64. This exists for convenience
// elsewhere since getting the address of a literal isn't possible.
func ToFloat64(val float64) *float64 {
	return &val
}
