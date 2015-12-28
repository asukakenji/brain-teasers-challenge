package lib

import "fmt"

var _ = fmt.Println // TODO: Delete this!

type prefix struct {
	text  string
	value int
}

/*
func (p prefix) String() string {
	return fmt.Sprintf("%q(%d)", p.text, p.value)
}
*/

type state struct {
	array    []int
	prefixes []prefix
}

/*
func (s *state) print() {
	fmt.Printf("array @ %p : %v\n", s.array)
	fmt.Printf("prefixes @ %p : %v\n", s.prefixes)
	fmt.Println()
}
*/

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

func compute(channel chan<- string, sum int, s *state) {
	last := s.array[len(s.array)-1]
	for _, pfx := range s.prefixes {
		for operator, function := range operators {
			if function(pfx.value, last) == sum {
				channel <- fmt.Sprintf("%s %s %d", pfx.text, operator, last)
			}
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
		var prefixesNew []prefix
		if len(s.prefixes) == 0 {
			text := fmt.Sprintf("%d", n1)
			value := n1
			prefixesNew = []prefix{{text, value}}
		} else {
			prefixesNew = make([]prefix, 0, len(s.prefixes)*len(operators))
			for _, pfx := range s.prefixes {
				for operator, function := range operators {
					text := fmt.Sprintf("%s %s %d", pfx.text, operator, n1)
					value := function(pfx.value, n1)
					prefixesNew = append(prefixesNew, prefix{text, value})
				}
			}
		}
		stateNew := &state{arrayNew, prefixesNew}
		compute(channel, sum, stateNew)
	}
}

func Compute(sum, digits int) <-chan string {
	channel := make(chan string)
	dc := digitCount(digits)
	capacity := 1 << (dc - 1)
	array := make([]int, 1, capacity)
	array[0] = digits
	prefixes := []prefix{}
	s := &state{array, prefixes}
	go func() {
		compute(channel, sum, s)
		close(channel)
	}()
	return channel
}
