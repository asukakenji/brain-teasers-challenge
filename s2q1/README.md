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

- (Not yet available)
  `lib/lib_test.go` contains unit tests and benchmarks.

- (Not yet available)
  After running the benchmark, the "best" implementation of a function is
  selected and used in `lib/lib.go`. The criteria include balancing execution
  speed and readability. The unoptimized versions are available in
  `lib/lib_test.go`, as benchmark case 0 (reference implementation) or case 1
  (unoptimized implementation).

## Approach

The MapReduce network is implemented with an "actor model", which is similar to
Erlang processes, and Akka actors in Scala. An actor is an entity that has its
own "unit of execution" (thread, or goroutine in this case). Actors communicate
to each other by sending messages (with Go channels in this case).

The following diagram pictures the goroutines setup when `<mapper-count>` is 3,
`<partitioner-count>` is 2, and `<reducer-count>` is 4:

```
                             main() / mapFunc()
                                     │
                  ┌──────────────────┼──────────────────┐
                  │                  │                  │
           mapperRoutine()    mapperRoutine()    mapperRoutine()
                  │                  │                  │
                  └──────────────────┼──────────────────┘
                                     │
                         ┌───────────┴───────────┐
                         │                       │
                partitionerRoutine()    partitionerRoutine()
                         │                       │
                         └───────────┬───────────┘
                                     │
       ┌───────────────────┬─────────┴─────────┬───────────────────┐
       │                   │                   │                   │ 
reducerRoutine()    reducerRoutine()    reducerRoutine()    reducerRoutine()
```

`<mapper-count>` (3 in this case) goroutines and `<mapper-count>` channels are
created. Each goroutine running `mapperRoutine()` has a dedicated channel, which
is referred to as mapper channel. The `main()` function reads the input file
line-by-line, and calls the `mapFunc()`. The `mapFunc()` sends the lines to the
mapper channels in a round-robin manner: line 0 goes to channel 0, line 1 goes
to channel 1, line 2 goes to channel 2, line 3 goes to channel 0, and so on.

`<partitioner-count>` (2 in this case) goroutines and 1 channel are created. All
goroutines running `partitionerRoutine()` share the same channel, which is
referred to as the partitioner channel. The `mapperRoutine()`s split the lines
into words, and send them to the partitioner channel. A word sent to the
partitioner channel is received by one (and only one) of the
`partitionerRoutine()`s selected arbitrarily.

`<reducer-count>` (4 in this case) goroutines and `<reducer-count>` channels are
created. Each goroutine running `reducerRoutine()` has a dedicated channel,
which is referred to as reducer channel. The `partitionerRoutine()`s calculate
a hash code from each word they receive, and send the word to a reducer channel
corresponding to the hash code. Therefore, the same word is sent to the same
reducer channel and is handled by the same goroutine. After all the words are
iterated, each of the `reducerRoutine()`s has its own dictionary. Given a word
in the original document, it exists in exactly one of the dictionaries.

In the query phase, a "reply channel" is attached to the query. The
`partitionerRoutine()` sends the query to the designated `reducerRoutine()` via
the designated reducer channel. The `reducerRoutine()` sends back the result
through the "reply channel".
