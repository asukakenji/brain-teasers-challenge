package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
	"strings"
)

// The messageType enum used in the message struct.
type messageType uint

const (
	MAP messageType = iota
	QUERY
)

// The result struct used for returning a result from a query.
// "word" is the word to be counted.
// "count" is the number of occurrences of the word.
type result struct {
	word  string
	count uint
}

// The message struct used for sending a query to the map-reduce network.
// "Type" is the type of the message, which could be MAP or QUERY.
// "content" is the content of the message:
// (1) If it is sent to the mapper, "content" is a line from the input document;
// (2) If it is sent to the partitioner or the reducer, "content" is a word,
//     which represents an occurrence of the word (when "Type" is MAP),
//     or the word to be found (when "Type" is QUERY);
// (3) If the "Type" is QUERY, and the "content" is an empty string (""),
//     the query means "show all values in the dictionary".
// "replyChannel" is used by the reducer to reply the result when "Type" is QUERY.
type message struct {
	// Since "type" is a keyword in Go, "Type" is used in the following line
	Type         messageType
	content      string
	replyChannel chan<- result
}

// "count" and "keywordToFind" are variables appear in the original Node.JS source.
var (
	count         = uint(0)
	keywordToFind = "TCP"
)

// mapperRoutine is the function executed by the mapper goroutines.
// Each mapper goroutine executes the same function, but with different parameters.
// When a line is received from mapperChannel, it is split by the function to words.
// Each word is then sent to partitionerChannel for further processing.
func mapperRoutine(
	mapperChannel <-chan string,
	partitionerChannel chan<- message,
	syncChannel chan<- bool) {

	for line := range mapperChannel {
		reader := strings.NewReader(line)
		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			word := scanner.Text()
			partitionerChannel <- message{MAP, word, nil}
		}
	}
	syncChannel <- true
}

// partitionerRoutine is the function executed by the partitioner goroutines.
// Each partitioner goroutine executes the same function, with same parameters.
// There are three types of requests:
// (1) If the "Type" is MAP, the hash code of "content" is calculated,
//     and the message is forworded to one of the reducerChannels determined by the hash code;
// (2) If the "Type" is QUERY and the "content" is not an empty string (""),
//     the hash code of "content" is calculated, and the message is forwarded like (1);
// (3) If the "Type" is QUERY and the "content" is an empty string (""),
//     the message is forwarded to all reducerChannels.
func partitionerRoutine(
	partitionerChannel <-chan message,
	reducerChannels []chan<- message,
	syncChannel chan<- bool) {

	// hashCode returns the FNV-1a hash of a string as an unsigned 32-bit integer.
	// What "FNV-1a" really is is not important.
	// This function could be treated as the same as Object#hashCode() in Java.
	hashCode := func(s string) uint32 {
		h := fnv.New32a()
		h.Write([]byte(s))
		return h.Sum32()
	}
	reducerCount := uint32(len(reducerChannels))

	for msg := range partitionerChannel {
		word := msg.content
		hc := hashCode(word)
		switch msg.Type {
		case MAP:
			reducerChannels[hc%reducerCount] <- msg
		case QUERY:
			if word == "" {
				for _, rc := range reducerChannels {
					rc <- msg
				}
			} else {
				reducerChannels[hc%reducerCount] <- msg
			}
		default:
			panic(fmt.Sprintf("Unknown case: %d", msg.Type))
		}
	}
	syncChannel <- true
}

// reducerRoutine is the function executed by the reducer goroutines.
// Each reducer goroutine executes the same function, but with different parameters.
// There are three types of requests:
// (1) If the "Type" is MAP, the count of the "content" is increased by 1;
// (2) If the "Type" is QUERY and the "content" is not an empty string (""),
//     the count of the "content" is sent to "replyChannel";
// (3) If the "Type" is QUERY and the "content" is an empty string (""),
//     the whole dictionary is sent to the "replyChannel",
//     and a "zero value" for the result struct is sent at the end to signal termination.
func reducerRoutine(
	reducerChannel <-chan message,
	syncChannel chan<- bool) {

	dictionary := make(map[string]uint)
	for msg := range reducerChannel {
		word := msg.content
		switch msg.Type {
		case MAP:
			dictionary[word]++
		case QUERY:
			rc := msg.replyChannel
			if word == "" {
				for word, count := range dictionary {
					rc <- result{word, count}
				}
				rc <- result{"", 0}
			} else {
				rc <- result{word, dictionary[word]}
			}
		default:
			panic(fmt.Sprintf("Unknown case: %d", msg.Type))
		}
	}
	syncChannel <- true
}

// setup sets-up all the channels required to build the map-reduce network, and
// creates all the functions required to interact with it. The channels are
// encapsulated in the closures and are not visible outside.
//
// The return values are:
// (1) the "map" function which accepts a key (line number) and a value (line);
// (2) the "query" function which accepts a keyword and returns its count;
// (3) the "queryAll" function which returns the whole dictionary by gathering
//     the local dictionaries from all reducers;
// (4) the "shutdown" function which should be called to gracefully terminate
//     the map-reduce network.
func setup(mapperCount, partitionerCount, reducerCount uint) (
	func(uint, string),
	func(string) uint,
	func() map[string]uint,
	func()) {

	// Setup the "syncChannel".
	//
	// The "syncChannel" is used for waiting the code to terminate gracefully.
	// When the channel-receiving code exits the loop, it sends "true" to
	// "syncChannel" before returning (aka. terminating the goroutine). The
	// caller waits for this message to ensure that the code is terminated
	// gracefully.
	syncChannel := make(chan bool)

	// Setup the "reducerChannels".
	//
	// The "reducerChannels" are used by "partitionerRoutine" / "reducerRoutine"
	// to communicate. "partitionerRoutine" knows all "reducerChannels" and
	// selects which one to communicate each time, while "reducerRoutine" knows
	// only one of them that are designated to the particular goroutine (similar
	// to thread) which executes it.
	//
	// "reducerChannelsForInput" reference to the same instances as
	// "reducerChannels", but with a different data type - they are "send-only".
	// Specifying the direction of channels is considered a good habit and helps
	// the code to be understood.
	reducerChannels := make([]chan message, reducerCount)
	reducerChannelsForInput := make([]chan<- message, reducerCount)
	for i, _ := range reducerChannels {
		reducerChannels[i] = make(chan message)
		reducerChannelsForInput[i] = reducerChannels[i]
		go reducerRoutine(reducerChannels[i], syncChannel)
	}

	// Setup the "partitionerChannel".
	//
	// The "partitionerChannel" is used by "mapperRoutine" / "partitionerRoutine"
	// to communicate.
	partitionerChannel := make(chan message)
	for i := uint(0); i < partitionerCount; i++ {
		go partitionerRoutine(partitionerChannel, reducerChannelsForInput, syncChannel)
	}

	// Setup the "mapperChannels".
	//
	// The "partitionerChannel" is used by "mapFunc" / "mapperRoutine" to
	// communicate.
	mapperChannels := make([]chan string, mapperCount)
	for i, _ := range mapperChannels {
		mapperChannels[i] = make(chan string)
		go mapperRoutine(mapperChannels[i], partitionerChannel, syncChannel)
	}

	// Since "map" is a keyword in Go, "mapFunc" is used in the following line
	mapFunc := func(lineNo uint, line string) {
		mapperChannels[lineNo%mapperCount] <- line
	}

	query := func(word string) uint {
		replyChannel := make(chan result)
		defer close(replyChannel)

		partitionerChannel <- message{QUERY, word, replyChannel}
		res := <-replyChannel
		count := res.count
		return count
	}

	queryAll := func() map[string]uint {
		replyChannel := make(chan result)
		defer close(replyChannel)

		partitionerChannel <- message{QUERY, "", replyChannel}
		dictionary := make(map[string]uint)
		channelCount := uint(0)
		for res := range replyChannel {
			if res.word == "" {
				channelCount++
				if channelCount == reducerCount {
					break
				}
			} else {
				dictionary[res.word] += res.count
			}
		}
		return dictionary
	}

	shutdown := func() {
		// Terminate the mapperRoutines gracefully
		for _, mc := range mapperChannels {
			close(mc)
			<-syncChannel
		}

		// Terminate the partitionerChannels gracefully
		close(partitionerChannel)
		for i := uint(0); i < partitionerCount; i++ {
			<-syncChannel
		}

		// Terminate the reducerRoutines gracefully
		for _, rc := range reducerChannels {
			close(rc)
			<-syncChannel
		}

		// Close syncChannel
		close(syncChannel)
	}

	return mapFunc, query, queryAll, shutdown
}

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
	mapFunc, query, queryAll, shutdown := setup(mapperCount, partitionerCount, reducerCount)
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
