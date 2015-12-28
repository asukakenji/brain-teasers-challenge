package lib

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

///////////
// Tests //
///////////

const (
	maxUint = ^uint(0)
	minUint = 0
	maxInt  = int(maxUint >> 1)
	minInt  = -maxInt - 1
)

const (
	numbersLen = 4093
)

var (
	numbers      [numbersLen]int
	numbersSlice = numbers[:]
)

func init() {
	seed := time.Now().UTC().UnixNano()
	src := rand.NewSource(seed)
	rng := rand.New(src)

	// Returns an int in the range [-9, +9].
	randInt := func() int {
		return rng.Intn(19) - 9
	}

	for i, _ := range numbers {
		numbers[i] = randInt()
	}
}

// TestIsEven checks the function isEven() with predefined test cases.
func TestIsEven(t *testing.T) {
	cases := []struct {
		n        int
		expected bool
	}{
		{minInt, true},
		{-3, false},
		{-2, true},
		{-1, false},
		{0, true},
		{1, false},
		{2, true},
		{3, false},
		{maxInt, false},
	}
	for _, c := range cases {
		got := isEven(c.n)
		if got != c.expected {
			t.Errorf("isEven(%d) = %t", c.n, got)
		}
	}
}

// TestPartitionPredefined checks the function Partition() with predefined test cases.
func TestPartition(t *testing.T) {
	cases := []struct {
		arr      []int
		expected []int
	}{
		{[]int{}, []int{}},
		{[]int{1}, []int{1}},
		{[]int{2}, []int{2}},
		{[]int{1, 2}, []int{2, 1}},
		{[]int{2, 1}, []int{2, 1}},
		{[]int{1, 3}, []int{1, 3}},
		{[]int{3, 1}, []int{3, 1}},
		{[]int{2, 4}, []int{2, 4}},
		{[]int{4, 2}, []int{4, 2}},
		{[]int{7, 7, 4, 0, 9, 8, 2, 4, 1, 9}, []int{4, 2, 4, 0, 8, 9, 7, 7, 1, 9}},
		{[]int{-4, -3, -2, -1, 0, 1, 2, 3, 4}, []int{-4, 4, -2, 2, 0, 1, -1, 3, -3}},
	}
	for _, c := range cases {
		got := make([]int, len(c.arr))
		copy(got, c.arr)
		Partition(got)
		if !reflect.DeepEqual(got, c.expected) {
			t.Errorf("Partition(%v) = %v, expected %v", c.arr, got, c.expected)
		}
	}
}

// TestPartitionPredefined checks the function Partition() with random test cases.
func TestPartitionRandom(t *testing.T) {
	isPartitioned := func(arr []int) bool {
		oddStarted := false
		for _, n := range arr {
			if !oddStarted {
				if !isEven(n) {
					oddStarted = true
				}
			} else {
				if isEven(n) {
					return false
				}
			}
		}
		return true
	}

	sort := func(arr []int) []int {
		arrCopy := make([]int, len(arr))
		copy(arrCopy, arr)
		sort.Ints(arrCopy)
		return arrCopy
	}

	haveSameElements := func(arr1, arr2 []int) bool {
		arr1Sorted := sort(arr1)
		arr2Sorted := sort(arr2)
		return reflect.DeepEqual(arr1Sorted, arr2Sorted)
	}

	got := make([]int, numbersLen)
	copy(got, numbersSlice)
	Partition(got)
	if !isPartitioned(got) || !haveSameElements(got, numbersSlice) {
		t.Errorf("Partition(%v) = %v, unexpected", numbersSlice, got)
	}
}

////////////////
// Benchmarks //
////////////////

const (
	numberssLen      = 8
	numberssItemsLen = 509
	numberssInterval = 512
)

var (
	numberss [numberssLen][]int
	primes   = [...]int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67}
	jj       int
)

func init() {
	seed := time.Now().UTC().UnixNano()
	src := rand.NewSource(seed)
	rng := rand.New(src)

	for i, _ := range numberss {
		low := numberssInterval * i
		high := low + numberssItemsLen
		numberss[i] = numbers[low:high:high]
	}
	jj = primes[rng.Intn(len(primes))]
}

// benchmarkIsEven is a skeleton for benchmarking the function isEven().
// BenchmarkIsEven1-8   	300000000	         5.16 ns/op
// BenchmarkIsEven2-8   	300000000	         4.75 ns/op    <- Sharp winner
func benchmarkIsEven(b *testing.B, isEvenFunc func(int) bool) {
	j := 0
	for i := 0; i < b.N; i++ {
		isEvenFunc(numbers[j])
		j = (j + jj) % numbersLen
	}
}

func BenchmarkIsEven1(b *testing.B) {
	benchmarkIsEven(b, func(n int) bool {
		return n%2 == 0
	})
}

func BenchmarkIsEven2(b *testing.B) {
	benchmarkIsEven(b, func(n int) bool {
		return n&1 == 0
	})
}

// benchmarkPartition is a skeleton for benchmarking the function Partition().
// BenchmarkPartition1-8	  500000	      2816 ns/op    <- Sharp loser
// BenchmarkPartition2-8	 1000000	      1918 ns/op
// BenchmarkPartition3-8	 1000000	      1957 ns/op
// BenchmarkPartition4-8	 1000000	      1920 ns/op
// BenchmarkPartition5-8	 1000000	      1927 ns/op
// BenchmarkPartition6-8	 1000000	      1908 ns/op
func benchmarkPartition(b *testing.B, partitionFunc func([]int) []int) {
	j := 0
	for i := 0; i < b.N; i++ {
		numberssItemCopy := make([]int, numberssItemsLen)
		copy(numberssItemCopy, numberss[j])
		partitionFunc(numberssItemCopy)
		j = (j + jj) % numberssLen
	}
}

func BenchmarkPartition1(b *testing.B) {
	benchmarkPartition(b, func(arr []int) []int {
		n := len(arr)
		i, j := 0, n-1
		for {
			for ; i < n; i++ {
				if !isEven(arr[i]) {
					break
				}
			}
			for ; j >= 0; j-- {
				if isEven(arr[j]) {
					break
				}
			}
			if i >= j {
				break
			}
			arr[i], arr[j] = arr[j], arr[i]
		}
		return arr
	})
}

func BenchmarkPartition2(b *testing.B) {
	benchmarkPartition(b, func(arr []int) []int {
		n := len(arr)
		i, j := 0, n-1
		for {
			for ; i < n && isEven(arr[i]); i++ {
			}
			for ; j >= 0 && !isEven(arr[j]); j-- {
			}
			if i >= j {
				break
			}
			arr[i], arr[j] = arr[j], arr[i]
		}
		return arr
	})
}

func BenchmarkPartition3(b *testing.B) {
	benchmarkPartition(b, func(arr []int) []int {
		n := len(arr)
		for i, j := 0, n-1; ; {
			for ; i < n && isEven(arr[i]); i++ {
			}
			for ; j >= 0 && !isEven(arr[j]); j-- {
			}
			if i >= j {
				break
			}
			arr[i], arr[j] = arr[j], arr[i]
		}
		return arr
	})
}

func BenchmarkPartition4(b *testing.B) {
	benchmarkPartition(b, func(arr []int) []int {
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
		return arr
	})
}

func BenchmarkPartition5(b *testing.B) {
	benchmarkPartition(b, func(arr []int) []int {
		for i, j := 0, len(arr)-1; ; {
			for ; i < j && isEven(arr[i]); i++ {
			}
			for ; j > i && !isEven(arr[j]); j-- {
			}
			if i >= j {
				break
			}
			arr[i], arr[j] = arr[j], arr[i]
		}
		return arr
	})
}

func BenchmarkPartition6(b *testing.B) {
	benchmarkPartition(b, func(arr []int) []int {
		for i, j := 0, len(arr)-1; ; {
			for ; i < j && isEven(arr[i]); i++ {
			}
			for ; j > i && !isEven(arr[j]); j-- {
			}
			if i >= j {
				break
			}
			temp := arr[j]
			arr[j] = arr[i]
			arr[i] = temp
		}
		return arr
	})
}
