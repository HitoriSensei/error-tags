package errtags

import "slices"

func isSubSlice[T comparable](sliceA []T, sliceB []T) bool {
	// This special case is not needed in this implementation as slices are never empty.
	//if len(sliceB) == 0 {
	//	return true
	//}

	if len(sliceB) > len(sliceA) {
		return false
	}

	first := sliceB[0]

	for i, itemA := range sliceA {
		if itemA != first {
			continue
		}

		return slices.Equal(sliceA[i:i+len(sliceB)], sliceB)
	}

	return false
}
