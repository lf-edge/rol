package utils

//RemoveElementFromSlice remove element from slice by value
func RemoveElementFromSlice[T comparable](slice []T, value T) []T {
	for index, elem := range slice {
		if elem == value {
			slice[index] = slice[len(slice)-1]
			return slice[:len(slice)-1]
		}
	}
	return slice
}

//SliceContainsElement check what slice contains value
func SliceContainsElement[T comparable](s []T, elem T) bool {
	for _, a := range s {
		if a == elem {
			return true
		}
	}
	return false
}

//SliceDiffElements determine slice diff from origin to modified
//
//Parameters:
//	original - original slice
//	modified - modified slice
//Returns:
//	T[] - deleted elements
//	T[] - added elements
func SliceDiffElements[T comparable](original []T, modified []T) ([]T, []T) {
	deletedElems := []T{}
deleted:
	for _, elemFromOriginal := range original {
		for _, elemFromModified := range modified {
			if elemFromOriginal == elemFromModified {
				continue deleted
			}
		}
		deletedElems = append(deletedElems, elemFromOriginal)
	}
	addedElems := []T{}
added:
	for _, elemFromModified := range modified {
		for _, elemFromOriginal := range original {
			if elemFromOriginal == elemFromModified {
				continue added
			}
		}
		addedElems = append(addedElems, elemFromModified)
	}
	return deletedElems, addedElems
}
