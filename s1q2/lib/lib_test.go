package lib

import (
	"math/rand"
	"testing"
	"time"
)

///////////
// Tests //
///////////

// TestDigitCount checks the function digitCount() with predefined test cases.
func TestDigitCount(t *testing.T) {
	cases := []struct {
		n        int
		expected uint
	}{
		{-12, 0},
		{-1, 0},
		{0, 1},
		{1, 1},
		{9, 1},
		{10, 2},
		{11, 2},
		{19, 2},
		{99, 2},
		{100, 3},
		{101, 3},
		{199, 3},
		{999, 3},
		{1234, 4},
	}
	for _, c := range cases {
		got := digitCount(c.n)
		if got != c.expected {
			t.Errorf("DigitCount(%d) = %d, expected %d", c.n, got, c.expected)
		}
	}
}

// TestSplit checks the function split() with predefined test cases.
func TestSplit(t *testing.T) {
	type returnType struct {
		n1 int
		n2 int
		ok bool
	}
	cases := []struct {
		n        int
		expected []returnType
	}{
		{-12, []returnType{{0, 0, false}, {0, 0, false}, {0, 0, false}}},
		{-1, []returnType{{0, 0, false}, {0, 0, false}, {0, 0, false}}},
		{0, []returnType{{0, 0, false}, {0, 0, false}, {0, 0, false}}},
		{1, []returnType{{0, 1, false}, {0, 1, false}, {0, 1, false}}},
		{9, []returnType{{0, 9, false}, {0, 9, false}, {0, 9, false}}},
		{10, []returnType{{1, 0, true}, {0, 10, false}, {0, 10, false}}},
		{11, []returnType{{1, 1, true}, {0, 11, false}, {0, 11, false}}},
		{19, []returnType{{1, 9, true}, {0, 19, false}, {0, 19, false}}},
		{99, []returnType{{9, 9, true}, {0, 99, false}, {0, 99, false}}},
		{100, []returnType{{10, 0, true}, {1, 0, true}, {0, 100, false}}},
		{101, []returnType{{10, 1, true}, {1, 1, true}, {0, 101, false}}},
		{199, []returnType{{19, 9, true}, {1, 99, true}, {0, 199, false}}},
		{999, []returnType{{99, 9, true}, {9, 99, true}, {0, 999, false}}},
		{1234, []returnType{{123, 4, true}, {12, 34, true}, {1, 234, true}, {0, 1234, false}}},
	}
	for _, c := range cases {
		sp := split(c.n)
		for iteration, expected := range c.expected {
			gotN1, gotN2, gotOk := sp()
			if gotN1 != expected.n1 || gotN2 != expected.n2 || gotOk != expected.ok {
				t.Errorf("Iteration %d of split(%d) = (%d, %d, %t), expected (%d, %d, %t)", iteration, c.n, gotN1, gotN2, gotOk, expected.n1, expected.n2, expected.ok)
			}
		}
	}
}

////////////////
// Benchmarks //
////////////////

const (
	numbersLen = 4093
)

var (
	numbers [numbersLen]int
	primes  = [...]int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67}
	jj      int
)

func init() {
	seed := time.Now().UTC().UnixNano()
	src := rand.NewSource(seed)
	rng := rand.New(src)

	// Returns an int in the range [-1000000000, +1000000000].
	randInt := func() int {
		return rng.Intn(2000000001) - 1000000000
	}

	for i, _ := range numbers {
		numbers[i] = randInt()
	}
	jj = primes[rng.Intn(len(primes))]
}

// benchmarkDigitCount is a skeleton for benchmarking the function digitCount().
// BenchmarkDigitCount1-8	100000000	        12.8 ns/op
// BenchmarkDigitCount2-8	100000000	        10.3 ns/op
// BenchmarkDigitCount3-8	100000000	        10.1 ns/op    <- Sharp winner
func benchmarkDigitCount(b *testing.B, digitCountFunc func(int) uint) {
	j := 0
	for i := 0; i < b.N; i++ {
		digitCountFunc(numbers[j])
		j = (j + jj) % numbersLen
	}
}

func BenchmarkDigitCount1(b *testing.B) {
	benchmarkDigitCount(b, func(n int) uint {
		if n < 0 {
			return 0
		}
		if n == 0 {
			return 1
		}
		result := uint(0)
		for n != 0 {
			n /= 10
			result++
		}
		return result
	})
}

func BenchmarkDigitCount2(b *testing.B) {
	benchmarkDigitCount(b, func(n int) uint {
		if n < 0 {
			return 0
		}
		result := uint(1)
		for n2 := 10; n2 <= n; n2 *= 10 {
			result++
		}
		return result
	})
}

func BenchmarkDigitCount3(b *testing.B) {
	benchmarkDigitCount(b, func(n int) uint {
		if n < 0 {
			return 0
		}
		result := uint(1)
		for n2 := 10; n2 <= n; n2, result = n2*10, result+1 {
		}
		return result
	})
}
