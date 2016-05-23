package utilities

// Finds the index of string in the string array
// returns -1 if not found
func FindIndexStringArr(cols []string, name string) int {
	for i, v := range cols {
		if name == v {
			return i
		}
	}
	return -1
}
