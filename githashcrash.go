package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	filler "githashcrash/filler/base"
	matcher "githashcrash/matcher/startswith"
	"githashcrash/worker"
)

func getStats(start time.Time, workers []worker.Worker) {
	var sum uint64
	for _, w := range workers {
		sum += w.Count()
	}
	elapsed := time.Since(start)
	log.Println("Time:", elapsed)
	log.Println("Tested:", sum)
	log.Println("HPS:", float64(sum)/elapsed.Seconds())
}

func run(hashRe string, obj []byte, seed []byte, threads int, placeholder []byte) worker.Result {
	matcher := matcher.New(hashRe)
	log.Println("Workers:", threads)
	results := make(chan worker.Result)
	var workers []worker.Worker
	for i := 0; i < threads; i++ {
		w := worker.NewW()
		workers = append(workers, w)
		filler := filler.New(append(seed[:2], byte(i)))
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
	defer getStats(start, workers)

	return <-results
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
	log.Println("Found:", result.Sha1)
	printRecreate(result)
}
