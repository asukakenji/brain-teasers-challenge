package lib

// isEven returns whether the parameter, n, is a even number.
func isEven(n int) bool {
	return n&1 == 0
}

// Partition partitions the parameter, arr, as described in the question:
// (1) After partitioning, even integers precede the odd integers in the array.
// (2) The array is operated in-place, with a constant amount of extra space.
// (3) The time complexity is O(n), where n is the length of arr.
func Partition(arr []int) {
	i, j := 0, len(arr)-1
	for {
		for ; i < j && isEven(arr[i]); i++ {
		}
		for ; j > i && !isEven(arr[j]); j-- {
		}
		if i >= j {
			break
		}
		arr[i], arr[j] = arr[j], arr[i]
	}
}
