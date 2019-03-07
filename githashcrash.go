package main

import (
	regexpmatcher "githashcrash/matcher/regexp"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
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

func run(hashRe string, obj []byte, seed []byte, threads int, placeholder []byte) Result {
	matcher := regexpmatcher.New(hashRe)
	var workers []*Worker
	for i := 0; i < threads; i++ {
		workers = append(workers, &Worker{0})
	}

	results := make(chan Result)
	for c, w := range workers {
		go w.work(matcher, obj, append(seed[:2], byte(c)), placeholder, results)
	}

	// Log stats during execution
	start := time.Now()
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			getStats(start, workers)
		}
	}()

	// Hang here until a result is sent
	extra := <-results
	getStats(start, workers)
	return extra
}

func main() {
	c := config{}
	parseFlags(&c)
	if c.cpuprofile != "" {
		f, err := os.Create(c.cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Printf("Got shutdown signal.")
		pprof.StopCPUProfile()
		os.Exit(1)
	}()

	result := run(c.fillerInput, c.object, c.seed, c.threads, c.placeholder)
	log.Println("Found:", result.sha1)
	printRecreate(result)
}
