package main

import (
	filler "githashcrash/filler/base"
	matcher "githashcrash/matcher/startswith"
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

type Filler interface {
	Fill(uint64)
	OutputBuffer() *[]byte
}

type Worker interface {
	Count() uint64
	Work(Matcher, Filler, []byte, []byte, chan Result)
}

type Result struct {
	sha1   string
	object []byte
}

func getStats(start time.Time, workers []Worker) {
	var sum uint64
	for _, w := range workers {
		sum += w.Count()
	}
	elapsed := time.Since(start)
	log.Println("Time:", elapsed)
	log.Println("Tested:", sum)
	log.Println("HPS:", float64(sum)/elapsed.Seconds())
}

func run(hashRe string, obj []byte, seed []byte, threads int, placeholder []byte) Result {
	matcher := matcher.New(hashRe)
	log.Println("Workers:", threads)
	var workers []Worker
	for i := 0; i < threads; i++ {
		workers = append(workers, NewW())
	}

	results := make(chan Result)
	for c, w := range workers {
		filler := filler.New(append(seed[:2], byte(c)))
		go w.Work(matcher, filler, obj, placeholder, results)
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
