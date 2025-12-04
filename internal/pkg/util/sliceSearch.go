package util

import "cmp"

func BinarySearch[T cmp.Ordered](slice []T, target T) int {
	low, high := 0, len(slice)-1
	for low <= high {
		mid := (low + high) / 2
		if slice[mid] == target {
			return mid
		} else if slice[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}
