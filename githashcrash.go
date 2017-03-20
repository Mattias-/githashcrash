package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type gitUser struct {
	name  string
	email string
	date  string
}

func parseUserLine(line string) gitUser {
	author := gitUser{}
	author_line := strings.Split(line, " ")
	author.date = strings.Join(author_line[len(author_line)-2:], " ")
	author.email = strings.Trim(author_line[len(author_line)-3 : len(author_line)-2][0], "<>")
	author.name = strings.Join(author_line[1:len(author_line)-3], " ")
	return author
}

func parseObj(obj []byte) (gitUser, gitUser) {
	lines := strings.Split(string(obj), "\n")
	var author gitUser
	var committer gitUser
	for _, line := range lines {
		if strings.HasPrefix(line, "author ") {
			author = parseUserLine(line)
		}
		if strings.HasPrefix(line, "committer ") {
			committer = parseUserLine(line)
		}
	}
	return author, committer
}

func printRecreate(author gitUser, committer gitUser, extra string) {
	fmt.Println("Recreate with:")
	envString := strings.Join([]string{
		"export",
		fmt.Sprintf("GIT_AUTHOR_DATE='%s'", author.date),
		fmt.Sprintf("GIT_AUTHOR_NAME='%s'", author.name),
		fmt.Sprintf("GIT_AUTHOR_EMAIL='%s'", author.email),
		fmt.Sprintf("GIT_COMMITTER_DATE='%s'", committer.date),
		fmt.Sprintf("GIT_COMMITTER_NAME='%s'", committer.name),
		fmt.Sprintf("GIT_COMMITTER_EMAIL='%s'", committer.email),
	}, " ")
	fmt.Printf("(%s; printf '%%s\\n%s' \"$(git show -s --format=%%B)\" | git commit --amend -F -)\n", envString, extra)
}

func worker(targetHash *regexp.Regexp, obj []byte, seed string, result chan string, tested chan int) {
	hashString := ""
	extra := ""
	i := 0
	for ; !targetHash.MatchString(hashString); i++ {
		extra = fmt.Sprintf("%s-%d", seed, i)
		new_obj_len := len(obj) + len(extra) + 1

		h := sha1.New()
		io.WriteString(h, "commit ")
		io.WriteString(h, strconv.Itoa(new_obj_len))
		io.WriteString(h, "\x00")
		h.Write(obj)
		io.WriteString(h, extra)
		io.WriteString(h, "\n")
		hashString = hex.EncodeToString(h.Sum(nil))

		if i%100000 == 0 {
			select {
			case tested <- 100000:
			default:
			}
		}
	}
	log.Println("Found:", hashString)
	result <- extra
}

func main() {
	args := os.Args
	var targetHash = regexp.MustCompile(args[1])

	var obj []byte
	if len(args) == 3 {
		obj = []byte(args[2])
	} else {
		obj, _ = exec.Command("git", "cat-file", "-p", "HEAD").Output()
	}
	if !bytes.HasSuffix(obj, []byte("\n")) {
		obj = append(obj, "\n"...)
	}
	author, committer := parseObj(obj)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := fmt.Sprintf("%x", r.Intn(99999))
	val, ok := os.LookupEnv("GITHASHCRASH_SEED")
	if ok {
		seed = val
	}

	threads := runtime.NumCPU()
	threads_val, threads_ok := os.LookupEnv("GITHASHCRASH_THREADS")
	if threads_ok {
		threads, _ = strconv.Atoi(threads_val)
	}
	log.Println("Threads:", threads)

	start := time.Now()
	tested := make(chan int, 100000)
	sum := 0
	go func() {
		for {
			sum += <-tested
		}
	}()
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for range ticker.C {
			log.Println("Tested", sum)
			elapsed := time.Since(start)
			log.Println("HPS:", float64(sum)/elapsed.Seconds())
		}
	}()

	results := make(chan string)
	for c := 0; c < threads; c++ {
		go worker(targetHash, obj, fmt.Sprintf("%s-%d", seed, c), results, tested)
	}
	extra := <-results
	ticker.Stop()
	elapsed := time.Since(start)
	log.Println("Time:", elapsed)
	log.Println("Commits tested:", sum)
	log.Println("Tests per second:", float64(sum)/elapsed.Seconds())

	printRecreate(author, committer, extra)
}
