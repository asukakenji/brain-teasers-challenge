package lib

import "fmt"

// The Expression struct contains the textual representation of an expression
// and the value of it.
// "Text" is the textual representation, such as "1 + 2 + 3" and "12 + 34".
// "Value" is the value of the expression, such as 6 and 46.
type Expression struct {
	Text  string
	Value int
}

// The state struct used in the recursion in the compute function.
// "array" contains the "splits", such as [1, 23, 456].
// "prefixes" contains the prefixes of the current "split", such as
// [{"1 + 23", 24}, {"1 - 23", -23}].
type state struct {
	array    []int
	prefixes []Expression
}

// The operator struct used to store arithmetic operators available.
// "text" is the textual representation of the operator.
// "function" is the actual operator as a function.
// This is not strictly necessary and could be hardcoded,
// but it is here to provide the flexibility for future extentions.
type operator struct {
	text     string
	function func(int, int) int
}

// The operators supported in this version. Only addition and subtraction are
// currently available.
var (
	operators = []operator{
		{"+", func(a, b int) int {
			return a + b
		}},
		{"-", func(a, b int) int {
			return a - b
		}},
	}
)

// digitCount returns 0 if n is negative; otherwise it returns the number of
// decimal digits n has.
//
// Examples:
// (1) digitCount(-1) => 0
// (2) digitCount(0) => 1
// (3) digitCount(1) => 1
// (4) digitCount(12321) => 5
func digitCount(n int) uint {
	if n < 0 {
		return 0
	}
	result := uint(1)
	for n2 := 10; n2 <= n; n2, result = n2*10, result+1 {
	}
	return result
}

// split returns a closure. When the closure is invoked, it returns a "split" of
// the decimal representation of n, and a bool representing whether the next
// invocation produces a new "split".
//
// Example (sp := split(12345)):
// (1) sp() => 1234, 5, true
// (2) sp() => 123, 45, true
// (3) sp() => 12, 345, true
// (4) sp() => 1, 2345, true
// (5) sp() => 0, 12345, false
// (6) sp() => 0, 12345, false
func split(n int) func() (int, int, bool) {
	if n <= 0 {
		return func() (int, int, bool) {
			return 0, 0, false
		}
	}
	d := int(10)
	return func() (int, int, bool) {
		n1, n2 := n/d, n%d
		if n1 == 0 {
			return n1, n2, false
		}
		d *= 10
		return n1, n2, true
	}
}

// compute is the private function called by Compute.
// It provides the recursive algorithm in the function.
func compute(channel chan<- Expression, s *state) {
	last := s.array[len(s.array)-1]
	for _, pfx := range s.prefixes {
		for _, op := range operators {
			text := fmt.Sprintf("%s %s %d", pfx.Text, op.text, last)
			value := op.function(pfx.Value, last)
			channel <- Expression{text, value}
		}
	}

	length := len(s.array)
	sp := split(s.array[length-1])
	for {
		n1, n2, ok := sp()
		if !ok {
			break
		}
		arrayNew := append(s.array[:length-1], n1, n2)
		var prefixesNew []Expression
		if len(s.prefixes) == 0 {
			text := fmt.Sprintf("%d", n1)
			value := n1
			prefixesNew = []Expression{{text, value}}
		} else {
			prefixesNew = make([]Expression, 0, len(s.prefixes)*len(operators))
			for _, pfx := range s.prefixes {
				for _, op := range operators {
					text := fmt.Sprintf("%s %s %d", pfx.Text, op.text, n1)
					value := op.function(pfx.Value, n1)
					prefixesNew = append(prefixesNew, Expression{text, value})
				}
			}
		}
		stateNew := &state{arrayNew, prefixesNew}
		compute(channel, stateNew)
	}
}

// Compute generates all permutations of inserting "+" or "-" between decimal
// digits of the digits parameter. It returns a channel of Expression. The
// permutations are sent to the channel one-by-one. After all permutations are
// sent, the channel is closed.
//
// Example (ch := Compute(123)):
// (1) <-ch => Expression{"12 + 3", 15}, true
// (2) <-ch => Expression{"12 - 3", 9}, true
// (3) <-ch => Expression{"1 + 23", 24}, true
// (4) <-ch => Expression{"1 - 23", -22}, true
// (5) <-ch => Expression{"1 + 2 + 3", 6}, true
// (6) <-ch => Expression{"1 + 2 - 3", 0}, true
// (7) <-ch => Expression{"1 - 2 + 3", 2}, true
// (8) <-ch => Expression{"1 - 2 - 3", -4}, true
// (9) <-ch => Expression{"", 0}, false
func Compute(digits int) <-chan Expression {
	channel := make(chan Expression)
	dc := digitCount(digits)
	capacity := 1 << (dc - 1)
	array := make([]int, 1, capacity)
	array[0] = digits
	prefixes := []Expression{}
	s := &state{array, prefixes}
	go func() {
		compute(channel, s)
		close(channel)
	}()
	return channel
}
