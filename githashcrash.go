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

func getStats(start time.Time, workers []*Worker) {
	var sum uint64 = 0
	for _, w := range workers {
		sum += w.i
	}
	elapsed := time.Since(start)
	log.Println("Time:", elapsed)
	log.Println("Tested:", sum)
	log.Println("HPS:", float64(sum)/elapsed.Seconds())
}

func run(hashRe string, obj []byte, seed string, threads int) string {
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

	results := make(chan string)
	for c, w := range workers {
		go w.worker(targetHash, obj, append([]byte(seed[:2]), byte(c)), results)
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

	extra := run(hashRe, obj, seed, threads)
	log.Println("Found:", extra)
	author, committer := parseObj(obj)
	printRecreate(author, committer, extra)

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
