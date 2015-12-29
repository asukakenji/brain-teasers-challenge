# Section 1 - Question 2

This is the solution for the compute to 100 challenge.

## Running the Source

#### Run the source directly (a temporary binary is built behind the hood):

```sh
cd brain-teasers-challenge
go run s1q2/main.go
```

#### Build the source and run the binary:

```sh
cd brain-teasers-challenge
go build s1q2/main.go
./main
```

#### Run the unit tests:

```sh
cd brain-teasers-challenge
go test ./s1q2/lib
```

#### Run the benchmarks:

```sh
cd brain-teasers-challenge
go test -bench . ./s1q2/lib
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

## Skills Used

In languages like Python and ECMAScript 6 / JavaScript, generator functionality
is available via the `yield` keyword. `yield` essentially suspends the execution
of the current function, and return to the caller, with an optional return
value. The execution context (call stack, program counter, etc) is not
destroyed. The execution could be resumed at a later time, at where it was
suspended.

Go does not have the `yield` keyword, but similar functionality could be
achieved in at least 2 ways: using closures and using goroutines and channels.

Closure-based generators use less resource, but could only apply to simpler
algorithms. For more complicated ones, it is difficult to "twist" the algorithm
to fit the "execution resumption" pattern. Since `return`, instead of `yield`,
is used to return control to the caller, the execution context will be
destroyed after return. Therefore, it has to be saved manually, so that in the
next invocation there is enough information to reproduce the same state as in
the last invocation before `return`-ing. In my solution, `split()` is a
closure-based generator.

Channel-based generators are more robust, but use more resources. A
`yield`-based algorithm could be ported to channel-based almost without any
changes. An extra advantage is that goroutines, like threads, are executed
concurrently, making the "producer" and "consumer" to run simultaneously. In
my solution, `Compute()` is a channel-based generator.

## Approach

#### split

```
  Input: 12345
 Output: 1234, 5, true
         123, 45, true
		 12, 345, true
		 1, 2345, true
		 0, 12345, false
		 0, 12345, false
		 ...
```

#### Compute

```
  Input: 12345
 Output: 1234, 5
         123, 45
         123, 4, 5
		 12, 345
		 12, 34, 5
		 12, 3, 45
		 12, 3, 4, 5
		 1, 2345
		 1, 234, 5
		 1, 23, 45
		 1, 23, 4,5
		 1, 2, 345
		 1, 2, 34, 5
		 1, 2, 3, 45
		 1, 2, 3, 4, 5
```
