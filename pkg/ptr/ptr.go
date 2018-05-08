package ptr

func ToInt(val int) *int {
	return &val
}

func ToInt64(val int64) *int64 {
	return &val
}

func ToFloat64(val float64) *float64 {
	return &val
}
