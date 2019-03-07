package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"runtime"
	"time"
)

type config struct {
	cpuprofile  string
	seed        []byte
	placeholder []byte
	object      []byte
	fillerInput string
	threads     int
}

func parseFlags(c *config) {
	flag.StringVar(&c.cpuprofile, "cpuprofile", "", "write cpu profile to `file`")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := flag.Int("seed", r.Intn(99999), "write memory profile to `file`")
	placeholder := flag.String("placeholder", "REPLACEME", "placeholder to mutate")
	flag.IntVar(&c.threads, "threads", runtime.NumCPU(), "threads")

	flag.Parse()
	args := flag.Args()
	c.fillerInput = args[0]
	obj, err := exec.Command("git", "cat-file", "-p", "HEAD").Output()
	if err != nil {
		log.Fatal("Could not run git command:", err)
	}
	c.object = obj
	c.seed = []byte(fmt.Sprintf("%x", seed))
	c.placeholder = []byte(*placeholder)
}
