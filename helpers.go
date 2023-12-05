package errtags

func sliceHasCommonSubset[T comparable](sliceA []T, sliceB []T) bool {
	for _, itemB := range sliceB {
		for _, itemA := range sliceA {
			if itemB == itemA {
				return true
			}
		}
	}

	return false
}
