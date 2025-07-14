package parser

// Ptr is a utility function to create a pointer to a value of type T.
func Ptr[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string](v T) *T {
	return &v
}
