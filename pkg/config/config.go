package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type Config struct {
	Cpuprofile   string
	Seed         []byte
	Placeholder  []byte
	Object       []byte
	MatcherInput string
	Threads      int
	MetricsPort  string
}

func ParseFlags(c *Config) {
	flag.StringVar(&c.Cpuprofile, "cpuprofile", "", "write cpu profile to `file`")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := flag.Int("seed", r.Intn(99999), "Seed bytes, default is actually random.")
	placeholder := flag.String("placeholder", "REPLACEME", "placeholder to mutate")
	flag.IntVar(&c.Threads, "threads", runtime.NumCPU(), "threads")
	flag.StringVar(&c.MetricsPort, "metrics-port", "", "Expose metrics on port.")

	flag.Parse()
	c.Seed = []byte(fmt.Sprintf("%x", seed))
	c.Placeholder = []byte(*placeholder)

	args := flag.Args()
	c.MatcherInput = args[0]
	if len(args) == 2 {
		var err error
		if args[1] == "-" {
			c.Object, err = ioutil.ReadAll(os.Stdin)
		} else {
			c.Object, err = ioutil.ReadFile(args[1])
		}
		if err != nil {
			log.Fatal(err)
		}
	} else {
		obj, err := exec.Command("git", "cat-file", "-p", "HEAD").Output()
		if err != nil {
			log.Fatal("Could not run git command:", err)
		}
		c.Object = obj
	}
}
