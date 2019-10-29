package config

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"runtime"
	"time"
)

type Config struct {
	Cpuprofile  string
	Seed        []byte
	Placeholder []byte
	Object      []byte
	FillerInput string
	Threads     int
}

func ParseFlags(c *Config) {
	flag.StringVar(&c.Cpuprofile, "cpuprofile", "", "write cpu profile to `file`")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := flag.Int("seed", r.Intn(99999), "Seed bytes, default is actually random.")
	placeholder := flag.String("placeholder", "REPLACEME", "placeholder to mutate")
	flag.IntVar(&c.Threads, "threads", runtime.NumCPU(), "threads")

	flag.Parse()
	args := flag.Args()
	c.FillerInput = args[0]
	obj, err := exec.Command("git", "cat-file", "-p", "HEAD").Output()
	if err != nil {
		log.Fatal("Could not run git command:", err)
	}
	c.Object = obj
	c.Seed = []byte(fmt.Sprintf("%x", seed))
	c.Placeholder = []byte(*placeholder)
}
