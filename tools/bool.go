package tools

// BoolBinary will return 1 if the value is true, 0 if the value is false.
func BoolBinary(val bool) uint8 {
	if val {
		return 1
	}
	return 0
}
