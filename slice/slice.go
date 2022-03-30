package slice

// Include slice or array include an item
func Include(length int, fn func(i int) bool) bool {
	for i := 0; i < length; i++ {
		if fn(i) {
			return true
		}
	}
	return false
}

// FindIndex The method returns the index of the first element in the array
// that satisfies the provided testing function. Otherwise,
// it returns -1, indicating that no element passed the test.
func FindIndex(length int, fn func(i int) bool) int {
	for i := 0; i < length; i++ {
		if fn(i) {
			return i
		}
	}
	return -1
}
