# Section 2 - Question 1

This is the solution for the MapReduce challenge.

## Running the Source

#### Run the source directly (a temporary binary is built behind the hood):

```sh
cd brain-teasers-challenge
go run s2q1/main.go s2q1/input.txt <mapper-count> <partitioner-count> <reducer-count>
```

Example:

```sh
cd brain-teasers-challenge
go run s2q1/main.go s2q1/input.txt 4 1 4
```

#### Build the source and run the binary:

```sh
cd brain-teasers-challenge
go build s2q1/main.go
./main s2q1/input.txt <mapper-count> <partitioner-count> <reducer-count>
```
Example:

```sh
cd brain-teasers-challenge
go build s2q1/main.go
./main s2q1/input.txt 4 1 4
```

#### Run the unit tests: (not yet available)

```sh
cd brain-teasers-challenge
go test ./s2q1/lib
```

#### Run the benchmarks: (not yet available)

```sh
cd brain-teasers-challenge
go test -bench . ./s2q1/lib
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

## Approach (TODO)
