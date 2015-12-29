# Section 1 - Question 1

This is the solution for the array partitioning challenge.

## Running the Source

#### Run the source directly (a temporary binary is built behind the hood):

```sh
cd brain-teasers-challenge
go run s1q1/main.go
```

#### Build the source and run the binary:

```sh
cd brain-teasers-challenge
go build s1q1/main.go
./main
```

#### Run the unit tests:

```sh
cd brain-teasers-challenge
go test ./s1q1/lib
```

#### Run the benchmarks:

```sh
cd brain-teasers-challenge
go test -bench . ./s1q1/lib
```

## Source Organization

- `main.go` contains the `main()` function and other functions which are
problem-specific and are not designed to be reusable for other projects.

- `lib/lib.go` contains public functions (with names starting with an uppercase
  letter) that are called by the `main` package, and private functions (with
  names starting with a lowercase letter) that are used internally by the
  library.

- `lib/lib_test.go` contains unit tests and benchmarks.

- After running the benchmark, the "best" implementation of a function is
  selected and used in `lib/lib.go`. The criteria include balancing execution
  speed and readability. The unoptimized versions are available in
  `lib/lib_test.go`, as benchmark case 0 (reference implementation) or case 1
  (unoptimized implementation).

## Approach

(0) Set `i = 0` and `j = len(arr) - 1`:
```
    ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
arr │ 7 │ 7 │ 4 │ 0 │ 9 │ 8 │ 2 │ 4 │ 1 │ 9 │
    └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
      ↑                                   ↑
      i                                   j
```

(1) Increse `i` until `arr[i]` is an odd number.
It is already done in this case.

(2) Decrease `j` until `arr[j]` is an even number:
```
    ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
arr │ 7 │ 7 │ 4 │ 0 │ 9 │ 8 │ 2 │ 4 │ 1 │ 9 │
    └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
      ↑                           ↑
      i                           j
```

(3) If `i >= j`, then it is done. It is not the case this time.

(4) Swap `arr[i]` and `arr[j]`:
```
    ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
arr │ 4 │ 7 │ 4 │ 0 │ 9 │ 8 │ 2 │ 7 │ 1 │ 9 │
    └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
      ↑                           ↑
      i                           j
```

(5) Repeat (1) ~ (4) above:
```
(1) Increase i
    ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
arr │ 4 │ 7 │ 4 │ 0 │ 9 │ 8 │ 2 │ 7 │ 1 │ 9 │
    └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
          ↑                       ↑
          i                       j

(2) Decrease j
    ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
arr │ 4 │ 7 │ 4 │ 0 │ 9 │ 8 │ 2 │ 7 │ 1 │ 9 │
    └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
          ↑                   ↑
          i                   j

(4) Swap
    ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
arr │ 4 │ 2 │ 4 │ 0 │ 9 │ 8 │ 7 │ 7 │ 1 │ 9 │
    └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
          ↑                   ↑
          i                   j

(1) Increase i
    ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
arr │ 4 │ 2 │ 4 │ 0 │ 9 │ 8 │ 7 │ 7 │ 1 │ 9 │
    └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
                      ↑       ↑
                      i       j

(2) Decrease j
    ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
arr │ 4 │ 2 │ 4 │ 0 │ 9 │ 8 │ 7 │ 7 │ 1 │ 9 │
    └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
                      ↑   ↑
                      i   j

(4) Swap
    ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
arr │ 4 │ 2 │ 4 │ 0 │ 8 │ 9 │ 7 │ 7 │ 1 │ 9 │
    └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
                      ↑   ↑
                      i   j

(1) Increase i
    ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
arr │ 4 │ 2 │ 4 │ 0 │ 8 │ 9 │ 7 │ 7 │ 1 │ 9 │
    └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
                         ↑ ↑
                         i j

(2) Decrease j
    ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
arr │ 4 │ 2 │ 4 │ 0 │ 8 │ 9 │ 7 │ 7 │ 1 │ 9 │
    └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
                      ↑   ↑
                      j   i

(3) Done!
```

## Common Mistakes

#### Unaware of negative integers in the input

```go
// This implementation is incorrect.
// isOdd(-1) should return true, but this implementation returns false.
// The root case is the result of the modulo operation could be negative.
// For example, -1 % 2 returns -1
func isOdd(n int) bool {
	return n%2 == 1
}
```
