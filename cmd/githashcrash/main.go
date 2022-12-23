package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/Mattias-/githashcrash/pkg/config"
	filler "github.com/Mattias-/githashcrash/pkg/filler/base"
	matcher "github.com/Mattias-/githashcrash/pkg/matcher/regexp"
	"github.com/Mattias-/githashcrash/pkg/worker"
	"github.com/Mattias-/githashcrash/pkg/worker/commitmsg"
)

var workers []worker.Worker

func printStats(start time.Time) {
	sum := hashCount()
	elapsed := time.Since(start).Round(time.Second)
	log.Println("Time:", elapsed.String())
	log.Println("Tested:", sum)
	log.Printf("%.2f MH/s\n", sum/1_000_000/elapsed.Seconds())
}

func hashCount() float64 {
	var sum float64
	for _, w := range workers {
		sum += float64(w.Count())
	}
	return sum
}

func main() {
	c := config.GetConfig()
	if c.Cpuprofile != "" {
		f, err := os.Create(c.Cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if c.MetricsPort != "" {
		go func() {
			http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(w, `# HELP hashcount_total How many Hashes has been tested.
# TYPE hashcount_total counter
hashcount_total`, hashCount())
			})
			log.Fatal(http.ListenAndServe(c.MetricsPort, nil)) // #nosec G114
		}()
	}

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Println("Got shutdown signal.")
		pprof.StopCPUProfile()
		os.Exit(1)
	}()

	log.Println("Workers:", c.Threads)

	results := make(chan worker.Result)
	for i := 0; i < c.Threads; i++ {
		filler := filler.New(append(c.Seed[:2], byte(i)))
		matcher := matcher.New(c.MatcherInput)
		w := commitmsg.NewWorker(matcher, filler, c.Object, c.Placeholder)
		workers = append(workers, w)
		go w.Work(results)
	}

	// Log stats during execution
	start := time.Now()
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for range ticker.C {
			printStats(start)
		}
	}()

	result := <-results

	ticker.Stop()
	printStats(start)

	log.Println("Found:", result.Sha1())
	fmt.Println(result.ShellRecreateCmd())
}
