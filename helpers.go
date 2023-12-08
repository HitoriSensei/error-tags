package errtags

func isSubset[T comparable](sliceA []T, sliceB []T) bool {
	for _, itemB := range sliceB {
		found := false
		for _, itemA := range sliceA {
			if itemB == itemA {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
