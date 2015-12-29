package main

import (
	"./lib"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// "count" and "keywordToFind" are variables appear in the original Node.JS source.
var (
	count         = uint(0)
	keywordToFind = "TCP"
)

// parseArgs parses the command line arguments.
// The return values are:
// (1) the path of "input.txt";
// (2) the mapper count;
// (3) the partitioner count;
// (4) the reducer count.
// If parseArgs succeeds, it returns the above values and a nil error;
// Otherwise, it return zero values for the above values and a non-nil error.
func parseArgs() (string, uint, uint, uint, error) {
	if len(os.Args) != 5 {
		return "", 0, 0, 0, fmt.Errorf("Incorrect number of command line arguments")
	}
	inputTxt := os.Args[1]
	mapperCount, err := strconv.ParseUint(os.Args[2], 10, 0)
	if err != nil {
		return "", 0, 0, 0, err
	}
	partitionerCount, err := strconv.ParseUint(os.Args[3], 10, 0)
	if err != nil {
		return "", 0, 0, 0, err
	}
	reducerCount, err := strconv.ParseUint(os.Args[4], 10, 0)
	if err != nil {
		return "", 0, 0, 0, err
	}
	return inputTxt, uint(mapperCount), uint(partitionerCount), uint(reducerCount), nil
}

// printUsage prints a usage reminder for this command to standard error.
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <input-file> <mapper-count> <partitioner-count> <reducer-count>\n", os.Args[0])
}

// printResult prints the global variable "count" to standard output.
// It is a function which appears in the original Node.JS source.
func printResult() {
	fmt.Println(count)
}

func main() {
	inputTxt, mapperCount, partitionerCount, reducerCount, err := parseArgs()
	if err != nil {
		printUsage()
		os.Exit(1)
		return
	}

	file, err := os.Open(inputTxt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file not found\n")
		os.Exit(2)
		return
	}
	defer file.Close()

	// Since "map" is a keyword in Go, "mapFunc" is used in the following line
	mapFunc, query, queryAll, shutdown := lib.Setup(mapperCount, partitionerCount, reducerCount)
	defer shutdown()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for lineNo := uint(0); scanner.Scan(); lineNo++ {
		line := scanner.Text()
		//fmt.Println(line)
		mapFunc(lineNo, line)
	}

	// OUTPUT ALL KEYS AND THEIR FINAL COUNT
	dictionary := queryAll()
	for word, count := range dictionary {
		fmt.Println(word, count)
	}

	// The global variable "count"
	count = query(keywordToFind)
	printResult()
}
