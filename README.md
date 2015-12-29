# Brain Teasers Challenge

## Installing Go

Go is available at:

https://golang.org/

Click "Download Go" and select a binary distribution for your platform.
Execute the binary and following the instructions.

## Downloading the source

#### With Git (preferred):

The source could be downloaded with Git by executing the following command:

```
git clone https://github.com/asukakenji/brain-teasers-challenge
```

#### As a ZIP via HTTP:

The latest version of the source is available at:

https://github.com/asukakenji/brain-teasers-challenge

Click "Download ZIP" on the right hand side of the screen to download the archive.
Unzip the archive to get the source code.

#### With `go get`:

The source cannot be downloaded with `go get`. If you attempt to do so, `go`
will complain that local imports are used. Since the packages are not assumed to
be fetched from the web and imported into a project other than this challenge, I
used `import "./lib"` instead of a fully-qualified package name like
`"github.com/asukakenji/brain-teasers-challenge/s1q1/lib"`.

## Source code directory structure

```
brain-teasers-challenge/
  ├─s1q1/
  │   ├─lib/
  │   │   ├─lib_test.go
  │   │   └─lib.go
  │   ├─main.go
  │   └─README.md
  ├─s1q2/
  │   ├─lib/
  │   │   ├─lib_test.go
  │   │   └─lib.go
  │   ├─main.go
  │   └─README.md
  ├─s2q1/
  │   ├─input.txt
  │   ├─lib/
  │   │   ├─lib_test.go
  │   │   └─lib.go
  │   ├─main.go
  │   └─README.md
  └─s2q2/
      ├─lib/
      │   ├─lib_test.go
      │   └─lib.go
      ├─main.go
      └─README.md
  
```

In the original source, the directory names do not correspond to the section and
question numbers in `README.md`. In my source, the directories are named `sMqN/`,
where M is the section number and N is the question number. Here is the mapping:

Question             | Original Directory | My Directory
-------------------- | ------------------ | ------------
Section 1 Question 1 | 1/                 | s1q1/
Section 1 Question 2 | 4/                 | s1q2/
Section 2 Question 1 | 3/                 | s2q1/
Section 2 Question 2 | 2/                 | s2q2/

## Running the source

Please see the `README.md` in the coresponding directory for further information.

## Some points to notice

- Go uses uppercase and lowercase to distinguish between public and private
  functions / methods / fields / etc. The declaration is public if the first
  letter of the identifier is uppercase; it is private otherwise.

- Private declarations are visible only to source code in the same package,
  and it is the primary means in Go to achieve information hiding.

- I also use uppercase for identifiers which collide with a Go keyword (e.g.
  `type` and `map`). I made sure that it does not leak information, because,
  for instance, the public declaration is enclosed by a private scope.

- A goroutine is similar to a thread. Goroutines are executed concurrently, but
  are a lot more lightweight than threads. They communicate through channels
  (preferred) and shared variables.

- In all my solutions, I tried to reduce the problem into smaller pieces.
  Functions are designed to focus on a single aspect and are highly testable.
  They are reusable for solving other problems too. This is in contrast to
  writing one big function to solve the whole problem, or having several smaller
  functions which are only usable to one problem (most probably because they use
  problem-specific parameter / return types).
