package lib

import "fmt"

type Expression struct {
	Text  string
	Value int
}

type state struct {
	array    []int
	prefixes []Expression
}

var (
	operators = map[string]func(int, int) int{
		"+": func(a, b int) int {
			return a + b
		},
		"-": func(a, b int) int {
			return a - b
		},
	}
)

func digitCount(n int) uint {
	if n < 0 {
		return 0
	}
	result := uint(1)
	for n2 := 10; n2 <= n; n2, result = n2*10, result+1 {
	}
	return result
}

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

func compute(channel chan<- Expression, s *state) {
	last := s.array[len(s.array)-1]
	for _, pfx := range s.prefixes {
		for operator, function := range operators {
			text := fmt.Sprintf("%s %s %d", pfx.Text, operator, last)
			value := function(pfx.Value, last)
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
				for operator, function := range operators {
					text := fmt.Sprintf("%s %s %d", pfx.Text, operator, n1)
					value := function(pfx.Value, n1)
					prefixesNew = append(prefixesNew, Expression{text, value})
				}
			}
		}
		stateNew := &state{arrayNew, prefixesNew}
		compute(channel, stateNew)
	}
}

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
