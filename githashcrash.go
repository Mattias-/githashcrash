package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"
)

// Matcher has a match function that returns true if
type Matcher interface {
	Match([]byte) bool
}

type Worker struct {
	i uint64
}

type Result struct {
	sha1   string
	object []byte
}

func getStats(start time.Time, workers []*Worker) {
	var sum uint64
	for _, w := range workers {
		sum += w.i
	}
	elapsed := time.Since(start)
	log.Println("Time:", elapsed)
	log.Println("Tested:", sum)
	log.Println("HPS:", float64(sum)/elapsed.Seconds())
}

func run(hashRe string, obj []byte, seed string, threads int, placeholder []byte) Result {
	var targetHash = regexp.MustCompile(hashRe)
	var workers []*Worker
	for i := 0; i < threads; i++ {
		workers = append(workers, &Worker{0})
	}

	start := time.Now()
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for range ticker.C {
			getStats(start, workers)
		}
	}()

	results := make(chan Result)
	for c, w := range workers {
		go w.work(targetHash, obj, append([]byte(seed[:2]), byte(c)), placeholder, results)
	}
	extra := <-results
	getStats(start, workers)
	ticker.Stop()
	return extra
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	args := os.Args
	hashRe := args[1]

	var obj []byte
	if len(args) == 3 {
		obj = []byte(args[2])
	} else {
		obj, _ = exec.Command("git", "cat-file", "-p", "HEAD").Output()
	}
	// A trailing newline might be lost if object is passed as argument.
	if !bytes.HasSuffix(obj, []byte("\n")) {
		obj = append(obj, "\n"...)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := fmt.Sprintf("%x", r.Intn(99999))
	val, ok := os.LookupEnv("GITHASHCRASH_SEED")
	if ok {
		seed = val
	}

	threads := runtime.NumCPU()
	threadsVal, threadsOk := os.LookupEnv("GITHASHCRASH_THREADS")
	if threadsOk {
		threads, _ = strconv.Atoi(threadsVal)
	}
	log.Println("Threads:", threads)

	placeholder := []byte("REPLACEME")
	placeholderVal, placeholderOk := os.LookupEnv("GITHASHCRASH_PLACEHOLDER")
	if placeholderOk {
		placeholder = []byte(placeholderVal)
	}

	result := run(hashRe, obj, seed, threads, placeholder)
	log.Println("Found:", result.sha1)
	printRecreate(result)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
